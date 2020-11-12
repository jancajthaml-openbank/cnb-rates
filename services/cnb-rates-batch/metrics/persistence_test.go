package metrics

import (
	"io/ioutil"
	"os"
	"testing"

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
			daysProcessed:   metrics.NewCounter(),
			monthsProcessed: metrics.NewCounter(),
		}

		entity.daysProcessed.Inc(1)
		entity.monthsProcessed.Inc(2)

		actual, err := entity.MarshalJSON()

		require.Nil(t, err)

		data := []byte("{\"daysProcessed\":1,\"monthsProcessed\":2}")

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
			daysProcessed:   metrics.NewCounter(),
			monthsProcessed: metrics.NewCounter(),
		}

		data := []byte("{")
		assert.NotNil(t, entity.UnmarshalJSON(data))
	}

	t.Log("happy path")
	{
		entity := Metrics{
			daysProcessed:   metrics.NewCounter(),
			monthsProcessed: metrics.NewCounter(),
		}

		data := []byte("{\"daysProcessed\":1,\"monthsProcessed\":2}")
		require.Nil(t, entity.UnmarshalJSON(data))

		assert.Equal(t, int64(1), entity.daysProcessed.Count())
		assert.Equal(t, int64(2), entity.monthsProcessed.Count())
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
		defer os.Remove("/tmp/metrics.batch.json")

		storage, _ := localfs.NewPlaintextStorage("/tmp")

		entity := Metrics{
			storage:         storage,
			daysProcessed:   metrics.NewCounter(),
			monthsProcessed: metrics.NewCounter(),
		}

		require.Nil(t, entity.Persist())

		expected, err := entity.MarshalJSON()
		require.Nil(t, err)

		actual, err := ioutil.ReadFile("/tmp/metrics.batch.json")
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
		defer os.Remove("/tmp/metrics.batch.json")

		old := Metrics{
			daysProcessed:   metrics.NewCounter(),
			monthsProcessed: metrics.NewCounter(),
		}

		old.daysProcessed.Inc(1)
		old.monthsProcessed.Inc(2)

		data, err := old.MarshalJSON()
		require.Nil(t, err)

		require.Nil(t, ioutil.WriteFile("/tmp/metrics.batch.json", data, 0444))

		storage, _ := localfs.NewPlaintextStorage("/tmp")

		entity := Metrics{
			storage:         storage,
			daysProcessed:   metrics.NewCounter(),
			monthsProcessed: metrics.NewCounter(),
		}

		require.Nil(t, entity.Hydrate())

		assert.Equal(t, int64(1), entity.daysProcessed.Count())
		assert.Equal(t, int64(2), entity.monthsProcessed.Count())
	}
}
