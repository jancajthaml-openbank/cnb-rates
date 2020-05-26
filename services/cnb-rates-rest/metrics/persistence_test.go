package metrics

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	localfs "github.com/jancajthaml-openbank/local-fs"
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

func TestPersist(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		assert.EqualError(t, entity.Persist(), "cannot persist nil reference")
	}

	t.Log("error when marshalling fails")
	{
		entity := Metrics{}
		assert.EqualError(t, entity.Persist(), "cannot marshall nil references")
	}

	t.Log("happy path")
	{
		defer os.Remove("/tmp/metrics.json")

		entity := Metrics{
			storage:        localfs.NewPlaintextStorage("/tmp"),
			gatewayLatency: metrics.NewTimer(),
			importLatency:  metrics.NewTimer(),
		}

		require.Nil(t, entity.Persist())

		expected, err := entity.MarshalJSON()
		require.Nil(t, err)

		actual, err := ioutil.ReadFile("/tmp/metrics.json")
		require.Nil(t, err)

		assert.Equal(t, expected, actual)
	}
}

func TestHydrate(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		assert.EqualError(t, entity.Hydrate(), "cannot hydrate nil reference")
	}

	t.Log("happy path")
	{
		defer os.Remove("/tmp/metrics.json")

		old := Metrics{
			gatewayLatency: metrics.NewTimer(),
			importLatency:  metrics.NewTimer(),
		}

		old.gatewayLatency.Update(time.Duration(1))
		old.importLatency.Update(time.Duration(2))

		data, err := old.MarshalJSON()
		require.Nil(t, err)

		require.Nil(t, ioutil.WriteFile("/tmp/metrics.json", data, 0444))

		entity := Metrics{
			storage:        localfs.NewPlaintextStorage("/tmp"),
			gatewayLatency: metrics.NewTimer(),
			importLatency:  metrics.NewTimer(),
		}

		require.Nil(t, entity.Hydrate())

		assert.Equal(t, float64(1), entity.gatewayLatency.Percentile(0.95))
		assert.Equal(t, float64(2), entity.importLatency.Percentile(0.95))
	}
}
