package metrics

import (
	"io/ioutil"
	"os"
	"testing"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersist(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		assert.EqualError(t, entity.Persist(), "cannot persist nil reference")
	}

	t.Log("error when marshalling fails")
	{
		entity := Metrics{}
		assert.EqualError(t, entity.Persist(), "json: error calling MarshalJSON for type *metrics.Metrics: cannot marshall nil references")
	}

	t.Log("error when cannot open tempfile for writing")
	{
		entity := Metrics{
			output:          "/sys/kernel/security",
			daysProcessed:   metrics.NewCounter(),
			monthsProcessed: metrics.NewCounter(),
		}

		assert.NotNil(t, entity.Persist())
	}

	t.Log("happy path")
	{
		tmpfile, err := ioutil.TempFile(os.TempDir(), "test_metrics_persist")

		require.Nil(t, err)
		defer os.Remove(tmpfile.Name())

		entity := Metrics{
			output:          tmpfile.Name(),
			daysProcessed:   metrics.NewCounter(),
			monthsProcessed: metrics.NewCounter(),
		}

		require.Nil(t, entity.Persist())

		expected, err := entity.MarshalJSON()
		require.Nil(t, err)

		actual, err := ioutil.ReadFile(tmpfile.Name())
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
		tmpfile, err := ioutil.TempFile(os.TempDir(), "test_metrics_hydrate")

		require.Nil(t, err)
		defer os.Remove(tmpfile.Name())

		old := Metrics{
			daysProcessed:   metrics.NewCounter(),
			monthsProcessed: metrics.NewCounter(),
		}

		old.daysProcessed.Inc(1)
		old.monthsProcessed.Inc(2)

		data, err := old.MarshalJSON()
		require.Nil(t, err)

		require.Nil(t, ioutil.WriteFile(tmpfile.Name(), data, 0444))

		entity := Metrics{
			output:          tmpfile.Name(),
			daysProcessed:   metrics.NewCounter(),
			monthsProcessed: metrics.NewCounter(),
		}

		require.Nil(t, entity.Hydrate())

		assert.Equal(t, int64(1), entity.daysProcessed.Count())
		assert.Equal(t, int64(2), entity.monthsProcessed.Count())
	}
}
