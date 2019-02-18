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

package cnb

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	//decimal "gopkg.in/inf.v0"

	"github.com/jancajthaml-openbank/cnb-rates/pkg/httpclient"
	"github.com/jancajthaml-openbank/cnb-rates/pkg/utils"

	log "github.com/sirupsen/logrus"
)

type CNB struct {
	exchangeRateFixing string
	fxMainDailyPath    string
	fxOtherMonthlyPath string
	client             *httpclient.HttpClient
	Kill               chan struct{}
}

func New(cacheDirectory, gateway string) (error, *CNB) {
	if cacheDirectory == "" {
		return fmt.Errorf("persistent directory cannot be empty"), nil
	}

	absPath, err := filepath.Abs(cacheDirectory)
	if err != nil {
		return fmt.Errorf("invalid root storage path: %s", err.Error()), nil
	}

	fxMainDailyPath := absPath + "/cnb-rates/raw/fx-main/daily"
	fxOtherMonthlyPath := absPath + "/cnb-rates/raw/fx-other/monthly"

	if err := os.MkdirAll(fxMainDailyPath, os.ModePerm); err != nil {
		return fmt.Errorf("could not assert storage %s: %v", fxMainDailyPath, err.Error()), nil
	}

	if err := os.MkdirAll(fxOtherMonthlyPath, os.ModePerm); err != nil {
		return fmt.Errorf("could not assert storage %s: %v", fxOtherMonthlyPath, err.Error()), nil
	}

	return nil, &CNB{
		exchangeRateFixing: gateway + "/en/financial_markets/foreign_exchange_market",
		client:             httpclient.New(strings.HasPrefix(gateway, "https://")),
		fxMainDailyPath:    fxMainDailyPath,
		fxOtherMonthlyPath: fxOtherMonthlyPath,
	}
}

func (a *CNB) Synchronize(killChan chan struct{}) bool {
	if a == nil {
		return true
	}

	func() {
		defer func() {
			recover()
			a.Kill = make(chan struct{})
		}()

		close(a.Kill)
	}()

	doneChan := make(chan struct{})

	go func() {
		log.Info("Synchronizing rates")

		log.Debug("Synchonizing main currencies fx daily rates")
		a.SyncMainRatesDaily()

		log.Debug("Synchonizing other currencies fx monthly rates")
		a.SyncOtherRatesMonthly()

		log.Info("Synchronized")
		close(doneChan)
	}()

	select {
	case <-doneChan:
		return true
	case <-killChan:
		return false
	}
}

func (a *CNB) SyncOtherRatesMonthly() error {
	now := time.Now()

	startDate := time.Date(1991, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	dates := GetMonthsBetween(startDate, endDate)

	if len(dates) > 0 {
		var wg sync.WaitGroup
		queue := make(chan time.Time, 128)

		worker := func() {
			wg.Add(1)
			defer wg.Done()
			for {
				select {
				case <-a.Kill:
					return
				case date, ok := <-queue:
					if !ok {
						return
					}
					cachePath := getDailyMonthlyPath(a, date)

					if utils.Exists(cachePath) {
						continue
					}

					var (
						err  error
						data io.Reader
					)

					if data, err = a.fetchOtherRateForMonth(date); err != nil {
						log.Warnf("CNB cloud monthly returned error %v", err)
						continue
					}

					if !utils.UpdateFile(cachePath, data) {
						log.Warnf("cannot store cache for %s at %s", date, cachePath)
						continue
					}

					log.Debugf("processed other fx for date %s", date.Format("01.2006"))
				}
			}
		}

		for i := 0; i < 2*runtime.NumCPU(); i++ {
			go worker()
		}

		for _, date := range dates {
			if utils.Exists(getDailyCachePath(a, date)) {
				continue
			}
			queue <- date
		}

		// FIXME does not process everything :/
		// wait till queue is empty somehow
		close(queue)
		wg.Wait()
	}

	return nil
}

func (a *CNB) SyncMainRatesDaily() error {
	now := time.Now()

	startDate := time.Date(1991, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	endDate := now.AddDate(0, 0, -1)

	dates := GetDatesBetween(startDate, endDate)

	if len(dates) > 0 {
		var wg sync.WaitGroup
		queue := make(chan time.Time, 128)

		worker := func() {
			wg.Add(1)
			defer wg.Done()
			for {
				select {
				case <-a.Kill:
					return
				case date, ok := <-queue:
					if !ok {
						return
					}
					cachePath := getDailyCachePath(a, date)

					if utils.Exists(cachePath) {
						continue
					}

					var (
						err  error
						data io.Reader
					)

					if data, err = a.fetchMainRateForDay(date); err != nil {
						log.Warnf("CNB cloud daily returned error %v", err)
						continue
					}

					if !utils.UpdateFile(cachePath, data) {
						log.Warnf("cannot store cache for %s at %s", date, cachePath)
						continue
					}

					log.Debugf("processed main fx for date %s", date.Format("02.01.2006"))
				}
			}
		}

		for i := 0; i < 2*runtime.NumCPU(); i++ {
			go worker()
		}

		for _, date := range dates {
			if utils.Exists(getDailyCachePath(a, date)) {
				continue
			}
			queue <- date
		}

		// FIXME does not process everything :/
		// wait till queue is empty somehow
		close(queue)
		wg.Wait()
	}

	cachePath := getDailyCachePath(a, now)

	if utils.Exists(cachePath) {
		return nil
	}

	var (
		err  error
		data io.Reader
	)

	if data, err = a.fetchMainRateForDay(now); err != nil {
		log.Warnf("CNB cloud returned error %v", err)
		return err
	}

	if !validateRates(now, &data) {
		return fmt.Errorf("bounce %s", now.Format("02.01.2006"))
	}

	if !utils.UpdateFile(cachePath, data) {
		return fmt.Errorf("cannot store cache for %s at %s", now.Format("02.01.2006"), cachePath)
	}

	return nil
}

func (a *CNB) fetchMainRateForDay(date time.Time) (io.Reader, error) {
	url := getUrlForDateMainFx(a, date)

	var (
		err  error
		data io.Reader
	)

	err = httpclient.Retry(10, time.Second, func() (err error) {
		data, err = a.client.Get(url)
		return
	})

	return data, err
}

func (a *CNB) fetchOtherRateForMonth(date time.Time) (io.Reader, error) {
	url := getUrlForDateOtherFx(a, date)

	var (
		err  error
		data io.Reader
	)

	err = httpclient.Retry(10, time.Second, func() (err error) {
		data, err = a.client.Get(url)
		return
	})

	return data, err
}

func validateRates(date time.Time, data *io.Reader) bool {
	if data == nil {
		return false
	}

	b := make([]byte, 20)
	if _, err := (*data).Read(b); err != nil {
		return false
	}

	chunk := make([]byte, 0)

	if b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		chunk = append(chunk, b[3:]...)
	} else {
		chunk = append(chunk, b...)
	}

	*data = io.MultiReader(bytes.NewReader(chunk), *data)

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
