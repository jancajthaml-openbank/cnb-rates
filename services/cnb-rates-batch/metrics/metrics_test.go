package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	entity := NewMetrics(ctx, "/tmp", time.Hour)

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
