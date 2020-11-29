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
	"strings"
	"time"
)

// Configuration of application
type Configuration struct {
	// CNBGateway represent cnb gateway uri
	CNBGateway string
	// RootStorage gives where to store journals
	RootStorage string
	// LogLevel ignorecase log level
	LogLevel string
	// MetricsContinuous determines if metrics should start from last state
	MetricsContinuous bool
	// MetricsRefreshRate represents interval in which in memory metrics should be
	// persisted to disk
	MetricsRefreshRate time.Duration
	// MetricsOutput represents output file for metrics persistence
	MetricsOutput string
}

// LoadConfig loads application configuration
func LoadConfig() Configuration {
	return Configuration{
		RootStorage:        envString("CNB_RATES_STORAGE", "/data"),
		CNBGateway:         envString("CNB_RATES_CNB_GATEWAY", "https://www.cnb.cz"),
		LogLevel:           strings.ToUpper(envString("CNB_RATES_LOG_LEVEL", "DEBUG")),
		MetricsContinuous:  envBoolean("CNB_RATES_METRICS_CONTINUOUS", true),
		MetricsRefreshRate: envDuration("CNB_RATES_METRICS_REFRESHRATE", time.Second),
		MetricsOutput:      envFilename("CNB_RATES_METRICS_OUTPUT", "/tmp/cnb-rates-batch-metrics"),
	}
}
