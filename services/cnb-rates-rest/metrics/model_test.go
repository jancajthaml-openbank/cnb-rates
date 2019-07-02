package metrics

import (
	"testing"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSON(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		_, err := entity.MarshalJSON()
		assert.EqualError(t, err, "cannot marshall nil")
	}

	t.Log("error when values are nil")
	{
		entity := Metrics{}
		_, err := entity.MarshalJSON()
		assert.EqualError(t, err, "cannot marshall nil references")
	}

	t.Log("happy path")
	{
		entity := Metrics{
			gatewayLatency: metrics.NewTimer(),
			importLatency:  metrics.NewTimer(),
		}

		entity.gatewayLatency.Update(time.Duration(1))
		entity.importLatency.Update(time.Duration(2))

		actual, err := entity.MarshalJSON()

		require.Nil(t, err)

		data := []byte("{\"gatewayLatency\":1,\"importLatency\":2}")

		assert.Equal(t, data, actual)
	}
}

func TestUnmarshalJSON(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		err := entity.UnmarshalJSON([]byte(""))
		assert.EqualError(t, err, "cannot unmarshall to nil")
	}

	t.Log("error when values are nil")
	{
		entity := Metrics{}
		err := entity.UnmarshalJSON([]byte(""))
		assert.EqualError(t, err, "cannot unmarshall to nil references")
	}

	t.Log("error on malformed data")
	{
		entity := Metrics{
			gatewayLatency: metrics.NewTimer(),
			importLatency:  metrics.NewTimer(),
		}

		data := []byte("{")
		assert.NotNil(t, entity.UnmarshalJSON(data))
	}

	t.Log("happy path")
	{
		entity := Metrics{
			gatewayLatency: metrics.NewTimer(),
			importLatency:  metrics.NewTimer(),
		}

		data := []byte("{\"gatewayLatency\":1,\"importLatency\":2}")
		require.Nil(t, entity.UnmarshalJSON(data))

		assert.Equal(t, float64(1), entity.gatewayLatency.Percentile(0.95))
		assert.Equal(t, float64(2), entity.importLatency.Percentile(0.95))
	}
}
