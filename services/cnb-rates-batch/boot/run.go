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

package boot

import (
	"context"
	"github.com/jancajthaml-openbank/cnb-rates-batch/support/host"
	"os/signal"
	"syscall"
)

// Stop stops all daemons
func (prog Program) Stop() {
	prog.pool.Stop()
}

// Start starts all daemons and blocks until INT or TERM signal is received
func (prog Program) Start() {
	log.Info().Msg("Program Starting")
	go prog.pool.Start(context.WithCancel(context.Background()))
	host.NotifyServiceReady()
	log.Info().Msg("Program Started")
	signal.Notify(prog.interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	select {
	case <-prog.interrupt:
		break
	case <-prog.pool.Done():
		break
	}
	log.Info().Msg("Program Stopping")
	if err := host.NotifyServiceStopping(); err != nil {
		log.Error().Msg(err.Error())
	}
	prog.pool.Stop()
	<-prog.pool.Done()
	log.Info().Msg("Program Stopped")
}
