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
	log.Info("Synchronizing rates")

	log.Debug("Synchonizing main currencies fx daily rates")
	client.SyncMainRatesDaily()

	log.Debug("Synchonizing other currencies fx monthly rates")
	client.SyncOtherRatesMonthly()

	log.Info("Synchronized")
	return
}

// SynchronizeRates synchronizes rate from CNB cloud
func SynchronizeRates(wg *sync.WaitGroup, terminationChan chan struct{}, params utils.RunParams) {

	// //panic: sync: WaitGroup is reused before previous Wait has returned

	//		Aug 13 07:49:45.155687 cnb-rates linux-amd64[166]: panic: sync: WaitGroup is reused before previous Wait has returned
	//		Aug 13 07:49:45.155687 cnb-rates linux-amd64[166]: goroutine 4 [running]:
	//		Aug 13 07:49:45.155687 cnb-rates linux-amd64[166]: sync.(*WaitGroup).Wait(0xc4203f0580)
	//		Aug 13 07:49:45.155687 cnb-rates linux-amd64[166]:         /usr/local/go/src/sync/waitgroup.go:131 +0xbb
	//		Aug 13 07:49:45.155687 cnb-rates linux-amd64[166]: github.com/jancajthaml-openbank/cnb-rates/pkg/cnb.(*CNB).SyncOtherRatesMonthly(0xc42016c190, 0x1, 0x1)
	//		Aug 13 07:49:45.156648 cnb-rates linux-amd64[166]:         /go/src/github.com/jancajthaml-openbank/cnb-rates/pkg/cnb/cnb.go:251 +0x3cb
	//		Aug 13 07:49:45.156648 cnb-rates linux-amd64[166]: github.com/jancajthaml-openbank/cnb-rates/pkg/cnb.do(0xc42016c190)
	//		Aug 13 07:49:45.156648 cnb-rates linux-amd64[166]:         /go/src/github.com/jancajthaml-openbank/cnb-rates/pkg/cnb/sync.go:31 +0xbd
	//		Aug 13 07:49:45.156648 cnb-rates linux-amd64[166]: github.com/jancajthaml-openbank/cnb-rates/pkg/cnb.SynchronizeRates(0xc420026030, 0xc420178000, 0xc420014042, 0x18, 0x0, 0x0, 0xc420016054, 0x4, 0xc420016032, 0x5, ...)
	//		Aug 13 07:49:45.156648 cnb-rates linux-amd64[166]:         /go/src/github.com/jancajthaml-openbank/cnb-rates/pkg/cnb/sync.go:60 +0x1a8
	//		Aug 13 07:49:45.156648 cnb-rates linux-amd64[166]: created by main.main
	//		Aug 13 07:49:45.156648 cnb-rates linux-amd64[166]:         /go/src/github.com/jancajthaml-openbank/cnb-rates/main.go:103 +0x381
	//		Aug 13 07:49:45.160428 cnb-rates systemd[1]: cnb-rates.service: Main process exited, code=exited, status=2/INVALIDARGUMENT
	//		Aug 13 07:49:45.175684 cnb-rates gracefull-stop[185]: not running

	defer wg.Done()

	log.Infof("Synchronizing each %s into %s/cnb-rates", params.SyncRate, params.RootStorage)

	ticker := time.NewTicker(params.SyncRate)
	defer ticker.Stop()

	err, client := New(params.RootStorage, params.Gateway)
	if err != nil {
		log.Errorf("Unable to create new client %v", err)
		return
	}

	//terminationChan chan struct{}

	//	killChan := make(chan struct{})

	if !client.Synchronize(terminationChan) {
		return
	}

	for {
		select {
		case <-ticker.C:
			if !client.Synchronize(terminationChan) {
				return
			}
		case <-terminationChan:
			//close(killChan)
			return
		}
	}
}
