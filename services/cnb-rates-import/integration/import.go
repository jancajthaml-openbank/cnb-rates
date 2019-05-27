// Copyright (c) 2016-2019, Jan Cajthaml <jan.cajthaml@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jancajthaml-openbank/cnb-rates-import/config"
	"github.com/jancajthaml-openbank/cnb-rates-import/metrics"
	"github.com/jancajthaml-openbank/cnb-rates-import/utils"

	localfs "github.com/jancajthaml-openbank/local-fs"
	log "github.com/sirupsen/logrus"
)

// CNBRatesImport represents cnb gateway rates import subroutine
type CNBRatesImport struct {
	utils.DaemonSupport
	cnbGateway string
	storage    *localfs.Storage
	metrics    *metrics.Metrics
	httpClient Client
}

// NewCNBRatesImport returns cnb rates import fascade
func NewCNBRatesImport(ctx context.Context, cfg config.Configuration, metrics *metrics.Metrics, storage *localfs.Storage) CNBRatesImport {
	return CNBRatesImport{
		DaemonSupport: utils.NewDaemonSupport(ctx),
		storage:       storage,
		cnbGateway:    cfg.CNBGateway,
		metrics:       metrics,
		httpClient:    NewClient(),
	}
}

func (cnb CNBRatesImport) syncMainRateToday(today time.Time) error {
	cachePath := utils.FXMainOfflinePath(today)
	if ok, err := cnb.storage.Exists(cachePath); err != nil {
		return err
	} else if ok {
		return nil
	}

	uri := cnb.cnbGateway + utils.GetUrlForDateMainFx(today)
	response, code, err := cnb.httpClient.Get(uri)
	if code != 200 && err == nil {
		return fmt.Errorf("CNB cloud error %d %+v", code, string(response))
	} else if err != nil {
		return fmt.Errorf("CNB cloud error %d %+v", code, err)
	}

	// FIXME try with backoff until hit
	if !validateRates(today, response) {
		return fmt.Errorf("today rate bounce %s", today.Format("02.01.2006"))
	}

	if cnb.storage.WriteFile(cachePath, response) != nil {
		return fmt.Errorf("cannot store cache for %s at %s", today.Format("02.01.2006"), cachePath)
	}

	cnb.metrics.DayImported()

	log.Debug("downloaded fx-main for today")
	return nil
}

func (cnb CNBRatesImport) syncOtherRates(day time.Time) error {
	cachePath := utils.FXOtherOfflinePath(day)
	if ok, err := cnb.storage.Exists(cachePath); err != nil {
		return fmt.Errorf("corrupted cache at %s with %+v", cachePath, err)
	} else if ok {
		return nil
	}

	uri := cnb.cnbGateway + utils.GetUrlForDateOtherFx(day)
	response, code, err := cnb.httpClient.Get(uri)
	if code != 200 && err == nil {
		return fmt.Errorf("CNB cloud error %d %+v", code, string(response))
	}
	if err != nil {
		return fmt.Errorf("CNB cloud error %d %+v", code, err)
	}

	if cnb.storage.WriteFile(cachePath, response) != nil {
		return fmt.Errorf("cannot store cache for %s at %s", day, cachePath)
	}

	log.Infof("downloaded fx-other for %s", day.Format("02.01.2006"))

	cnb.metrics.DayImported()
	return nil
}

