package services

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
)

type SendOptions struct {
	BaseURL string
	HashKey string
}

func UpdateMetrics(stat *Stat, counter int64) {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	stat.Lock.Lock()
	defer stat.Lock.Unlock()

	stat.Data.Counters = make([]common.CounterItem, 0)
	stat.Data.Gauges = make([]common.GaugeItem, 0)

	stat.AddGauge("Alloc", float64(ms.Alloc))
	stat.AddGauge("BuckHashSys", float64(ms.BuckHashSys))
	stat.AddGauge("Frees", float64(ms.Frees))
	stat.AddGauge("BuckHashSys", float64(ms.BuckHashSys))
	stat.AddGauge("GCCPUFraction", float64(ms.GCCPUFraction))
	stat.AddGauge("GCSys", float64(ms.GCSys))
	stat.AddGauge("HeapAlloc", float64(ms.HeapAlloc))
	stat.AddGauge("HeapIdle", float64(ms.HeapIdle))
	stat.AddGauge("HeapInuse", float64(ms.HeapInuse))
	stat.AddGauge("HeapObjects", float64(ms.HeapObjects))
	stat.AddGauge("HeapReleased", float64(ms.HeapReleased))
	stat.AddGauge("HeapSys", float64(ms.HeapSys))
	stat.AddGauge("LastGC", float64(ms.LastGC))
	stat.AddGauge("Lookups", float64(ms.Lookups))
	stat.AddGauge("MCacheInuse", float64(ms.MCacheInuse))
	stat.AddGauge("MCacheSys", float64(ms.MCacheSys))
	stat.AddGauge("MSpanInuse", float64(ms.MSpanInuse))
	stat.AddGauge("MSpanSys", float64(ms.MSpanSys))
	stat.AddGauge("Mallocs", float64(ms.Mallocs))
	stat.AddGauge("NextGC", float64(ms.NextGC))
	stat.AddGauge("NumForcedGC", float64(ms.NumForcedGC))
	stat.AddGauge("NumGC", float64(ms.NumGC))
	stat.AddGauge("OtherSys", float64(ms.OtherSys))
	stat.AddGauge("PauseTotalNs", float64(ms.PauseTotalNs))
	stat.AddGauge("StackInuse", float64(ms.StackInuse))
	stat.AddGauge("StackSys", float64(ms.StackSys))
	stat.AddGauge("Sys", float64(ms.Sys))
	stat.AddGauge("TotalAlloc", float64(ms.TotalAlloc))
	stat.AddGauge("RandomValue", rand.Float64()) //nolint:gosec //rand хватает
	stat.AddCounter("PollCount", counter)
}

func SendMetricsBatch(client *resty.Client, stat *Stat, opts *SendOptions) error {
	if len(stat.Data.Counters) == 0 && len(stat.Data.Gauges) == 0 {
		return nil
	}

	statCopy := stat.Data.ToMetricsSlice()
	out, err := easyjson.Marshal(statCopy)
	if err != nil {
		return fmt.Errorf("JSON marshalling error: %w", err)
	}

	compressedOut, err := common.GZIPCompress(out)
	if err != nil {
		return fmt.Errorf("data compress error: %w", err)
	}

	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip")

	if opts.HashKey != "" {
		out := common.HashSHA256(compressedOut, []byte(opts.HashKey))
		request.SetHeader("HashSHA256", out)
	}

	_, err = request.SetBody(compressedOut).
		Post("http://" + opts.BaseURL + "/updates/")

	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}

	return nil
}
