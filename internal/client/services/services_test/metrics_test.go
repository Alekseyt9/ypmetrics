package services_test

import (
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/client/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetrics(t *testing.T) {
	var counter int64
	data := services.NewMetricsData()

	testsGauge := []string{
		"Alloc", "BuckHashSys", "Frees", "BuckHashSys",
		"GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle",
		"HeapInuse", "HeapObjects", "HeapReleased", "HeapSys",
		"LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC",
		"NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs",
		"StackInuse", "StackSys", "Sys", "TotalAlloc",
		"RandomValue",
	}

	testsCounter := []string{
		"PollCount",
	}

	err := services.UpdateMetrics(data, counter)
	require.NoError(t, err)
	stat := data.StatRuntime

	for _, name := range testsGauge {
		t.Run(name, func(t *testing.T) {
			stat.Lock.RLock()
			defer stat.Lock.RUnlock()
			_, ok := stat.FindGauge(name)
			assert.True(t, ok)
		})
	}

	for _, name := range testsCounter {
		t.Run(name, func(t *testing.T) {
			stat.Lock.RLock()
			defer stat.Lock.RUnlock()
			_, ok := stat.FindCounter(name)
			assert.True(t, ok)
		})
	}
}
