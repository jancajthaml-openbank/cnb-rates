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

package daemon

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/jancajthaml-openbank/cnb-rates-batch/config"

	localfs "github.com/jancajthaml-openbank/local-fs"
	log "github.com/sirupsen/logrus"
)

// Batch represents batch subroutine
type Batch struct {
	Support
	storage *localfs.Storage
	metrics *Metrics
}

// NewBatch returns batch fascade
func NewBatch(ctx context.Context, cfg config.Configuration, metrics *Metrics, storage *localfs.Storage) Batch {
	return Batch{
		Support: NewDaemonSupport(ctx),
		storage: storage,
		metrics: metrics,
	}
}

// WaitReady wait for batch to be ready
func (batch Batch) WaitReady(deadline time.Duration) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch x := e.(type) {
			case string:
				err = fmt.Errorf(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("unknown panic")
			}
		}
	}()

	ticker := time.NewTicker(deadline)
	select {
	case <-batch.IsReady:
		ticker.Stop()
		err = nil
		return
	case <-ticker.C:
		err = fmt.Errorf("cnb-rates-batch daemon was not ready within %v seconds", deadline)
		return
	}
}

// Start handles everything needed to start batch daemon
func (batch Batch) Start() {
	defer batch.MarkDone()

	log.Info("Start cnb-rates-batch daemon")
	batch.MarkReady()

	//batch.importRoundtrip()

	log.Info("Stopping cnb-rates-batch daemon")
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	log.Info("Stop cnb-rates-batch daemon")

	return
}
