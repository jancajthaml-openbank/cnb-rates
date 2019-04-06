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

package daemon

import (
	"context"
	"fmt"
	"sync"
	"syscall"
	"time"

	"github.com/jancajthaml-openbank/cnb-rates-batch/config"
	"github.com/jancajthaml-openbank/cnb-rates-batch/utils"

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

/*
func (batch Batch) ProcessNewMonthly() {
	log.Infof("Processing new Monthly rates")

	monthlyUnprocessed, _ := utils.GetMonthlyUnprocessedFiles(batch.storage)
	log.Info("Monthly")
	log.Infof("Unprocessed: %+v", monthlyUnprocessed)

}
*/

func (batch Batch) ProcessNewFXMain(wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Info("Processing new fx-main rates")

	days, err := utils.GetFXMainUnprocessedFiles(batch.storage)
	if err != nil {
		return err
	}

	cachePath := utils.FXMainOfflineDirectory() + "/"

	for _, day := range days {
		if batch.ctx.Err() != nil {
			return nil
		}

		log.Infof("Processing new fx-main for %s", day)

		reader, err := batch.storage.GetFileReader(cachePath + day)
		if err != nil {
			log.Warnf("error parse main-fx CSV data for day %s, %+v\n", day, err)
			continue
		}

		result, err := utils.ParseCSV(day, reader)
		if err != nil {
			log.Warnf("error parse fx-main CSV data for day %s, %+v\n", day, err)
			continue
		}

		bytes, err := result.MarshalJSON()
		if err != nil {
			log.Warnf("error marshall fx-main data for day %s, %+v\n", day, err)
			continue
		}

		err = batch.storage.WriteFile(utils.FXMainDailyCachePath(result.Date), bytes)
		if err != nil {
			log.Warnf("error write cache fail fx-main data for day %s, %+v\n", day, err)
			continue
		}
	}

	return nil
}

func (batch Batch) ProcessNewFXOther(wg *sync.WaitGroup) error {
	defer wg.Done()

	log.Info("Processing new fx-other rates")

	days, err := utils.GetFXOtherUnprocessedFiles(batch.storage)
	if err != nil {
		return err
	}

	cachePath := utils.FXOtherOfflineDirectory() + "/"
	for _, day := range days {
		if batch.ctx.Err() != nil {
			return nil
		}

		log.Infof("Processing new fx-other for %s", day)

		reader, err := batch.storage.GetFileReader(cachePath + day)
		if err != nil {
			log.Warnf("error parse fx-other CSV data for day %s, %+v\n", day, err)
			continue
		}

		result, err := utils.ParseCSV(day, reader)
		if err != nil {
			log.Warnf("error parse fx-other CSV data for day %s, %+v\n", day, err)
			continue
		}

		bytes, err := result.MarshalJSON()
		if err != nil {
			log.Warnf("error marshall fx-other data for mo %s, %+v\n", day, err)
			continue
		}

		err = batch.storage.WriteFile(utils.FXOtherDailyCachePath(result.Date), bytes)
		if err != nil {
			log.Warnf("error write cache fail fx-other data for day %s, %+v\n", day, err)
			continue
		}
	}

	return nil
}

func (batch Batch) ProcessNewFX() {
	var wg sync.WaitGroup

	wg.Add(1)
	go batch.ProcessNewFXMain(&wg)

	wg.Add(1)
	go batch.ProcessNewFXOther(&wg)

	wg.Wait()
}

// Start handles everything needed to start batch daemon
func (batch Batch) Start() {
	defer batch.MarkDone()

	batch.MarkReady()

	select {
	case <-batch.canStart:
		break
	case <-batch.Done():
		return
	}

	log.Info("Start cnb-rates-batch daemon")

	batch.ProcessNewFX()

	log.Info("Stopping cnb-rates-batch daemon")
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	log.Info("Stop cnb-rates-batch daemon")

	return
}
