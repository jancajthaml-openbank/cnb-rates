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

package metrics

import (
	"bytes"
	"fmt"
	"github.com/jancajthaml-openbank/cnb-rates-rest/utils"
	"os"
	"strconv"
	"time"
)

// MarshalJSON serializes Metrics as json bytes
func (metrics *Metrics) MarshalJSON() ([]byte, error) {
	if metrics == nil {
		return nil, fmt.Errorf("cannot marshall nil")
	}

	if metrics.gatewayLatency == nil || metrics.importLatency == nil {
		return nil, fmt.Errorf("cannot marshall nil references")
	}

	var buffer bytes.Buffer

	buffer.WriteString("{\"gatewayLatency\":")
	buffer.WriteString(strconv.FormatFloat(metrics.gatewayLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"importLatency\":")
	buffer.WriteString(strconv.FormatFloat(metrics.importLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}

// UnmarshalJSON deserializes Metrics from json bytes
func (metrics *Metrics) UnmarshalJSON(data []byte) error {
	if metrics == nil {
		return fmt.Errorf("cannot unmarshall to nil")
	}

	if metrics.gatewayLatency == nil || metrics.importLatency == nil {
		return fmt.Errorf("cannot unmarshall to nil references")
	}

	aux := &struct {
		GatewayLatency float64 `json:"gatewayLatency"`
		ImportLatency  float64 `json:"importLatency"`
	}{}

	if err := utils.JSON.Unmarshal(data, &aux); err != nil {
		return err
	}

	metrics.gatewayLatency.Update(time.Duration(aux.GatewayLatency))
	metrics.importLatency.Update(time.Duration(aux.ImportLatency))

	return nil
}

// Persist saved metrics state to storage
func (metrics *Metrics) Persist() error {
	if metrics == nil {
		return fmt.Errorf("cannot persist nil reference")
	}
	data, err := utils.JSON.Marshal(metrics)
	if err != nil {
		return err
	}
	err = metrics.storage.WriteFile("metrics.json", data)
	if err != nil {
		return err
	}
	err = os.Chmod(metrics.storage.Root+"/metrics.json", 0644)
	if err != nil {
		return err
	}
	return nil
}

// Hydrate loads metrics state from storage
func (metrics *Metrics) Hydrate() error {
	if metrics == nil {
		return fmt.Errorf("cannot hydrate nil reference")
	}
	data, err := metrics.storage.ReadFileFully("metrics.json")
	if err != nil {
		return err
	}
	err = utils.JSON.Unmarshal(data, metrics)
	if err != nil {
		return err
	}
	return nil
}
