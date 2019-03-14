// Copyright (c) 2016-2018, Jan Cajthaml <jan.cajthaml@gmail.com>
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

package daemon

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
	"github.com/jancajthaml-openbank/cnb-rates-import/http"
	"github.com/jancajthaml-openbank/cnb-rates-import/utils"

	localfs "github.com/jancajthaml-openbank/local-fs"
	log "github.com/sirupsen/logrus"
)

// CNBRatesImport represents cnb gateway rates import subroutine
type CNBRatesImport struct {
	Support
	cnbGateway string
	storage    *localfs.Storage
	metrics    *Metrics
	httpClient http.Client
}

// NewCNBRatesImport returns cnb rates import fascade
func NewCNBRatesImport(ctx context.Context, cfg config.Configuration, metrics *Metrics, storage *localfs.Storage) CNBRatesImport {
	return CNBRatesImport{
		Support:    NewDaemonSupport(ctx),
		storage:    storage,
		cnbGateway: cfg.CNBGateway,
		metrics:    metrics,
		httpClient: http.NewClient(),
	}
}

func (cnb CNBRatesImport) syncOtherRatesMonthly() error {
	now := time.Now()

	startDate := time.Date(1991, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	dates := utils.GetMonthsBetween(startDate, endDate)

	if len(dates) == 0 {
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
				if cnb.ctx.Err() != nil {
					wg.Done()
					continue
				}

				// FIXME to function that accepts wg reference
				cachePath := utils.MonthlyCachePath(date)

				ok, err := cnb.storage.Exists(cachePath)
				if err != nil {
					log.Warnf("corrupted cache at %s with %+v", cachePath, err)
					wg.Done()
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

				uri := cnb.cnbGateway + utils.GetUrlForDateOtherFx(date)

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

				log.Debugf("downloaded other fx for date %s", date.Format("01.2006"))
				wg.Done()
			}
		}
	}

	for i := 0; i < 4*runtime.NumCPU(); i++ {
		go worker()
	}

	for _, date := range dates {
		cachePath := utils.MonthlyCachePath(date)
		ok, err := cnb.storage.Exists(cachePath)
		if err != nil {
			log.Warnf("corrupted cache at %s with %+v", cachePath, err)
			continue
		}
		if ok {
			continue
		}

		wg.Add(1)
		queue <- date
	}

	wg.Wait()
	close(queue)

	return nil
}

func (cnb CNBRatesImport) syncMainRatesDaily() error {
	now := time.Now()

	startDate := time.Date(1991, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	endDate := now.AddDate(0, 0, -1)

	dates := utils.GetDatesBetween(startDate, endDate)

	if len(dates) > 0 {
		var wg sync.WaitGroup
		queue := make(chan time.Time, 128)

		worker := func() {
			for {
				select {

				case date, ok := <-queue:
					if !ok {
						return
					}
					if cnb.ctx.Err() != nil {
						wg.Done()
						continue
					}

					// FIXME to function that accepts wg reference
					cachePath := utils.DailyCachePath(date)

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

					log.Debugf("downloaded main fx for date %s", date.Format("02.01.2006"))
					wg.Done()
				}
			}
		}

		for i := 0; i < 4*runtime.NumCPU(); i++ {
			go worker()
		}

		for _, date := range dates {
			cachePath := utils.DailyCachePath(date)
			ok, err := cnb.storage.Exists(cachePath)
			if err != nil {
				log.Warnf("corrupted cache at %s with %+v", cachePath, err)
				continue
			}
			if ok {
				continue
			}

			wg.Add(1)
			queue <- date
		}

		wg.Wait()
		close(queue)
	}

	cachePath := utils.DailyCachePath(now)

	ok, err := cnb.storage.Exists(cachePath)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	var (
		response []byte
		code     int
	)

	uri := cnb.cnbGateway + utils.GetUrlForDateMainFx(now)
	response, code, err = cnb.httpClient.Get(uri)
	if code != 200 && err == nil {
		return fmt.Errorf("CNB cloud error %d %+v", code, string(response))
	} else if err != nil {
		return fmt.Errorf("CNB cloud error %d %+v", code, err)
	}

	if !validateRates(now, response) {
		return fmt.Errorf("bounce %s", now.Format("02.01.2006"))
	}

	if cnb.storage.WriteFile(cachePath, response) != nil {
		return fmt.Errorf("cannot store cache for %s at %s", now.Format("02.01.2006"), cachePath)
	}

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
	expected, err := time.Parse("02.Jan 2006", parts[0])
	if err != nil {
		return false
	}
	return date.Year() == expected.Year() && date.Month() == expected.Month() && date.Day() == expected.Day()
}

func (cnb CNBRatesImport) importRoundtrip() {
	log.Debug("Synchonizing main currencies fx daily rates")
	cnb.syncMainRatesDaily()

	log.Debug("Synchonizing other currencies fx monthly rates")
	cnb.syncOtherRatesMonthly()
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

	log.Infof("Start cnb-rates-import daemon, sync %v now", cnb.cnbGateway)
	cnb.MarkReady()

	cnb.importRoundtrip()

	log.Info("Stopping cnb-rates-import daemon")
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	log.Info("Stop cnb-rates-import daemon")

	return
}
