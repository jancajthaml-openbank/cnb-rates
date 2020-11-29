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
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// MarshalJSON serializes Metrics as json bytes
func (metrics *Metrics) MarshalJSON() ([]byte, error) {
	if metrics == nil {
		return nil, fmt.Errorf("cannot marshall nil")
	}

	if metrics.gatewayLatency == nil || metrics.importLatency == nil ||
		metrics.daysImported == nil || metrics.monthsImported == nil {
		return nil, fmt.Errorf("cannot marshall nil references")
	}

	var buffer bytes.Buffer

	buffer.WriteString("{\"gatewayLatency\":")
	buffer.WriteString(strconv.FormatFloat(metrics.gatewayLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"importLatency\":")
	buffer.WriteString(strconv.FormatFloat(metrics.importLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"daysImported\":")
	buffer.WriteString(strconv.FormatInt(metrics.daysImported.Count(), 10))
	buffer.WriteString(",\"monthsImported\":")
	buffer.WriteString(strconv.FormatInt(metrics.monthsImported.Count(), 10))
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}

// UnmarshalJSON deserializes Metrics from json bytes
func (metrics *Metrics) UnmarshalJSON(data []byte) error {
	if metrics == nil {
		return fmt.Errorf("cannot unmarshall to nil")
	}

	if metrics.gatewayLatency == nil || metrics.importLatency == nil ||
		metrics.daysImported == nil || metrics.monthsImported == nil {
		return fmt.Errorf("cannot unmarshall to nil references")
	}

	aux := &struct {
		GatewayLatency float64 `json:"gatewayLatency"`
		ImportLatency  float64 `json:"importLatency"`
		DaysImported   int64   `json:"daysImported"`
		MonthsImported int64   `json:"monthsImported"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	metrics.gatewayLatency.Update(time.Duration(aux.GatewayLatency))
	metrics.importLatency.Update(time.Duration(aux.ImportLatency))
	metrics.daysImported.Clear()
	metrics.daysImported.Inc(aux.DaysImported)
	metrics.monthsImported.Clear()
	metrics.monthsImported.Inc(aux.MonthsImported)

	return nil
}

// Persist saved metrics state to storage
func (metrics *Metrics) Persist() error {
	if metrics == nil {
		return fmt.Errorf("cannot persist nil reference")
	}
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	err = metrics.storage.WriteFile("metrics.import.json", data)
	if err != nil {
		return err
	}
	err = metrics.storage.Chmod("metrics.import.json", 0644)
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
	data, err := metrics.storage.ReadFileFully("metrics.import.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, metrics)
	if err != nil {
		return err
	}
	return nil
}
