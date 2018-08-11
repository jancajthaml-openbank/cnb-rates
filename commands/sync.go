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

package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/jancajthaml-openbank/cnb-rates/pkg/cnb"
)

func CmdSyncRates(c *cli.Context) error {
	fmt.Println("Start")

	client := cnb.New(c.GlobalString("output-directory"))

	fmt.Println("> synchonizing main currencies fx daily rates")
	client.SyncMainRatesDaily()

	fmt.Println("> synchonizing other currencies fx monthly rates")
	client.SyncOtherRatesMonthly()

	fmt.Println("> synchonizing main currencies fx yearly rates")
	client.SyncMainRatesYearly()

	fmt.Println("End")
	return nil
}
