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
	"os"

	"github.com/codegangsta/cli"
)

func GlobalFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "output-directory, o",
			Usage: "directory where cache is stored",
		},
	}
}

func All() []cli.Command {
	return []cli.Command{
		{
			Name:   "sync",
			Usage:  "import rates until today",
			Action: try(CmdSyncRates),
		},
	}
}

func NotFound(c *cli.Context, command string) {
	cli.ShowAppHelp(c)
	os.Exit(2)
}

func try(fn func(c *cli.Context) error) func(c *cli.Context) error {
	return func(c *cli.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("Runtime Error")
				fmt.Println("command failed with", err)
				fmt.Println(r)
			}
		}()

		if err = fn(c); err != nil {
			panic(err)
		}
		return nil
	}
}
