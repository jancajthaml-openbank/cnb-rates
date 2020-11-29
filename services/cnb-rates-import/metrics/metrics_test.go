package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	entity := NewMetrics("/tmp", false)
	delay := 1e8
	delta := 1e8

	t.Log("TimeGatewayLatency properly times gateway latency")
	{
		require.Equal(t, int64(0), entity.gatewayLatency.Count())
		entity.TimeGatewayLatency(func() {
			time.Sleep(time.Duration(delay))
		})
		assert.Equal(t, int64(1), entity.gatewayLatency.Count())
		assert.InDelta(t, entity.gatewayLatency.Percentile(0.95), delay, delta)
	}

	t.Log("TimeImportLatency properly times import latency")
	{
		require.Equal(t, int64(0), entity.importLatency.Count())
		entity.TimeImportLatency(func() {
			time.Sleep(time.Duration(delay))
		})
		assert.Equal(t, int64(1), entity.importLatency.Count())
		assert.InDelta(t, entity.importLatency.Percentile(0.95), delay, delta)
	}

	t.Log("DayImported properly increments number of imported days")
	{
		require.Equal(t, int64(0), entity.daysImported.Count())
		entity.DayImported()
		assert.Equal(t, int64(1), entity.daysImported.Count())
	}

	t.Log("MonthImported properly increments number of imported months")
	{
		require.Equal(t, int64(0), entity.monthsImported.Count())
		entity.MonthImported()
		assert.Equal(t, int64(1), entity.monthsImported.Count())
	}
}
