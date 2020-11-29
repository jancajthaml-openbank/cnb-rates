package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	entity := NewMetrics("/tmp", false)

	t.Log("DayProcessed properly increments number of processed days")
	{
		require.Equal(t, int64(0), entity.daysProcessed.Count())
		entity.DayProcessed()
		assert.Equal(t, int64(1), entity.daysProcessed.Count())
	}

	t.Log("MonthProcessed properly increments number of processed months")
	{
		require.Equal(t, int64(0), entity.monthsProcessed.Count())
		entity.MonthProcessed()
		assert.Equal(t, int64(1), entity.monthsProcessed.Count())
	}
}
