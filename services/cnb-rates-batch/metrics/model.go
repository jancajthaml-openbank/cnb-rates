// Copyright (c) 2016-2019, Jan Cajthaml <jan.cajthaml@gmail.com>
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
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jancajthaml-openbank/cnb-rates-batch/utils"
	metrics "github.com/rcrowley/go-metrics"
)

// Metrics represents metrics subroutine
type Metrics struct {
	utils.DaemonSupport
	output          string
	refreshRate     time.Duration
	daysProcessed   metrics.Counter
	monthsProcessed metrics.Counter
}

// NewMetrics returns metrics fascade
func NewMetrics(ctx context.Context, output string, refreshRate time.Duration) Metrics {
	return Metrics{
		DaemonSupport:   utils.NewDaemonSupport(ctx),
		output:          output,
		refreshRate:     refreshRate,
		daysProcessed:   metrics.NewCounter(),
		monthsProcessed: metrics.NewCounter(),
	}
}

// MarshalJSON serialises Metrics as json bytes
func (metrics *Metrics) MarshalJSON() ([]byte, error) {
	if metrics == nil {
		return nil, fmt.Errorf("cannot marshall nil")
	}

	if metrics.daysProcessed == nil || metrics.monthsProcessed == nil {
		return nil, fmt.Errorf("cannot marshall nil references")
	}

	var buffer bytes.Buffer

	buffer.WriteString("{\"daysProcessed\":")
	buffer.WriteString(strconv.FormatInt(metrics.daysProcessed.Count(), 10))
	buffer.WriteString(",\"monthsProcessed\":")
	buffer.WriteString(strconv.FormatInt(metrics.monthsProcessed.Count(), 10))
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}

// UnmarshalJSON deserializes Metrics from json bytes
func (metrics *Metrics) UnmarshalJSON(data []byte) error {
	if metrics == nil {
		return fmt.Errorf("cannot unmarshall to nil")
	}

	if metrics.daysProcessed == nil || metrics.monthsProcessed == nil {
		return fmt.Errorf("cannot unmarshall to nil references")
	}

	aux := &struct {
		DaysProcessed   int64 `json:"daysProcessed"`
		MonthsProcessed int64 `json:"monthsProcessed"`
	}{}

	if err := utils.JSON.Unmarshal(data, &aux); err != nil {
		return err
	}

	metrics.daysProcessed.Clear()
	metrics.daysProcessed.Inc(aux.DaysProcessed)

	metrics.monthsProcessed.Clear()
	metrics.monthsProcessed.Inc(aux.MonthsProcessed)

	return nil
}