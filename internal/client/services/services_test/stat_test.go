package services_test

import (
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/client/services"
	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOrUpdateGauge(t *testing.T) {
	stat := &services.Stat{
		Data: &common.MetricItems{
			Gauges:   []common.GaugeItem{},
			Counters: []common.CounterItem{},
		},
	}

	stat.AddOrUpdateGauge("test_gauge", 1.1)
	require.Len(t, stat.Data.Gauges, 1)
	assert.Equal(t, 1.1, stat.Data.Gauges[0].Value)

	stat.AddOrUpdateGauge("test_gauge", 2.1)
	require.Len(t, stat.Data.Gauges, 1)
	assert.Equal(t, 2.1, stat.Data.Gauges[0].Value)
}

func TestAddOrUpdateCounter(t *testing.T) {
	stat := &services.Stat{
		Data: &common.MetricItems{
			Gauges:   []common.GaugeItem{},
			Counters: []common.CounterItem{},
		},
	}

	stat.AddOrUpdateCounter("test_counter", 1)
	require.Len(t, stat.Data.Counters, 1)
	assert.Equal(t, int64(1), stat.Data.Counters[0].Value)

	stat.AddOrUpdateCounter("test_counter", 2)
	require.Len(t, stat.Data.Counters, 1)
	assert.Equal(t, int64(2), stat.Data.Counters[0].Value)
}

func TestFindGauge(t *testing.T) {
	stat := &services.Stat{
		Data: &common.MetricItems{
			Gauges:   []common.GaugeItem{{Name: "test_gauge", Value: 123.45}},
			Counters: []common.CounterItem{},
		},
	}

	gauge, ok := stat.FindGauge("test_gauge")
	require.True(t, ok)
	assert.Equal(t, 123.45, gauge.Value)

	_, ok = stat.FindGauge("non_existent_gauge")
	assert.False(t, ok)
}

func TestFindCounter(t *testing.T) {
	stat := &services.Stat{
		Data: &common.MetricItems{
			Gauges:   []common.GaugeItem{},
			Counters: []common.CounterItem{{Name: "test_counter", Value: 123}},
		},
	}

	counter, ok := stat.FindCounter("test_counter")
	require.True(t, ok)
	assert.Equal(t, int64(123), counter.Value)

	_, ok = stat.FindCounter("non_existent_counter")
	assert.False(t, ok)
}
