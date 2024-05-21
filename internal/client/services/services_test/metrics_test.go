package services_test

import (
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/client/services"
	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetrics(t *testing.T) {
	var counter int64
	stat := &services.Stat{
		Data: &common.MetricItems{
			Counters: make([]common.CounterItem, 0),
			Gauges:   make([]common.GaugeItem, 0),
		},
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

	services.UpdateMetrics1(stat, counter)

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
