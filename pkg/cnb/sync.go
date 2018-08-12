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

package cnb

import (
	"sync"
	"time"

	"github.com/jancajthaml-openbank/cnb-rates/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func do(client *CNB) {
	log.Debug("synchonizing main currencies fx daily rates")

	client.SyncMainRatesDaily()

	log.Debug("synchonizing other currencies fx monthly rates")
	client.SyncOtherRatesMonthly()

	log.Debug("synchonizing main currencies fx yearly rates")
	client.SyncMainRatesYearly()

	return
}

// SynchronizeRates synchronizes rate from CNB cloud
func SynchronizeRates(wg *sync.WaitGroup, terminationChan chan struct{}, params utils.RunParams) {
	defer wg.Done()

	ticker := time.NewTicker(params.SyncRate)
	defer ticker.Stop()

	err, client := New(params.RootStorage)
	if err != nil {
		log.Errorf("Unable to create new client %v", err)
		return
	}

	log.Debugf("Synchronizing each %v into %v", params.SyncRate, params.RootStorage)

	for {
		select {
		case <-ticker.C:
			do(client)
		case <-terminationChan:
			return
		}
	}
}
