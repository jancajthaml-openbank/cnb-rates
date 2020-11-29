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

package batch

import (
	"sync"
	"syscall"

	"github.com/jancajthaml-openbank/cnb-rates-batch/metrics"
	"github.com/jancajthaml-openbank/cnb-rates-batch/utils"

	localfs "github.com/jancajthaml-openbank/local-fs"
)

// Batch represents batch subroutine
type Batch struct {
	storage localfs.Storage
	metrics *metrics.Metrics
}

// NewBatch returns batch fascade
func NewBatch(rootStorage string, metrics *metrics.Metrics) *Batch {
	storage, err := localfs.NewPlaintextStorage(rootStorage)
	if err != nil {
		log.Error().Msgf("Failed to ensure storage %+v", err)
		return nil
	}
	return &Batch{
		storage:       storage,
		metrics:       metrics,
	}
}

func (batch *Batch) ProcessNewFXMain(wg *sync.WaitGroup) error {
	if batch == nil {
		return nil
	}

	defer wg.Done()
	log.Info().Msg("Processing new fx-main rates")

	days, err := utils.GetFXMainUnprocessedFiles(batch.storage)
	if err != nil {
		return err
	}

	cachePath := utils.FXMainOfflineDirectory() + "/"

	for _, day := range days {
		log.Info().Msgf("Processing new fx-main for %s", day)
		data, err := batch.storage.ReadFileFully(cachePath + day)
		if err != nil {
			log.Warn().Msgf("error parse main-fx CSV data for day %s, %+v\n", day, err)
			continue
		}
		result, err := utils.ParseCSV(day, data)
		if err != nil {
			log.Warn().Msgf("error parse fx-main CSV data for day %s, %+v\n", day, err)
			continue
		}
		bytes, err := result.MarshalJSON()
		if err != nil {
			log.Warn().Msgf("error marshall fx-main data for day %s, %+v\n", day, err)
			continue
		}
		err = batch.storage.WriteFile(utils.FXMainDailyCachePath(result.Date), bytes)
		if err != nil {
			log.Warn().Msgf("error write cache fail fx-main data for day %s, %+v\n", day, err)
			continue
		}
		batch.metrics.DayProcessed()
	}

	return nil
}

func (batch *Batch) ProcessNewFXOther(wg *sync.WaitGroup) error {
	if batch == nil {
		return nil
	}

	defer wg.Done()

	log.Info().Msg("Processing new fx-other rates")

	days, err := utils.GetFXOtherUnprocessedFiles(batch.storage)
	if err != nil {
		return err
	}

	cachePath := utils.FXOtherOfflineDirectory() + "/"
	for _, day := range days {
		log.Info().Msgf("Processing new fx-other for %s", day)
		data, err := batch.storage.ReadFileFully(cachePath + day)
		if err != nil {
			log.Warn().Msgf("error parse fx-other CSV data for day %s, %+v\n", day, err)
			continue
		}
		result, err := utils.ParseCSV(day, data)
		if err != nil {
			log.Warn().Msgf("error parse fx-other CSV data for day %s, %+v\n", day, err)
			continue
		}
		bytes, err := result.MarshalJSON()
		if err != nil {
			log.Warn().Msgf("error marshall fx-other data for mo %s, %+v\n", day, err)
			continue
		}
		err = batch.storage.WriteFile(utils.FXOtherDailyCachePath(result.Date), bytes)
		if err != nil {
			log.Warn().Msgf("error write cache fail fx-other data for day %s, %+v\n", day, err)
			continue
		}
		batch.metrics.DayProcessed()
	}

	return nil
}


// Setup hydrates metrics from storage
func (batch *Batch) Setup() error {
	return nil
}

// Done returns always finished
func (batch *Batch) Done() <-chan interface{} {
	done := make(chan interface{})
	close(done)
	return done
}

// Cancel does nothing
func (batch *Batch) Cancel() {
}

// Work represents metrics worker work
func (batch *Batch) Work() {
	defer syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	if batch == nil {
		return
	}

	// FIXME make cancelable
	var wg sync.WaitGroup

	wg.Add(1)
	go batch.ProcessNewFXMain(&wg)

	wg.Add(1)
	go batch.ProcessNewFXOther(&wg)

	wg.Wait()


}
