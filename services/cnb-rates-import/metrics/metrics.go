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

package metrics

import (
	"context"
	"time"

	"github.com/jancajthaml-openbank/cnb-rates-import/utils"
	localfs "github.com/jancajthaml-openbank/local-fs"
	metrics "github.com/rcrowley/go-metrics"
)

// Metrics holds metrics counters
type Metrics struct {
	utils.DaemonSupport
	storage        localfs.PlaintextStorage
	refreshRate    time.Duration
	daysImported   metrics.Counter
	monthsImported metrics.Counter
	gatewayLatency metrics.Timer
	importLatency  metrics.Timer
}

// NewMetrics returns blank metrics holder
func NewMetrics(ctx context.Context, output string, refreshRate time.Duration) Metrics {
	return Metrics{
		DaemonSupport:  utils.NewDaemonSupport(ctx, "metrics"),
		storage:        localfs.NewPlaintextStorage(output),
		refreshRate:    refreshRate,
		daysImported:   metrics.NewCounter(),
		monthsImported: metrics.NewCounter(),
		gatewayLatency: metrics.NewTimer(),
		importLatency:  metrics.NewTimer(),
	}
}

// TimeGatewayLatency measure execution of gateway sync
func (metrics *Metrics) TimeGatewayLatency(f func()) {
	metrics.gatewayLatency.Time(f)
}

// TimeImportLatency measure execution of import
func (metrics *Metrics) TimeImportLatency(f func()) {
	metrics.importLatency.Time(f)
}

// DayImported increments days imported by one
func (metrics *Metrics) DayImported() {
	metrics.daysImported.Inc(1)
}

// MonthImported increments months imported by one
func (metrics *Metrics) MonthImported() {
	metrics.monthsImported.Inc(1)
}

// Start handles everything needed to start metrics daemon
func (metrics Metrics) Start() {
	ticker := time.NewTicker(metrics.refreshRate)
	defer ticker.Stop()

	if err := metrics.Hydrate(); err != nil {
		log.Warn().Msg(err.Error())
	}

	metrics.Persist()
	metrics.MarkReady()

	select {
	case <-metrics.CanStart:
		break
	case <-metrics.Done():
		metrics.MarkDone()
		return
	}

	log.Info().Msgf("Start metrics daemon, update each %v into %v", metrics.refreshRate, metrics.storage.Root)

	go func() {
		for {
			select {
			case <-metrics.Done():
				metrics.Persist()
				metrics.MarkDone()
				return
			case <-ticker.C:
				metrics.Persist()
			}
		}
	}()

	metrics.WaitStop()
	log.Info().Msg("Stop metrics daemon")
}
