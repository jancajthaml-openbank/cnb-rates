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

package main

import (
	"os"
	"os/signal"
	"syscall"

	"bufio"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/jancajthaml-openbank/cnb-rates/pkg/cnb"
	"github.com/jancajthaml-openbank/cnb-rates/pkg/utils"
)

func init() {
	viper.SetEnvPrefix("CNB_RATES")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("log.level", "DEBUG")
	viper.SetDefault("storage", "/data")
	viper.SetDefault("sync.rate", "1m")
	viper.SetDefault("gateway", "https://www.cnb.cz")

	log.SetFormatter(new(utils.LogFormat))
}

func validParams(params utils.RunParams) bool {
	if os.MkdirAll(params.RootStorage, os.ModePerm) != nil {
		log.Error("unable to assert storage directory")
		return false
	}

	// FIXME validate sync rate duration

	// FIXME validate that gateway is proper url

	return true
}

func loadParams() utils.RunParams {
	return utils.RunParams{
		RootStorage: viper.GetString("storage"),
		Log:         viper.GetString("log"),
		LogLevel:    viper.GetString("log.level"),
		SyncRate:    viper.GetDuration("sync.rate"),
		Gateway:     viper.GetString("gateway"),
	}
}

func main() {
	log.Print(">>> Setup <<<")

	params := loadParams()
	if !validParams(params) {
		return
	}

	if params.Log == "" {
		log.SetOutput(os.Stdout)
	} else if file, err := os.OpenFile(params.Log, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644); err == nil {
		defer file.Close()
		log.SetOutput(bufio.NewWriter(file))
	} else {
		log.SetOutput(os.Stdout)
		log.Warnf("Unable to create %s: %v", params.Log, err)
	}

	if level, err := log.ParseLevel(params.LogLevel); err == nil {
		log.Printf("Log level set to %v", strings.ToUpper(params.LogLevel))
		log.SetLevel(level)
	} else {
		log.Warnf("Invalid log level %v, using level WARN", params.LogLevel)
		log.SetLevel(log.WarnLevel)
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	log.Print(">>> Starting <<<")

	var wg sync.WaitGroup

	terminationChan := make(chan struct{})
	wg.Add(1)
	go cnb.SynchronizeRates(&wg, terminationChan, params)

	log.Print(">>> Started <<<")

	<-exitSignal

	log.Print(">>> Terminating <<<")
	close(terminationChan)
	wg.Wait()

	log.Print(">>> Terminated <<<")
}
