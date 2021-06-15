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
	"fmt"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jancajthaml-openbank/cnb-rates-import/integration/cnb"
	"github.com/jancajthaml-openbank/cnb-rates-import/support/timeshift"

	localfs "github.com/jancajthaml-openbank/local-fs"
)

// CNBRatesImport represents cnb gateway rates import subroutine
type CNBRatesImport struct {
	cnbGateway string
	storage    localfs.Storage
	client     *cnb.Client
}

// NewCNBRatesImport returns cnb rates import fascade
func NewCNBRatesImport(gateway string, rootStorage string) *CNBRatesImport {
	storage, err := localfs.NewPlaintextStorage(rootStorage)
	if err != nil {
		log.Error().Msgf("Failed to ensure storage %+v", err)
		return nil
	}
	return &CNBRatesImport{
		storage: storage,
		client:  cnb.NewClient(gateway),
	}
}

func (rates *CNBRatesImport) syncMainRateToday(today time.Time) error {
	if rates == nil {
		return nil
	}
	cachePath := FXMainOfflinePath(today)
	if ok, err := rates.storage.Exists(cachePath); err != nil {
		return err
	} else if ok {
		return nil
	}

	resp, err := rates.client.GetMainFxFor(today)
	if err != nil {
		return fmt.Errorf("sync Main Rates error %+v", err)
	}

	// FIXME try with backoff until hit
	if !validateRates(today, resp) {
		return fmt.Errorf("today rate bounce %s", today.Format("02.01.2006"))
	}

	if rates.storage.WriteFile(cachePath, resp) != nil {
		return fmt.Errorf("cannot store cache for %s at %s", today.Format("02.01.2006"), cachePath)
	}

	log.Debug().Msg("downloaded fx-main for today")
	return nil
}

func (rates *CNBRatesImport) syncOtherRates(day time.Time) error {
	if rates == nil {
		return nil
	}

	cachePath := FXOtherOfflinePath(day)
	if ok, err := rates.storage.Exists(cachePath); err != nil {
		return fmt.Errorf("corrupted cache at %s with %+v", cachePath, err)
	} else if ok {
		return nil
	}

	resp, err := rates.client.GetMainFxFor(day)
	if err != nil {
		return fmt.Errorf("sync Other Rates error %+v", err)
	}

	if rates.storage.WriteFile(cachePath, resp) != nil {
		return fmt.Errorf("cannot store cache for %s at %s", day, cachePath)
	}

	log.Info().Msgf("downloaded fx-other for %s", day.Format("02.01.2006"))
	return nil
}

func (rates *CNBRatesImport) syncMainRates(days []time.Time) error {
	if rates == nil {
		return nil
	}

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

				cachePath := FXMainOfflinePath(date)

				ok, err := rates.storage.Exists(cachePath)
				if err != nil {
					wg.Done()
					log.Warn().Msgf("corrupted cache at %s with %+v", cachePath, err)
					continue
				}
				if ok {
					wg.Done()
					continue
				}

				resp, err := rates.client.GetMainFxFor(date)
				if err != nil {
					log.Warn().Msgf("sync Main Rates error %+v", err)
					wg.Done()
					continue
				}

				if rates.storage.WriteFile(cachePath, resp) != nil {
					log.Warn().Msgf("cannot store cache for %s at %s", date, cachePath)
					wg.Done()
					continue
				}

				log.Info().Msgf("downloaded fx-main for day %s", date.Format("02.01.2006"))
				wg.Done()
			}
		}
	}

	for i := 0; i < 4*runtime.NumCPU(); i++ {
		go worker()
	}

	for _, day := range days {
		cachePath := FXMainOfflinePath(day)
		ok, err := rates.storage.Exists(cachePath)
		if err != nil {
			log.Warn().Msgf("corrupted cache at %s with %+v", cachePath, err)
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

func (rates *CNBRatesImport) importRoundtrip() {
	if rates == nil {
		return
	}

	now := time.Now()

	if err := rates.syncMainRateToday(now); err != nil {
		log.Warn().Msg(err.Error())
	}

	fxMainHistoryStart := time.Date(1991, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	fxOtherHistoryStart := time.Date(2004, time.Month(5), 1, 0, 0, 0, 0, time.UTC)
	today := now.AddDate(0, 0, -1)
	yesterday := today.AddDate(0, 0, -1)

	months := timeshift.GetMonthsBetween(fxMainHistoryStart, today)
	for _, month := range months {
		// FIXME context cancel check
		//if rates.IsCanceled() {
		//return
		//}

		currentMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
		nextMonth := time.Date(month.Year(), month.Month()+1, 0, 0, 0, 0, 0, time.UTC)
		nextMonth.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

		var days []time.Time
		if nextMonth.After(yesterday) {
			days = timeshift.GetDatesBetween(currentMonth, yesterday)
		} else {
			days = timeshift.GetDatesBetween(currentMonth, nextMonth)
		}

		if len(days) <= 2 {
			continue
		}

		if !currentMonth.Before(fxOtherHistoryStart) {
			lastDay := days[len(days)-1]

			log.Debug().Msgf("Synchonizing other fx rates for %s", lastDay.Format("02.01.2006"))
			// FIXME must be last day of `currentMonth` for there fx fates
			if err := rates.syncOtherRates(lastDay); err != nil {
				log.Warn().Msg(err.Error())
			}
		}

		log.Debug().Msgf("Synchonizing main fx rates from %s to %s", days[0].Format("02.01.2006"), days[len(days)-1].Format("02.01.2006"))
		rates.syncMainRates(days)
	}
}

func (rates *CNBRatesImport) Setup() error {
	return nil
}

func (rates *CNBRatesImport) Done() <-chan interface{} {
	done := make(chan interface{})
	close(done)
	return done
}

func (rates *CNBRatesImport) Cancel() {
}

func (rates *CNBRatesImport) Work() {
	defer syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	rates.importRoundtrip()
}
