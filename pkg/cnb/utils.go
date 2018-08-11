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
	"time"
	//decimal "gopkg.in/inf.v0"
)

func GetDatesForYear(year string) []time.Time {
	date, err := time.Parse("2006", year)
	if err != nil {
		return nil
	}

	dates := make([]time.Time, 0)

	for i := 1; i < 13; i++ {
		startDate := time.Date(date.Year(), time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(time.Nanosecond * -1)
		for ; startDate.Before(endDate); startDate = startDate.AddDate(0, 0, 1) {
			switch startDate.Weekday() {
			case time.Sunday, time.Saturday:
				continue
			default:
				dates = append(dates, startDate)
			}
		}
	}

	return dates
}

func GetMonthsBetween(startDate time.Time, endDate time.Time) []time.Time {
	dates := make([]time.Time, 0)

	for ; startDate.Before(endDate); startDate = startDate.AddDate(0, 1, 0) {
		switch startDate.Weekday() {
		case time.Sunday, time.Saturday:
			continue
		default:
			dates = append(dates, startDate.AddDate(0, 1, 0).Add(time.Nanosecond*-1))
		}
	}
	return dates
}

func GetDatesBetween(startDate time.Time, endDate time.Time) []time.Time {
	dates := make([]time.Time, 0)

	for ; startDate.Before(endDate); startDate = startDate.AddDate(0, 0, 1) {
		switch startDate.Weekday() {
		case time.Sunday, time.Saturday:
			continue
		default:
			dates = append(dates, startDate)
		}
	}

	return dates
}

/*
func unmarshallRates(day time.Time, data []byte) (ExchangeRates, error) {
  var result ExchangeRates

  bom := []byte{0xEF, 0xBB, 0xBF}
  if bytes.HasPrefix(data, bom) {
    data = data[3:]
  }

  idx := 0

  for i := 0; ; {
    j := bytes.IndexByte(data[i:], '\n')
    if j < 0 {
      break
    }
    if idx == 0 {
      chunks := strings.Split(string(data[i:i+j]), " #")
      date, err := time.Parse("02.Jan 2006", chunks[0])
      if err != nil {
        return result, err
      }
      if date.Day() != day.Day() || date.Month() != day.Month() || date.Year() != day.Year() {
        // FIXME return error bounce
        return result, fmt.Errorf("no data for %s got %s instead", day.Format("02.01.2006"), date.Format("02.01.2006"))
      }
      result.Date = date
      id, err := strconv.Atoi(chunks[1])
      if err != nil {
        return result, err
      }
      result.Id = id
    } else {
      chunks := strings.Split(string(data[i:i+j]), "|")
      if len(chunks) == 5 && chunks[0] != "Country" {
        fromAmount, ok := new(decimal.Dec).SetString(chunks[2])
        if !ok || fromAmount == nil {
          return result, fmt.Errorf("invalid amount %s", chunks[2])
        }

        toAmount, ok := new(decimal.Dec).SetString(chunks[4])
        if !ok || toAmount == nil {
          return result, fmt.Errorf("invalid amount %s", chunks[4])
        }

        rate := ExchangeRate{
          From: Money{
            Amount:   fromAmount,
            Currency: chunks[3],
          },
          To: Money{
            Amount:   toAmount,
            Currency: "CZK",
          },
        }
        rate.Normalize()

        result.Rates = append(result.Rates, rate)
      }
    }

    j += i
    i = j + 1
    if data[j-1] == '\r' {
      j--
    }
    idx++
  }

  return result, nil
}*/

func getUrlForDateMainFx(a *CNB, date time.Time) string {
	return a.exchangeRateFixing + "/exchange_rate_fixing/daily.txt?date=" + date.Format("02.01.2006")
}

func getUrlForDateOtherFx(a *CNB, date time.Time) string {
	return a.exchangeRateFixing + "/other_currencies_fx_rates/fx_rates.txt?month=" + date.Format("01") + "&year=" + date.Format("2006")
}

func getUrlForYearMainFx(a *CNB, year time.Time) string {
	return a.exchangeRateFixing + "/exchange_rate_fixing/year.txt?year=" + year.Format("2006")
}

func getDailyCachePath(a *CNB, date time.Time) string {
	return a.fxMainDailyPath + "/" + date.Format("02.01.2006")
}

func getDailyMonthlyPath(a *CNB, date time.Time) string {
	return a.fxOtherMonthlyPath + "/" + date.Format("01.2006")
}

func getYearlyCachePath(a *CNB, year time.Time) string {
	return a.fxMainYearlyPath + "/" + year.Format("2006")
}
