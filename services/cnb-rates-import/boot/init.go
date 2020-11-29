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
	"time"
	"os"

	"github.com/jancajthaml-openbank/cnb-rates-import/config"
	"github.com/jancajthaml-openbank/cnb-rates-import/metrics"
	"github.com/jancajthaml-openbank/cnb-rates-import/support/concurrent"
	"github.com/jancajthaml-openbank/cnb-rates-import/support/logging"
	"github.com/jancajthaml-openbank/cnb-rates-import/integration"
)

// Program encapsulate program
type Program struct {
	interrupt chan os.Signal
	cfg       config.Configuration
	daemons   []concurrent.Daemon
}

// Register daemon into program
func (prog *Program) Register(daemon concurrent.Daemon) {
	if prog == nil || daemon == nil {
		return
	}
	prog.daemons = append(prog.daemons, daemon)
}

// NewProgram returns new program
func NewProgram() Program {
	return Program{
		interrupt: make(chan os.Signal, 1),
		cfg:       config.LoadConfig(),
		daemons:   make([]concurrent.Daemon, 0),
	}
}

// Setup setups program
func (prog *Program) Setup() {
	if prog == nil {
		return
	}

	logging.SetupLogger(prog.cfg.LogLevel)

	metricsWorker := metrics.NewMetrics(
		prog.cfg.MetricsOutput,
		prog.cfg.MetricsContinuous,
	)

	cnbWorker := integration.NewCNBRatesImport(
		prog.cfg.CNBGateway,
		prog.cfg.RootStorage,
		metricsWorker,
	)

	prog.Register(concurrent.NewScheduledDaemon(
		"metrics",
		metricsWorker,
		prog.cfg.MetricsRefreshRate,
	))

	prog.Register(concurrent.NewScheduledDaemon(
		"cnb",
		cnbWorker,
		time.Second,
	))

}
