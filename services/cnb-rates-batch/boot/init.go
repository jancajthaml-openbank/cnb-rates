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

package boot

import (
	"context"
	"os"

	"github.com/jancajthaml-openbank/cnb-rates-batch/config"
	"github.com/jancajthaml-openbank/cnb-rates-batch/daemon"
	"github.com/jancajthaml-openbank/cnb-rates-batch/utils"

	localfs "github.com/jancajthaml-openbank/local-fs"
	log "github.com/sirupsen/logrus"
)

// Application encapsulate initialized application
type Application struct {
	cfg       config.Configuration
	interrupt chan os.Signal
	metrics   daemon.Metrics
	batch     daemon.Batch
	cancel    context.CancelFunc
}

// Initialize application
func Initialize() Application {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.GetConfig()

	utils.SetupLogger(cfg.LogLevel)

	log.Info(">>> Setup <<<")

	metrics := daemon.NewMetrics(ctx, cfg)

	storage := localfs.NewStorage(cfg.RootStorage)

	batch := daemon.NewBatch(ctx, cfg, &metrics, &storage)

	return Application{
		cfg:       cfg,
		interrupt: make(chan os.Signal, 1),
		metrics:   metrics,
		batch:     batch,
		cancel:    cancel,
	}
}
