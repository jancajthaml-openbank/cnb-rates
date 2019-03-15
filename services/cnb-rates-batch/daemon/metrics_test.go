package daemon

import (
	"context"
	"testing"
	"time"

	"github.com/jancajthaml-openbank/cnb-rates-batch/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFilename(t *testing.T) {
	assert.Equal(t, "/a/b/c.batch.d", getFilename("/a/b/c.d"))
}

func TestMetricsPersist(t *testing.T) {
	cfg := config.Configuration{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	entity := NewMetrics(ctx, cfg)
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
}
