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

	"github.com/jancajthaml-openbank/cnb-rates-import/config"
	"github.com/jancajthaml-openbank/cnb-rates-import/integration"
	"github.com/jancajthaml-openbank/cnb-rates-import/metrics"
	"github.com/jancajthaml-openbank/cnb-rates-import/utils"

	localfs "github.com/jancajthaml-openbank/local-fs"
)

// Program encapsulate initialized application
type Program struct {
	cfg       config.Configuration
	interrupt chan os.Signal
	metrics   metrics.Metrics
	cnb       integration.CNBRatesImport
	cancel    context.CancelFunc
}

// Initialize application
func Initialize() Program {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.GetConfig()

	utils.SetupLogger(cfg.LogLevel)

	storage := localfs.NewStorage(cfg.RootStorage)
	metricsDaemon := metrics.NewMetrics(ctx, cfg.MetricsOutput, cfg.MetricsRefreshRate)

	cnb := integration.NewCNBRatesImport(ctx, cfg, &metricsDaemon, &storage)

	return Program{
		cfg:       cfg,
		interrupt: make(chan os.Signal, 1),
		metrics:   metricsDaemon,
		cnb:       cnb,
		cancel:    cancel,
	}
}
