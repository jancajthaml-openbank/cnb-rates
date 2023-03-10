// Copyright (c) 2016-2023, Jan Cajthaml <jan.cajthaml@gmail.com>
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
	"github.com/jancajthaml-openbank/cnb-rates-import/support/env"
	"strings"
)

// Configuration of application
type Configuration struct {
	// CNBGateway represent cnb gateway uri
	CNBGateway string
	// RootStorage gives where to store journals
	RootStorage string
	// LogLevel ignorecase log level
	LogLevel string
}

// LoadConfig loads application configuration
func LoadConfig() Configuration {
	return Configuration{
		RootStorage: env.String("CNB_RATES_STORAGE", "/data") + "/rates/cnb",
		CNBGateway:  env.String("CNB_RATES_CNB_GATEWAY", "https://www.cnb.cz"),
		LogLevel:    strings.ToUpper(env.String("CNB_RATES_LOG_LEVEL", "DEBUG")),
	}
}
