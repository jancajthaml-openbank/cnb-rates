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

package main

import (
	"os"
	"os/signal"
	"sort"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/jancajthaml-openbank/cnb-rates/commands"
	"github.com/sirupsen/logrus"
)

func main() {
	app := cli.NewApp()
	app.Name = "cnb-rates"
	app.Author = "Jan Cajthaml <jan.cajthaml@gmail.com>"
	app.Usage = "CNB daily exchange rate importer"

	app.Flags = commands.GlobalFlags()
	app.Commands = commands.All()
	app.CommandNotFound = commands.NotFound

	app.Before = beforeHook
	app.After = afterHook

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func beforeHook(c *cli.Context) error {
	setLogLevel(c)
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-exitSignal
		afterHook(nil)
		os.Exit(2)
	}()

	return nil
}

func afterHook(c *cli.Context) error {
	return nil
}

func setLogLevel(c *cli.Context) {
	if c.GlobalBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
