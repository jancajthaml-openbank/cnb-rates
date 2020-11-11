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

package config

import (
	"os"
	"strings"
	"time"

	"github.com/jancajthaml-openbank/cnb-rates-batch/utils"
)

// Configuration of application
type Configuration struct {
	// CNBGateway represent cnb gateway uri
	CNBGateway string
	// RootStorage gives where to store journals
	RootStorage string
	// LogLevel ignorecase log level
	LogLevel string
	// MetricsRefreshRate represents interval in which in memory metrics should be
	// persisted to disk
	MetricsRefreshRate time.Duration
	// MetricsOutput represents output file for metrics persistence
	MetricsOutput string
}

// GetConfig loads application configuration
func GetConfig() Configuration {
	logLevel := strings.ToUpper(envString("CNB_RATES_LOG_LEVEL", "DEBUG"))
	rootStorage := envString("CNB_RATES_STORAGE", "/data")
	cnbGateway := envString("CNB_RATES_CNB_GATEWAY", "https://www.cnb.cz")
	metricsOutput := envFilename("CNB_RATES_METRICS_OUTPUT", "/tmp")
	metricsRefreshRate := envDuration("CNB_RATES_METRICS_REFRESHRATE", time.Second)

	if rootStorage == "" {
		log.Error().Msg("missing required parameter to run")
		panic("missing required parameter to run")
	}

	rootStorage = rootStorage + "/rates/cnb"

	if os.MkdirAll(rootStorage+"/"+utils.FXMainDailyCacheDirectory(), os.ModePerm) != nil {
		log.Error().Msg("unable to assert fx-main daily cache directory")
		panic("unable to assert fx-main daily cache directory")
	}

	if os.MkdirAll(rootStorage+"/"+utils.FXMainMonthlyCacheDirectory(), os.ModePerm) != nil {
		log.Error().Msg("unable to assert fx-main monthly cache directory")
		panic("unable to assert fx-main monthly cache directory")
	}

	if os.MkdirAll(rootStorage+"/"+utils.FXMainYearlyCacheDirectory(), os.ModePerm) != nil {
		log.Error().Msg("unable to assert fx-main yearly cache directory")
		panic("unable to assert fx-main yearly cache directory")
	}

	if os.MkdirAll(rootStorage+"/"+utils.FXOtherDailyCacheDirectory(), os.ModePerm) != nil {
		log.Error().Msg("unable to assert fx-other daily cache directory")
		panic("unable to assert fx-other daily cache directory")
	}

	if os.MkdirAll(rootStorage+"/"+utils.FXOtherMonthlyCacheDirectory(), os.ModePerm) != nil {
		log.Error().Msg("unable to assert fx-other monthly cache directory")
		panic("unable to assert fx-other monthly cache directory")
	}

	if os.MkdirAll(rootStorage+"/"+utils.FXOtherYearlyCacheDirectory(), os.ModePerm) != nil {
		log.Error().Msg("unable to assert fx-other yearly cache directory")
		panic("unable to assert fx-other yearly cache directory")
	}

	return Configuration{
		RootStorage:        rootStorage,
		CNBGateway:         cnbGateway,
		LogLevel:           logLevel,
		MetricsRefreshRate: metricsRefreshRate,
		MetricsOutput:      metricsOutput,
	}
}
