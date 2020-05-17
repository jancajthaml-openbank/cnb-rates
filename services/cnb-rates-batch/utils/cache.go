// Copyright (c) 2016-2020, Jan Cajthaml <jan.cajthaml@gmail.com>
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

package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"bytes"
	"strings"
	"time"

	"github.com/jancajthaml-openbank/cnb-rates-batch/model"

	localfs "github.com/jancajthaml-openbank/local-fs"
)

func difference(a, b []string) []string {
	mb := make(map[string]interface{}, 0)
	for _, x := range b {
		mb[x] = nil
	}
	ab := make([]string, 0)
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}

func ParseCSV(filename string, data []byte) (model.ExchangeFixing, error) {
	result := model.ExchangeFixing{}
	result.Rates = make([]model.Exchange, 0)

	rc := bytes.NewReader(data)

	r := csv.NewReader(rc)
	r.Comment = '#'

	dateLine, err := r.Read()
	if err != nil {
		return result, err
	}
	parts := strings.Split(dateLine[0], " #")
	date, err := time.Parse("02 Jan 2006", parts[0])
	if err != nil {
		return result, err
	}

	if filename != date.Format("02.01.2006") {
		date, _ = time.Parse("02.01.2006", filename)
		result.Date = date
		return result, nil
	}

	result.Date = date

	r.Comma = '|'
	r.FieldsPerRecord = 5

	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return result, fmt.Errorf("line %+v error: %+v", rec, err)
		}

		if rec[0] == "Country" {
			continue
		}
		entry := model.Exchange{}
		entry.UnmarshalText(rec)
		result.Rates = append(result.Rates, entry)
	}

	return result, nil
}

func GetFXMainUnprocessedFiles(storage *localfs.PlaintextStorage) ([]string, error) {
	raw, err := storage.ListDirectory(FXMainOfflineDirectory(), true)
	if err != nil {
		return nil, err
	}
	processed, err := storage.ListDirectory(FXMainDailyCacheDirectory(), true)
	if err != nil {
		return nil, err
	}
	return difference(raw, processed), nil
}

func GetFXOtherUnprocessedFiles(storage *localfs.PlaintextStorage) ([]string, error) {
	raw, err := storage.ListDirectory(FXOtherOfflineDirectory(), true)
	if err != nil {
		return nil, err
	}
	processed, err := storage.ListDirectory(FXOtherDailyCacheDirectory(), true)
	if err != nil {
		return nil, err
	}
	return difference(raw, processed), nil
}
