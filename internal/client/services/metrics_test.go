package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateMetrics(t *testing.T) {
	var counter int64 = 0
	gMap := make(map[string]float64)
	cMap := make(map[string]int64)

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

	UpdateMetrics(gMap, cMap, counter)

	for _, name := range testsGauge {
		t.Run(name, func(t *testing.T) {
			_, ok := gMap[name]
			assert.True(t, ok)
		})
	}

	for _, name := range testsCounter {
		t.Run(name, func(t *testing.T) {
			_, ok := cMap[name]
			assert.True(t, ok)
		})
	}

}
