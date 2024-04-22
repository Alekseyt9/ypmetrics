package services_test

import (
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/client/services"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetrics(t *testing.T) {
	var counter int64
	stat := &services.Stat{
		CounterMap: make(map[string]int64),
		GaugeMap:   make(map[string]float64),
	}

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

	services.UpdateMetrics(stat, counter)

	for _, name := range testsGauge {
		t.Run(name, func(t *testing.T) {
			stat.GaugeLock.RLock()
			defer stat.GaugeLock.RUnlock()
			_, ok := stat.GaugeMap[name]
			assert.True(t, ok)
		})
	}

	for _, name := range testsCounter {
		t.Run(name, func(t *testing.T) {
			stat.CounterLock.RLock()
			defer stat.CounterLock.RUnlock()
			_, ok := stat.CounterMap[name]
			assert.True(t, ok)
		})
	}
}
