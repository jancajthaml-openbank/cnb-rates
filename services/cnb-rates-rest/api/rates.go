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

package api

import (
	"encoding/json"
	"fmt"
	"github.com/jancajthaml-openbank/cnb-rates-rest/persistence"
	localfs "github.com/jancajthaml-openbank/local-fs"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetRates return existing tokens of given currency
func GetRates(storage localfs.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

		currency := c.Param("currency")
		if currency == "" {
			return fmt.Errorf("missing currency")
		}

		rates, err := persistence.LoadRates(storage, currency)
		if err != nil {
			return err
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)

		for idx, rate := range rates {
			chunk, err := json.Marshal(rate)
			if err != nil {
				return err
			}
			if idx == len(rates)-1 {
				c.Response().Write(chunk)
			} else {
				c.Response().Write(chunk)
				c.Response().Write([]byte("\n"))
			}
			c.Response().Flush()
		}

		return nil
	}
}
