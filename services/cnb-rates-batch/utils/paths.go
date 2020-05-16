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

import "time"

func FXMainDailyCacheDirectory() string {
	return "processed/fx-main/d"
}

func FXMainMonthlyCacheDirectory() string {
	return "processed/fx-main/m"
}

func FXMainYearlyCacheDirectory() string {
	return "processed/fx-main/y"
}

func FXOtherDailyCacheDirectory() string {
	return "processed/fx-other/d"
}

func FXOtherMonthlyCacheDirectory() string {
	return "processed/fx-other/m"
}

func FXOtherYearlyCacheDirectory() string {
	return "processed/fx-other/y"
}

func FXMainDailyCachePath(date time.Time) string {
	return FXMainDailyCacheDirectory() + "/" + date.Format("02.01.2006")
}

func FXMainMonthlyCachePath(date time.Time) string {
	return FXMainMonthlyCacheDirectory() + "/" + date.Format("01.2006")
}

func FXMainYearlyCachePath(date time.Time) string {
	return FXMainYearlyCacheDirectory() + "/" + date.Format("2006")
}

func FXOtherDailyCachePath(date time.Time) string {
	return FXOtherDailyCacheDirectory() + "/" + date.Format("02.01.2006")
}

func FXOtherMonthlyCachePath(date time.Time) string {
	return FXOtherMonthlyCacheDirectory() + "/" + date.Format("01.2006")
}

func FXOtherYearlyCachePath(date time.Time) string {
	return FXOtherYearlyCacheDirectory() + "/" + date.Format("2006")
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