func (cnb CNBRatesImport) syncMainRates(days []time.Time) error {
	if len(days) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	queue := make(chan time.Time, 128)

	worker := func() {
		for {
			select {

			case date, ok := <-queue:
				if !ok {
					return
				}
				if cnb.IsDone() {
					wg.Done()
					continue
				}

				cachePath := utils.FXMainOfflinePath(date)

				ok, err := cnb.storage.Exists(cachePath)
				if err != nil {
					wg.Done()
					log.Warnf("corrupted cache at %s with %+v", cachePath, err)
					continue
				}
				if ok {
					wg.Done()
					continue
				}

				var (
					response []byte
					code     int
				)

				uri := cnb.cnbGateway + utils.GetUrlForDateMainFx(date)
				response, code, err = cnb.httpClient.Get(uri)
				if code != 200 && err == nil {
					log.Warnf("CNB cloud error %d %+v", code, string(response))
					wg.Done()
					continue
				} else if err != nil {
					log.Warnf("CNB cloud error %d %+v", code, err)
					wg.Done()
					continue
				}

				if cnb.storage.WriteFile(cachePath, response) != nil {
					log.Warnf("cannot store cache for %s at %s", date, cachePath)
					wg.Done()
					continue
				}

				log.Infof("downloaded fx-main for day %s", date.Format("02.01.2006"))
				cnb.metrics.DayImported()
				wg.Done()
			}
		}
	}

	for i := 0; i < 4*runtime.NumCPU(); i++ {
		go worker()
	}

	for _, day := range days {
		cachePath := utils.FXMainOfflinePath(day)
		ok, err := cnb.storage.Exists(cachePath)
		if err != nil {
			log.Warnf("corrupted cache at %s with %+v", cachePath, err)
			continue
		}
		if ok {
			continue
		}

		wg.Add(1)
		queue <- day
	}

	wg.Wait()
	close(queue)

	return nil
}

func validateRates(date time.Time, data []byte) bool {
	if data == nil || len(data) < 20 {
		return false
	}

	b := data[0:20]
	chunk := make([]byte, 0)

	if b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		chunk = append(chunk, b[3:]...)
	} else {
		chunk = append(chunk, b...)
	}

	chunk = append(chunk, data[20:]...)

	j := bytes.IndexByte(chunk, '\n')
	if j < 0 {
		return false
	}

	parts := strings.Split(string(chunk[:j]), " #")
	expected, err := time.Parse("02 Jan 2006", parts[0])
	if err != nil {
		return false
	}
	return date.Year() == expected.Year() && date.Month() == expected.Month() && date.Day() == expected.Day()
}

func (cnb CNBRatesImport) importRoundtrip() {
	now := time.Now()

	if err := cnb.syncMainRateToday(now); err != nil {
		log.Warnf(err.Error())
	}

	fxMainHistoryStart := time.Date(1991, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	fxOtherHistoryStart := time.Date(2004, time.Month(5), 1, 0, 0, 0, 0, time.UTC)
	today := now.AddDate(0, 0, -1)

	months := utils.GetMonthsBetween(fxMainHistoryStart, today)
	for _, month := range months {
		if cnb.IsDone() {
			return
		}

		currentMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
		nextMonth := time.Date(month.Year(), month.Month()+1, 0, 0, 0, 0, 0, time.UTC)
		nextMonth.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

		days := utils.GetDatesBetween(currentMonth, nextMonth)
		if len(days) <= 2 {
			continue
		}

		if !currentMonth.Before(fxOtherHistoryStart) {
			lastDay := days[len(days)-1]

			log.Debugf("Synchonizing other fx rates for %s", lastDay.Format("02.01.2006"))
			// FIXME must be last day of `currentMonth` for ther fx fates
			if err := cnb.syncOtherRates(lastDay); err != nil {
				log.Warnf(err.Error())
			}
		}

		log.Debugf("Synchonizing main fx rates from %s to %s", days[0].Format("02.01.2006"), days[len(days)-1].Format("02.01.2006"))
		cnb.syncMainRates(days)
	}
}

// WaitReady wait for cnb rates import to be ready
func (cnb CNBRatesImport) WaitReady(deadline time.Duration) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch x := e.(type) {
			case string:
				err = fmt.Errorf(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("unknown panic")
			}
		}
	}()

	ticker := time.NewTicker(deadline)
	select {
	case <-cnb.IsReady:
		ticker.Stop()
		err = nil
		return
	case <-ticker.C:
		err = fmt.Errorf("cnb-rates-import daemon was not ready within %v seconds", deadline)
		return
	}
}

// Start handles everything needed to start cnb rates import daemon
func (cnb CNBRatesImport) Start() {
	defer cnb.MarkDone()

	cnb.MarkReady()

	select {
	case <-cnb.CanStart:
		break
	case <-cnb.Done():
		return
	}

	log.Infof("Start cnb-rates-import daemon, sync %v now", cnb.cnbGateway)

	cnb.importRoundtrip()

	log.Info("Stopping cnb-rates-import daemon")
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	log.Info("Stop cnb-rates-import daemon")

	return
}
