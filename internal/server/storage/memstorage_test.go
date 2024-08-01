package storage_test

import (
	"context"
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/Alekseyt9/ypmetrics/internal/server/filedump"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGaugeStorage(t *testing.T) {
	tests := []struct {
		name   string
		metric string
		set    float64
		want   float64
	}{
		{
			name:   "test1",
			set:    1.0,
			want:   1.0,
			metric: "m1",
		},
		{
			name:   "test2",
			set:    -1.0,
			want:   -1.0,
			metric: "m1",
		},
		{
			name:   "test3",
			set:    0,
			want:   0,
			metric: "m2",
		},
		{
			name:   "test4",
			set:    -0.0000000000000001,
			want:   -0.0000000000000001,
			metric: "m2",
		},
	}

	store := storage.NewMemStorage()
	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := store.SetGauge(ctx, test.metric, test.set)
			require.NoError(t, err)
			var v float64
			v, err = store.GetGauge(ctx, test.metric)
			require.NoError(t, err)
			assert.InDelta(t, test.want, v, 0.01)
		})
	}
}

func TestCounterStorage(t *testing.T) {
	tests := []struct {
		name   string
		metric string
		set    int64
		want   int64
	}{
		{
			name:   "test1",
			set:    1,
			want:   1,
			metric: "m1",
		},
		{
			name:   "test2",
			set:    2,
			want:   3,
			metric: "m1",
		},
		{
			name:   "test3",
			set:    0,
			want:   0,
			metric: "m2",
		},
		{
			name:   "test4",
			set:    0,
			want:   0,
			metric: "m2",
		},
	}

	store := storage.NewMemStorage()
	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := store.SetCounter(ctx, test.metric, test.set)
			require.NoError(t, err)
			var v int64
			v, err = store.GetCounter(ctx, test.metric)
			require.NoError(t, err)
			assert.InDelta(t, test.want, v, 0.01)
		})
	}
}

func TestGauges(t *testing.T) {
	store := storage.NewMemStorage()
	ctx := context.Background()

	gauges := []common.GaugeItem{
		{Name: "gauge1", Value: 1.1},
		{Name: "gauge2", Value: 2.2},
		{Name: "gauge3", Value: 3.3},
	}

	err := store.SetGauges(ctx, gauges)
	require.NoError(t, err, "SetGauges failed")

	result, err := store.GetGauges(ctx)
	require.NoError(t, err, "GetGauges failed")

	assert.Equal(t, len(gauges), len(result), "Expected number of gauges does not match")
	for _, gauge := range gauges {
		found := false
		for _, res := range result {
			if res.Name == gauge.Name && res.Value == gauge.Value {
				found = true
				break
			}
		}
		assert.True(t, found, "Gauge %v not found in result", gauge)
	}
}

func TestCounters(t *testing.T) {
	store := storage.NewMemStorage()
	ctx := context.Background()

	counters := []common.CounterItem{
		{Name: "counter1", Value: 100},
		{Name: "counter2", Value: 200},
		{Name: "counter3", Value: 300},
	}

	err := store.SetCounters(ctx, counters)
	require.NoError(t, err, "SetCounters failed")

	result, err := store.GetCounters(ctx)
	require.NoError(t, err, "GetCounters failed")

	assert.Equal(t, len(counters), len(result), "Expected number of counters does not match")
	for _, counter := range counters {
		found := false
		for _, res := range result {
			if res.Name == counter.Name && res.Value == counter.Value {
				found = true
				break
			}
		}
		assert.True(t, found, "Counter %v not found in result", counter)
	}
}

func TestLoadFromDump(t *testing.T) {
	store := storage.NewMemStorage()
	dump := &filedump.FileDump{
		CounterData: map[string]int64{"counter1": 10, "counter2": 20},
		GaugeData:   map[string]float64{"gauge1": 1.1, "gauge2": 2.2},
	}

	store.LoadFromDump(dump)

	counter, err := store.GetCounter(context.Background(), "counter1")
	require.NoError(t, err)
	assert.Equal(t, int64(10), counter)

	gauge, err := store.GetGauge(context.Background(), "gauge1")
	require.NoError(t, err)
	assert.Equal(t, 1.1, gauge)
}

func TestSaveToDump(t *testing.T) {
	store := storage.NewMemStorage()
	store.SetCounter(context.Background(), "counter1", 10)
	store.SetGauge(context.Background(), "gauge1", 1.1)

	dump := &filedump.FileDump{
		CounterData: make(map[string]int64),
		GaugeData:   make(map[string]float64),
	}
	store.SaveToDump(dump)

	assert.Equal(t, int64(10), dump.CounterData["counter1"])
	assert.Equal(t, 1.1, dump.GaugeData["gauge1"])
}

func TestSaveToFile(t *testing.T) {
	store := storage.NewMemStorage()
	store.SetCounter(context.Background(), "counter1", 10)
	store.SetGauge(context.Background(), "gauge1", 1.1)

	filePath := "test_dump.json"
	err := store.SaveToFile(filePath)
	require.NoError(t, err)

	dump := &filedump.FileDump{
		CounterData: make(map[string]int64),
		GaugeData:   make(map[string]float64),
	}
	dc := filedump.NewController()
	err = dc.Load(dump, filePath)
	require.NoError(t, err)

	assert.Equal(t, int64(10), dump.CounterData["counter1"])
	assert.Equal(t, 1.1, dump.GaugeData["gauge1"])
}

func TestLoadFromFile(t *testing.T) {
	store := storage.NewMemStorage()

	dump := &filedump.FileDump{
		CounterData: map[string]int64{"counter1": 10, "counter2": 20},
		GaugeData:   map[string]float64{"gauge1": 1.1, "gauge2": 2.2},
	}
	filePath := "test_load_dump.json"
	dc := filedump.NewController()
	err := dc.Save(dump, filePath)
	require.NoError(t, err)

	err = store.LoadFromFile(filePath)
	require.NoError(t, err)

	counter, err := store.GetCounter(context.Background(), "counter1")
	require.NoError(t, err)
	assert.Equal(t, int64(10), counter)

	gauge, err := store.GetGauge(context.Background(), "gauge1")
	require.NoError(t, err)
	assert.Equal(t, 1.1, gauge)
}

func TestPing(t *testing.T) {
	store := storage.NewMemStorage()
	err := store.Ping(context.Background())
	assert.NoError(t, err)
}
