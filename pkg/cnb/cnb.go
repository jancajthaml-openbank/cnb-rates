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
)

type CNB struct {
	exchangeRateFixing string
	fxMainDailyPath    string
	fxMainYearlyPath   string
	fxOtherMonthlyPath string
	client             *httpclient.HttpClient
	Kill chan struct{}
}

func (a *CNB) Stop() {
	if a == nil {
		return
	}

	defer func() {
		recover()
		a.Kill = make(chan struct{})
	}()

	close(a.Kill)
}

func New(cacheDirectory, gateway string) (error, *CNB) {
	if cacheDirectory == "" {
		return fmt.Errorf("persistent directory cannot be empty"), nil
	}

	absPath, err := filepath.Abs(cacheDirectory)
	if err != nil {
		return fmt.Errorf("invalid root storage path: %s", err.Error()), nil
	}

	cacheDirectory = absPath

	fxMainDailyPath := cacheDirectory + "/raw/daily"
	fxMainYearlyPath := cacheDirectory + "/raw/yearly"
	fxOtherMonthlyPath := cacheDirectory + "/raw/monthly"

	if err := os.MkdirAll(fxMainDailyPath, os.ModePerm); err != nil {
		return fmt.Errorf("could not assert storage: %s", err.Error()), nil
	}

	if err := os.MkdirAll(fxOtherMonthlyPath, os.ModePerm); err != nil {
		return fmt.Errorf("could not assert storage: %s", err.Error()), nil
	}

	if err := os.MkdirAll(fxMainYearlyPath, os.ModePerm); err != nil {
		return fmt.Errorf("could not assert storage: %s", err.Error()), nil
	}

	return nil, &CNB{
		exchangeRateFixing: gateway + "/en/financial_markets/foreign_exchange_market",
		client:             httpclient.New(strings.HasPrefix(gateway, "https://")),
		fxMainDailyPath:    fxMainDailyPath,
		fxMainYearlyPath:   fxMainYearlyPath,
		fxOtherMonthlyPath: fxOtherMonthlyPath,
	}
}

func (a *CNB) SyncMainRatesYearly() error {
	consolidated := make(map[string]interface{})
	for _, item := range utils.ListDirectory(a.fxMainYearlyPath) {
		consolidated[item] = nil
	}

	all := utils.ListDirectory(a.fxMainDailyPath)

	byYear := make(map[string][]string)

	for _, item := range all {
		year := item[6:]
		if _, ok := consolidated[year]; !ok {
			byYear[year] = append(byYear[year], item)
		}
	}

	years := make([]time.Time, 0)

	for year, items := range byYear {
		dates := GetDatesForYear(year)

		allDatesPresent := func(x []time.Time, y []string) bool {
			if len(x) != len(y) {
				return false
			}
			diff := make(map[string]int, len(x))
			for _, _x := range x {
				diff[_x.Format("02.01.2006")]++
			}
			for _, _y := range y {
				if _, ok := diff[_y]; !ok {
					return false
				}
				diff[_y] -= 1
				if diff[_y] == 0 {
					delete(diff, _y)
				}
			}
			if len(diff) == 0 {
				return true
			}
			fmt.Println(diff)
			return false
		}

		if allDatesPresent(dates, items) {
			date, err := time.Parse("2006", year)
			if err != nil {
				return fmt.Errorf("Invalid file found at " + a.fxMainYearlyPath)
			}
			years = append(years, date)
		}
	}

	if len(years) > 0 {
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
					cachePath := getYearlyCachePath(a, date)

					if utils.Exists(cachePath) {
						continue
					}

					var (
						err  error
						data io.Reader
					)

					if data, err = a.fetchMainRateForYear(date); err != nil {
						fmt.Println("CNB cloud yearly returned error", err)
						continue
					}
					if !utils.UpdateFile(cachePath, data) {
						fmt.Printf("cannot store cache for %s at %s\n", date, cachePath)
					}

					fmt.Println("processed year", date.Format("2006"))
				}
			}
		}

		for i := 0; i < 2*runtime.NumCPU(); i++ {
			go worker()
		}

		for _, year := range years {
			if utils.Exists(getYearlyCachePath(a, year)) {
				continue
			}
			queue <- year
		}

		// FIXME does not process everything :/
		// wait till queue is empty somehow
		close(queue)
		wg.Wait()
	}

	return nil
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
						fmt.Println("CNB cloud monthly returned error", err)
						continue
					}

					if !utils.UpdateFile(cachePath, data) {
						fmt.Printf("cannot store cache for %s at %s\n", date, cachePath)
					}

					fmt.Println("processed date", date.Format("01.2006"))
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
						fmt.Println("CNB cloud daily returned error", err)
						continue
					}

					if !utils.UpdateFile(cachePath, data) {
						fmt.Printf("cannot store cache for %s at %s\n", date, cachePath)
					}

					fmt.Println("processed date", date.Format("02.01.2006"))
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
		fmt.Println("CNB cloud returned error", err)
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

/*
func (a *CNB) GetRateForDay(day time.Time) (ExchangeRates, error) {
	var (
		rates ExchangeRates
		data  []byte
		err   error
		ok    bool
	)

	cachePath := getDailyCachePatch(a, day)

	if data, ok = utils.ReadFileFully(cachePath); ok {
		return unmarshallRates(day, data)
	}

	if data, err = a.fetchRateForDay(day); err != nil {
		return rates, err
	}

	rates, err = unmarshallRates(day, data)
	if err != nil {
		return rates, err
	}

	if !utils.UpdateFile(cachePath, data) {
		return rates, fmt.Errorf("cannot store cache for %s at %s", day.Format("02.01.2006"), cachePath)
	}

	return rates, err
}*/

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

func (a *CNB) fetchMainRateForYear(date time.Time) (io.Reader, error) {
	url := getUrlForYearMainFx(a, date)

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
