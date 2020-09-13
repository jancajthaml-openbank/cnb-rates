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

package boot

import (
	"context"
	"os"

	"github.com/jancajthaml-openbank/cnb-rates-rest/api"
	"github.com/jancajthaml-openbank/cnb-rates-rest/config"
	"github.com/jancajthaml-openbank/cnb-rates-rest/logging"
	"github.com/jancajthaml-openbank/cnb-rates-rest/metrics"
	"github.com/jancajthaml-openbank/cnb-rates-rest/utils"

	localfs "github.com/jancajthaml-openbank/local-fs"
)

// Program encapsulate initialized application
type Program struct {
	interrupt chan os.Signal
	cfg       config.Configuration
	daemons   []utils.Daemon
	cancel    context.CancelFunc
}

// Initialize application
func Initialize() Program {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.GetConfig()

	logging.SetupLogger(cfg.LogLevel)

	storage := localfs.NewPlaintextStorage(
		cfg.RootStorage,
	)
	metricsDaemon := metrics.NewMetrics(
		ctx,
		cfg.MetricsOutput,
		cfg.MetricsRefreshRate,
	)
	restDaemon := api.NewServer(
		ctx,
		cfg.ServerPort,
		cfg.ServerCert,
		cfg.ServerKey,
		&storage,
	)

	var daemons = make([]utils.Daemon, 0)
	daemons = append(daemons, metricsDaemon)
	daemons = append(daemons, restDaemon)

	return Program{
		interrupt: make(chan os.Signal, 1),
		cfg:       cfg,
		daemons:   daemons,
		cancel:    cancel,
	}
}
