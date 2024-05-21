package services

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type SendOptions struct {
	BaseURL string
	HashKey string
}

func UpdateMetrics1(stat *Stat, counter int64) error {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	stat.Lock.Lock()
	defer stat.Lock.Unlock()

	stat.AddOrUpdateGauge("Alloc", float64(ms.Alloc))
	stat.AddOrUpdateGauge("BuckHashSys", float64(ms.BuckHashSys))
	stat.AddOrUpdateGauge("Frees", float64(ms.Frees))
	stat.AddOrUpdateGauge("BuckHashSys", float64(ms.BuckHashSys))
	stat.AddOrUpdateGauge("GCCPUFraction", float64(ms.GCCPUFraction))
	stat.AddOrUpdateGauge("GCSys", float64(ms.GCSys))
	stat.AddOrUpdateGauge("HeapAlloc", float64(ms.HeapAlloc))
	stat.AddOrUpdateGauge("HeapIdle", float64(ms.HeapIdle))
	stat.AddOrUpdateGauge("HeapInuse", float64(ms.HeapInuse))
	stat.AddOrUpdateGauge("HeapObjects", float64(ms.HeapObjects))
	stat.AddOrUpdateGauge("HeapReleased", float64(ms.HeapReleased))
	stat.AddOrUpdateGauge("HeapSys", float64(ms.HeapSys))
	stat.AddOrUpdateGauge("LastGC", float64(ms.LastGC))
	stat.AddOrUpdateGauge("Lookups", float64(ms.Lookups))
	stat.AddOrUpdateGauge("MCacheInuse", float64(ms.MCacheInuse))
	stat.AddOrUpdateGauge("MCacheSys", float64(ms.MCacheSys))
	stat.AddOrUpdateGauge("MSpanInuse", float64(ms.MSpanInuse))
	stat.AddOrUpdateGauge("MSpanSys", float64(ms.MSpanSys))
	stat.AddOrUpdateGauge("Mallocs", float64(ms.Mallocs))
	stat.AddOrUpdateGauge("NextGC", float64(ms.NextGC))
	stat.AddOrUpdateGauge("NumForcedGC", float64(ms.NumForcedGC))
	stat.AddOrUpdateGauge("NumGC", float64(ms.NumGC))
	stat.AddOrUpdateGauge("OtherSys", float64(ms.OtherSys))
	stat.AddOrUpdateGauge("PauseTotalNs", float64(ms.PauseTotalNs))
	stat.AddOrUpdateGauge("StackInuse", float64(ms.StackInuse))
	stat.AddOrUpdateGauge("StackSys", float64(ms.StackSys))
	stat.AddOrUpdateGauge("Sys", float64(ms.Sys))
	stat.AddOrUpdateGauge("TotalAlloc", float64(ms.TotalAlloc))
	stat.AddOrUpdateGauge("RandomValue", rand.Float64()) //nolint:gosec //rand хватает
	stat.AddOrUpdateCounter("PollCount", counter)

	return nil
}

func UpdateMetrics2(stat *Stat, counter int64) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	cpuUtilizations, err := cpu.Percent(1*time.Second, true)
	if err != nil {
		return err
	}

	stat.Lock.Lock()
	defer stat.Lock.Unlock()

	stat.AddOrUpdateGauge("TotalMemory", float64(v.Total))
	stat.AddOrUpdateGauge("FreeMemory", float64(v.Free))

	for i, utilization := range cpuUtilizations {
		stat.AddOrUpdateGauge("CPUutilization"+strconv.Itoa(i), float64(utilization))
	}

	return nil
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
