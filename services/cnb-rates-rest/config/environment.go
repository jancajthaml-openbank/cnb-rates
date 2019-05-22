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
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func loadConfFromEnv() Configuration {
	logLevel := strings.ToUpper(getEnvString("CNB_RATES_LOG_LEVEL", "DEBUG"))
	secrets := getEnvString("CNB_RATES_SECRETS", "")
	rootStorage := getEnvString("CNB_RATES_STORAGE", "/data")
	port := getEnvInteger("CNB_RATES_HTTP_PORT", 4011)

	if secrets == "" || rootStorage == "" {
		log.Fatal("missing required parameter to run")
	}

	return Configuration{
		RootStorage: rootStorage,
		ServerPort:  port,
		SecretsPath: secrets,
		LogLevel:    logLevel,
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
