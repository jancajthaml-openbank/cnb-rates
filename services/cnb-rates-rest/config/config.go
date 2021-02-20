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

package config

import (
	"github.com/jancajthaml-openbank/cnb-rates-rest/support/env"
	"strings"
)

// Configuration of application
type Configuration struct {
	// RootStorage gives where to store journals
	RootStorage string
	// ServerPort is port which server is bound to
	ServerPort int
	// ServerKey path to server tls key file
	ServerKey string
	// ServerCert path to server tls cert file
	ServerCert string
	// LogLevel ignorecase log level
	LogLevel string
}

// LoadConfig loads application configuration
func LoadConfig() Configuration {
	return Configuration{
		RootStorage: env.String("CNB_RATES_STORAGE", "/data"),
		ServerPort:  env.Int("CNB_RATES_HTTP_PORT", 4011),
		ServerKey:   env.String("CNB_RATES_SERVER_KEY", ""),
		ServerCert:  env.String("CNB_RATES_SERVER_CERT", ""),
		LogLevel:    strings.ToUpper(env.String("CNB_RATES_LOG_LEVEL", "DEBUG")),
	}
}
