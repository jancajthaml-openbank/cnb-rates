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

package model

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	money "gopkg.in/inf.v0"
)

type ExchangeFixing struct {
	Rates []Exchange
	Date  time.Time
}

type Exchange struct {
	Currency string
	Rate     *money.Dec
}

func (entity *Exchange) UnmarshalText(data []string) error {
	if entity == nil {
		return fmt.Errorf("cannot unmarshall to nil pointer")
	}
	if len(data) < 5 {
		return fmt.Errorf("invalid data")
	}
	amount, ok := new(money.Dec).SetString(data[2])
	if !ok {
		return fmt.Errorf("invalid amount %s", data[2])
	}
	rate, ok := new(money.Dec).SetString(data[4])
	if !ok {
		return fmt.Errorf("invalid rate %s", data[4])
	}
	if len(data[3]) != 3 ||
		!((data[3][0] >= 'A' && data[3][0] <= 'Z') &&
			(data[3][1] >= 'A' && data[3][1] <= 'Z') &&
			(data[3][2] >= 'A' && data[3][2] <= 'Z')) {
		return fmt.Errorf("invalid currency %s", data[3])
	}
	entity.Currency = data[3]
	entity.Rate = new(money.Dec).QuoRound(rate, amount, 35, money.RoundHalfEven)
	return nil
}

// MarshalJSON serializes ExchangeFixing as json
func (entity ExchangeFixing) MarshalJSON() ([]byte, error) {
	pivot := len(entity.Rates) - 1
	var buffer bytes.Buffer
	for idx, rate := range entity.Rates {
		if idx == pivot {
			buffer.WriteString("\"" + rate.Currency + "\":\"" + strings.TrimRight(strings.TrimRight(rate.Rate.String(), "0"), ".") + "\"")
		} else {
			buffer.WriteString("\"" + rate.Currency + "\":\"" + strings.TrimRight(strings.TrimRight(rate.Rate.String(), "0"), ".") + "\",")
		}
	}
	return []byte("{" + buffer.String() + "}"), nil
}
