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
		assert.NotNil(t, err)
	}

	t.Log("error when values are nil")
	{
		entity := Metrics{}
		_, err := entity.MarshalJSON()
		assert.NotNil(t, err)
	}

	t.Log("happy path")
	{
		entity := Metrics{
			daysImported:   metrics.NewCounter(),
			monthsImported: metrics.NewCounter(),
			gatewayLatency: metrics.NewTimer(),
			importLatency:  metrics.NewTimer(),
		}

		entity.gatewayLatency.Update(time.Duration(1))
		entity.importLatency.Update(time.Duration(2))
		entity.daysImported.Inc(3)
		entity.monthsImported.Inc(4)

		actual, err := entity.MarshalJSON()

		require.Nil(t, err)

		data := []byte("{\"gatewayLatency\":1,\"importLatency\":2,\"daysImported\":3,\"monthsImported\":4}")

		assert.Equal(t, data, actual)
	}
}

func TestUnmarshalJSON(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		err := entity.UnmarshalJSON([]byte(""))
		assert.NotNil(t, err)
	}

	t.Log("error when values are nil")
	{
		entity := Metrics{}
		err := entity.UnmarshalJSON([]byte(""))
		assert.NotNil(t, err)
	}

	t.Log("error on malformed data")
	{
		entity := Metrics{
			daysImported:   metrics.NewCounter(),
			monthsImported: metrics.NewCounter(),
			gatewayLatency: metrics.NewTimer(),
			importLatency:  metrics.NewTimer(),
		}

		data := []byte("{")
		assert.NotNil(t, entity.UnmarshalJSON(data))
	}

	t.Log("happy path")
	{
		entity := Metrics{
			daysImported:   metrics.NewCounter(),
			monthsImported: metrics.NewCounter(),
			gatewayLatency: metrics.NewTimer(),
			importLatency:  metrics.NewTimer(),
		}

		data := []byte("{\"gatewayLatency\":1,\"importLatency\":2,\"daysImported\":3,\"monthsImported\":4}")
		require.Nil(t, entity.UnmarshalJSON(data))

		assert.Equal(t, float64(1), entity.gatewayLatency.Percentile(0.95))
		assert.Equal(t, float64(2), entity.importLatency.Percentile(0.95))
		assert.Equal(t, int64(3), entity.daysImported.Count())
		assert.Equal(t, int64(4), entity.monthsImported.Count())
	}
}

func TestPersist(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		assert.NotNil(t, entity.Persist())
	}

	t.Log("error when marshaling fails")
	{
		entity := Metrics{}
		assert.NotNil(t, entity.Persist())
	}

	t.Log("happy path")
	{
		defer os.Remove("/tmp/metrics.import.json")

		storage, _ := localfs.NewPlaintextStorage("/tmp")

		entity := Metrics{
			storage:        storage,
			daysImported:   metrics.NewCounter(),
			monthsImported: metrics.NewCounter(),
			gatewayLatency: metrics.NewTimer(),
			importLatency:  metrics.NewTimer(),
		}

		require.Nil(t, entity.Persist())

		expected, err := entity.MarshalJSON()
		require.Nil(t, err)

		actual, err := ioutil.ReadFile("/tmp/metrics.import.json")
		require.Nil(t, err)

		assert.Equal(t, expected, actual)
	}
}

func TestHydrate(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		assert.NotNil(t, entity.Hydrate())
	}

	t.Log("happy path")
	{
		defer os.Remove("/tmp/metrics.import.json")

		old := Metrics{
			daysImported:   metrics.NewCounter(),
			monthsImported: metrics.NewCounter(),
			gatewayLatency: metrics.NewTimer(),
			importLatency:  metrics.NewTimer(),
		}

		old.gatewayLatency.Update(time.Duration(1))
		old.importLatency.Update(time.Duration(2))
		old.daysImported.Inc(3)
		old.monthsImported.Inc(4)

		data, err := old.MarshalJSON()
		require.Nil(t, err)

		require.Nil(t, ioutil.WriteFile("/tmp/metrics.import.json", data, 0444))

		storage, _ := localfs.NewPlaintextStorage("/tmp")

		entity := Metrics{
			storage:        storage,
			daysImported:   metrics.NewCounter(),
			monthsImported: metrics.NewCounter(),
			gatewayLatency: metrics.NewTimer(),
			importLatency:  metrics.NewTimer(),
		}

		require.Nil(t, entity.Hydrate())

		assert.Equal(t, float64(1), entity.gatewayLatency.Percentile(0.95))
		assert.Equal(t, float64(2), entity.importLatency.Percentile(0.95))
		assert.Equal(t, int64(3), entity.daysImported.Count())
		assert.Equal(t, int64(4), entity.monthsImported.Count())
	}
}
