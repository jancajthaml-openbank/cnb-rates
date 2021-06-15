// Copyright (c) 2016-2021, Jan Cajthaml <jan.cajthaml@gmail.com>
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
	"fmt"
	"io/ioutil"
	"time"

	"github.com/jancajthaml-openbank/cnb-rates-import/support/http"
)

// Client represents CNB gateway bound http client
type Client struct {
	gateway    string
	httpClient http.Client
}

// NewClient returns new CNB integration client
func NewClient(gateway string) *Client {
	return &Client{
		gateway:    gateway,
		httpClient: http.NewClient(),
	}
}

// GetOtherFxFor returns other fx for given date
func (client *Client) GetOtherFxFor(date time.Time) ([]byte, error) {
	if client == nil {
		return nil, fmt.Errorf("nil deference")
	}

	uri := client.gateway + "/en/financial_markets/foreign_exchange_market/other_currencies_fx_rates/fx_rates.txt?month=" + date.Format("1") + "&year=" + date.Format("2006")

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("method GetOtherFxFor(%+v) %+v", date, err)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("method GetOtherFxFor(%+v) %+v", date, err)
	} else if err == nil && resp.StatusCode != 200 {
		return nil, fmt.Errorf("method GetOtherFxFor(%+v) invalid http status %s", date, resp.Status)
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("method GetOtherFxFor(%+v) invalid response %+v", date, err)
	}

	return response, nil
}

// GetMainFxFor returns main fx for given date
func (client *Client) GetMainFxFor(date time.Time) ([]byte, error) {
	if client == nil {
		return nil, fmt.Errorf("nil deference")
	}

	uri := client.gateway + "/en/financial_markets/foreign_exchange_market/exchange_rate_fixing/daily.txt?date=" + date.Format("02+01+2006")

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("method GetMainFxFor(%+v) %+v", date, err)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("method GetMainFxFor(%+v) %+v", date, err)
	} else if err == nil && resp.StatusCode != 200 {
		return nil, fmt.Errorf("method GetMainFxFor(%+v) invalid http status %s", date, resp.Status)
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("method GetMainFxFor(%+v) invalid response %+v", date, err)
	}

	return response, nil
}
