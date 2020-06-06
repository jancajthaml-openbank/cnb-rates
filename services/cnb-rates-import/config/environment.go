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

	"github.com/jancajthaml-openbank/cnb-rates-import/utils"
)

func loadConfFromEnv() Configuration {
	logLevel := strings.ToUpper(getEnvString("CNB_RATES_LOG_LEVEL", "DEBUG"))
	rootStorage := getEnvString("CNB_RATES_STORAGE", "/data")
	cnbGateway := getEnvString("CNB_RATES_CNB_GATEWAY", "https://www.cnb.cz")
	metricsOutput := getEnvFilename("CNB_RATES_METRICS_OUTPUT", "/tmp")
	metricsRefreshRate := getEnvDuration("CNB_RATES_METRICS_REFRESHRATE", time.Second)

	if rootStorage == "" {
		log.Fatal("missing required parameter to run")
	}

	if os.MkdirAll(rootStorage+"/rates/cnb/"+utils.FXMainOfflineDirectory(), os.ModePerm) != nil {
		log.Fatal("unable to assert daily offline directory")
	}

	if os.MkdirAll(rootStorage+"/rates/cnb/"+utils.FXOtherOfflineDirectory(), os.ModePerm) != nil {
		log.Fatal("unable to assert monthly offline directory")
	}

	return Configuration{
		RootStorage:        rootStorage + "/rates/cnb",
		CNBGateway:         cnbGateway,
		LogLevel:           logLevel,
		MetricsRefreshRate: metricsRefreshRate,
		MetricsOutput:      metricsOutput,
	}
}

func getEnvFilename(key string, fallback string) string {
	var value = strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	value = filepath.Clean(value)
	if os.MkdirAll(value, os.ModePerm) != nil {
		return fallback
	}
	return value
}
func getEnvString(key string, fallback string) string {
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
		log.Errorf("invalid value of variable %s", key)
		return fallback
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
		log.Errorf("invalid value of variable %s", key)
		return fallback
	}
	return cast
}
