// Package services provides various services for the client, including handling metrics.
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

// SendOptions holds the options for sending metrics.
type SendOptions struct {
	BaseURL string
	HashKey string
}

// MetricsData holds runtime and gopsutil metrics.
type MetricsData struct {
	StatRuntime  *Stat
	StatGopsutil *Stat
}

type metricUpdate struct {
	name  string
	value float64
}

// NewMetricsData creates a new MetricsData instance.
// Returns a pointer to MetricsData.
func NewMetricsData() *MetricsData {
	statRuntime := &Stat{
		Data: &common.MetricItems{
			Counters: make([]common.CounterItem, 0),
			Gauges:   make([]common.GaugeItem, 0),
		},
	}
	statGopsutil := &Stat{
		Data: &common.MetricItems{
			Counters: make([]common.CounterItem, 0),
			Gauges:   make([]common.GaugeItem, 0),
		},
	}
	return &MetricsData{StatRuntime: statRuntime, StatGopsutil: statGopsutil}
}

// UpdateMetrics updates the runtime and gopsutil metrics in the given MetricsData instance.
// Parameters:
//   - data: the MetricsData instance to update
//   - counter: the current poll count
//
// Returns an error if any of the metric updates fail.
func UpdateMetrics(data *MetricsData, counter int64) error {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	gauges := getMetricsGaugeData(&ms)

	stat := data.StatRuntime
	stat.Lock.Lock()
	defer stat.Lock.Unlock()

	for _, u := range gauges {
		stat.AddOrUpdateGauge(u.name, u.value)
	}
	stat.AddOrUpdateCounter("PollCount", counter)

	err := addCPUMetrics(data)
	if err != nil {
		return err
	}

	return nil
}

func addCPUMetrics(data *MetricsData) error {
	stat := data.StatGopsutil
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

func getMetricsGaugeData(ms *runtime.MemStats) []metricUpdate {
	items := []metricUpdate{
		{"Alloc", float64(ms.Alloc)},
		{"BuckHashSys", float64(ms.BuckHashSys)},
		{"Frees", float64(ms.Frees)},
		{"BuckHashSys", float64(ms.BuckHashSys)},
		{"GCCPUFraction", float64(ms.GCCPUFraction)},
		{"GCSys", float64(ms.GCSys)},
		{"HeapAlloc", float64(ms.HeapAlloc)},
		{"HeapIdle", float64(ms.HeapIdle)},
		{"HeapInuse", float64(ms.HeapInuse)},
		{"HeapObjects", float64(ms.HeapObjects)},
		{"HeapReleased", float64(ms.HeapReleased)},
		{"HeapSys", float64(ms.HeapSys)},
		{"LastGC", float64(ms.LastGC)},
		{"Lookups", float64(ms.Lookups)},
		{"MCacheInuse", float64(ms.MCacheInuse)},
		{"MCacheSys", float64(ms.MCacheSys)},
		{"MSpanInuse", float64(ms.MSpanInuse)},
		{"MSpanSys", float64(ms.MSpanSys)},
		{"Mallocs", float64(ms.Mallocs)},
		{"NextGC", float64(ms.NextGC)},
		{"NumForcedGC", float64(ms.NumForcedGC)},
		{"NumGC", float64(ms.NumGC)},
		{"OtherSys", float64(ms.OtherSys)},
		{"PauseTotalNs", float64(ms.PauseTotalNs)},
		{"StackInuse", float64(ms.StackInuse)},
		{"StackSys", float64(ms.StackSys)},
		{"Sys", float64(ms.Sys)},
		{"TotalAlloc", float64(ms.TotalAlloc)},
		{"RandomValue", rand.Float64()}, //nolint:gosec //rand хватает
	}

	return items
}

// SendMetricsBatch sends a batch of metrics to the specified server.
// It serializes the metrics data, optionally hashes it, and sends it via HTTP POST request.
// Parameters:
//   - client: the Resty client used to send the HTTP request
//   - stat: the Stat instance containing the metrics data
//   - opts: the SendOptions containing the base URL and hash key
//
// Returns an error if the sending process fails or the server responds with a non-200 status code.
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
