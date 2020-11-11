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

	"github.com/jancajthaml-openbank/cnb-rates-batch/batch"
	"github.com/jancajthaml-openbank/cnb-rates-batch/config"
	"github.com/jancajthaml-openbank/cnb-rates-batch/logging"
	"github.com/jancajthaml-openbank/cnb-rates-batch/metrics"
	"github.com/jancajthaml-openbank/cnb-rates-batch/utils"
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

	metricsDaemon := metrics.NewMetrics(
		ctx,
		cfg.MetricsOutput,
		cfg.MetricsRefreshRate,
	)
	batchDaemon := batch.NewBatch(
		ctx,
		cfg.RootStorage,
		metricsDaemon,
	)

	var daemons = make([]utils.Daemon, 0)
	daemons = append(daemons, metricsDaemon)
	daemons = append(daemons, batchDaemon)

	return Program{
		interrupt: make(chan os.Signal, 1),
		cfg:       cfg,
		daemons:   daemons,
		cancel:    cancel,
	}
}
