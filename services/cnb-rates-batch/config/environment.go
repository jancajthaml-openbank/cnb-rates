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

package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jancajthaml-openbank/cnb-rates-batch/utils"
)

func loadConfFromEnv() Configuration {
	logOutput := getEnvString("CNB_RATES_LOG", "")
	logLevel := strings.ToUpper(getEnvString("CNB_RATES_LOG_LEVEL", "DEBUG"))
	rootStorage := getEnvString("CNB_RATES_STORAGE", "/data")
	cnbGateway := getEnvString("CNB_RATES_CNB_GATEWAY", "https://www.cnb.cz")
	metricsOutput := getEnvString("CNB_RATES_METRICS_OUTPUT", "")
	metricsRefreshRate := getEnvDuration("CNB_RATES_METRICS_REFRESHRATE", time.Second)

	if rootStorage == "" {
		log.Fatal("missing required parameter to run")
	}

	rootStorage = rootStorage + "/rates/cnb"

	if metricsOutput != "" && os.MkdirAll(filepath.Dir(metricsOutput), os.ModePerm) != nil {
		log.Fatal("unable to assert metrics output")
	}

	if os.MkdirAll(rootStorage+"/"+utils.FXMainDailyCacheDirectory(), os.ModePerm) != nil {
		log.Fatal("unable to assert fx-main daily cache directory")
	}

	if os.MkdirAll(rootStorage+"/"+utils.FXMainMonthlyCacheDirectory(), os.ModePerm) != nil {
		log.Fatal("unable to assert fx-main monthly cache directory")
	}

	if os.MkdirAll(rootStorage+"/"+utils.FXMainYearlyCacheDirectory(), os.ModePerm) != nil {
		log.Fatal("unable to assert fx-main yearly cache directory")
	}

	if os.MkdirAll(rootStorage+"/"+utils.FXOtherDailyCacheDirectory(), os.ModePerm) != nil {
		log.Fatal("unable to assert fx-other daily cache directory")
	}

	if os.MkdirAll(rootStorage+"/"+utils.FXOtherMonthlyCacheDirectory(), os.ModePerm) != nil {
		log.Fatal("unable to assert fx-other monthly cache directory")
	}

	if os.MkdirAll(rootStorage+"/"+utils.FXOtherYearlyCacheDirectory(), os.ModePerm) != nil {
		log.Fatal("unable to assert fx-other yearly cache directory")
	}

	return Configuration{
		RootStorage:        rootStorage,
		CNBGateway:         cnbGateway,
		LogOutput:          logOutput,
		LogLevel:           logLevel,
		MetricsRefreshRate: metricsRefreshRate,
		MetricsOutput:      metricsOutput,
	}
}

func getEnvString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInteger(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	cast, err := strconv.Atoi(value)
	if err != nil {
		log.Panicf("invalid value of variable %s", key)
	}
	return cast
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	cast, err := time.ParseDuration(value)
	if err != nil {
		log.Panicf("invalid value of variable %s", key)
	}
	return cast
}
