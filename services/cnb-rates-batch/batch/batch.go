// Copyright (c) 2016-2023, Jan Cajthaml <jan.cajthaml@gmail.com>
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

	"github.com/jancajthaml-openbank/cnb-rates-batch/utils"

	localfs "github.com/jancajthaml-openbank/local-fs"
)

// Batch represents batch subroutine
type Batch struct {
	storage localfs.Storage
}

// NewBatch returns batch fascade
func NewBatch(rootStorage string) *Batch {
	storage, err := localfs.NewPlaintextStorage(rootStorage)
	if err != nil {
		log.Error().Err(err).Msg("Failed to ensure storage")
		return nil
	}
	return &Batch{
		storage: storage,
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
			log.Warn().Err(err).Msgf("failed to parse main-fx CSV data for day %s", day)
			continue
		}
		result, err := utils.ParseCSV(day, data)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to parse fx-main CSV data for day %s", day)
			continue
		}
		bytes, err := result.MarshalJSON()
		if err != nil {
			log.Warn().Err(err).Msgf("failed to marshall fx-main data for day %s", day)
			continue
		}
		err = batch.storage.WriteFile(utils.FXMainDailyCachePath(result.Date), bytes)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to write cache fx-main data for day %s", day)
			continue
		}
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
			log.Warn().Err(err).Msgf("failed to parse fx-other CSV data for day %s", day)
			continue
		}
		result, err := utils.ParseCSV(day, data)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to parse fx-other CSV data for day %s", day)
			continue
		}
		bytes, err := result.MarshalJSON()
		if err != nil {
			log.Warn().Err(err).Msgf("failed to marshall fx-other data for mo %s", day)
			continue
		}
		err = batch.storage.WriteFile(utils.FXOtherDailyCachePath(result.Date), bytes)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to write cache fail fx-other data for day %s", day)
			continue
		}
	}

	return nil
}

// Setup does nothing
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

// Work represents batch work
func (batch *Batch) Work() {
	defer syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	if batch == nil {
		return
	}

	// FIXME make cancelable
	var wg sync.WaitGroup
	wg.Add(2)

	go batch.ProcessNewFXMain(&wg)
	go batch.ProcessNewFXOther(&wg)

	wg.Wait()
}
