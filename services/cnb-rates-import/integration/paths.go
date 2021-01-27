// Copyright (c) 2016-2021, Jan Cajthaml <jan.cajthaml@gmail.com>
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

import "time"

func GetUrlForDateMainFx(date time.Time) string {
	return "/en/financial_markets/foreign_exchange_market/exchange_rate_fixing/daily.txt?date=" + date.Format("02+01+2006")
}

func GetUrlForDateOtherFx(date time.Time) string {
	return "/en/financial_markets/foreign_exchange_market/other_currencies_fx_rates/fx_rates.txt?month=" + date.Format("1") + "&year=" + date.Format("2006")
}

func FXMainOfflineDirectory() string {
	return "raw/fx-main"
}

func FXOtherOfflineDirectory() string {
	return "raw/fx-other"
}

func FXMainOfflinePath(date time.Time) string {
	return FXMainOfflineDirectory() + "/" + date.Format("02.01.2006")
}

func FXOtherOfflinePath(date time.Time) string {
	return FXOtherOfflineDirectory() + "/" + date.Format("02.01.2006")
}
