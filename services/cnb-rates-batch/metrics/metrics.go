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

	localfs "github.com/jancajthaml-openbank/local-fs"
	"github.com/jancajthaml-openbank/cnb-rates-batch/utils"
	metrics "github.com/rcrowley/go-metrics"
)

// Metrics holds metrics counters
type Metrics struct {
	utils.DaemonSupport
	storage         localfs.PlaintextStorage
	refreshRate     time.Duration
	daysProcessed   metrics.Counter
	monthsProcessed metrics.Counter
}

// NewMetrics returns blank metrics holder
func NewMetrics(ctx context.Context, output string, refreshRate time.Duration) Metrics {
	return Metrics{
		DaemonSupport:   utils.NewDaemonSupport(ctx, "metrics"),
		storage:         localfs.NewPlaintextStorage(output),
		refreshRate:     refreshRate,
		daysProcessed:   metrics.NewCounter(),
		monthsProcessed: metrics.NewCounter(),
	}
}

// DayProcessed increments days processed by one
func (metrics *Metrics) DayProcessed() {
	metrics.daysProcessed.Inc(1)
}

// MonthProcessed increments months processed by one
func (metrics *Metrics) MonthProcessed() {
	metrics.monthsProcessed.Inc(1)
}

// Start handles everything needed to start metrics daemon
func (metrics Metrics) Start() {
	ticker := time.NewTicker(metrics.refreshRate)
	defer ticker.Stop()

	if err := metrics.Hydrate(); err != nil {
		log.Warn(err.Error())
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

	log.Infof("Start metrics daemon, update each %v into %v", metrics.refreshRate, metrics.storage.Root)

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
	log.Info("Stop metrics daemon")
}
