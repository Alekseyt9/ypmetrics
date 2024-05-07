package services

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"

	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
)

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

func SendMetricsBatch(client *resty.Client, baseURL string, stat *Stat) error {
	if len(stat.Data.Counters) == 0 && len(stat.Data.Gauges) == 0 {
		return nil
	}

	copy := copyStat(stat)
	out, err := easyjson.Marshal(copy)
	if err != nil {
		return fmt.Errorf("JSON marshalling error: %w", err)
	}

	compressedOut, err := common.GZIPCompress(out)
	if err != nil {
		return fmt.Errorf("data compress error: %w", err)
	}

	_, err = client.R().
		SetHeader("Content-Type", "Content-Type: application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(compressedOut).
		Post("http://" + baseURL + "/updates")

	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}

	return nil
}

func SendMetricsJSON(client *resty.Client, baseURL string, stat *Stat) error {
	copy := copyStat(stat)

	for _, item := range copy.Gauges {
		data := common.Metrics{
			ID:    item.Name,
			MType: "gauge",
			Value: &item.Value,
		}

		out, err := easyjson.Marshal(data)
		if err != nil {
			return fmt.Errorf("JSON marshalling error: %w", err)
		}

		compressedOut, err := common.GZIPCompress(out)
		if err != nil {
			return fmt.Errorf("data compress error: %w", err)
		}

		_, err = client.R().
			SetHeader("Content-Type", "Content-Type: application/json").
			SetHeader("Content-Encoding", "gzip").
			SetHeader("Accept-Encoding", "gzip").
			SetBody(compressedOut).
			Post("http://" + baseURL + "/update")

		if err != nil {
			return fmt.Errorf("error executing request: %w", err)
		}
	}

	for _, item := range copy.Counters {
		data := common.Metrics{
			ID:    item.Name,
			MType: "counter",
			Delta: &item.Value,
		}
		out, err := easyjson.Marshal(data)
		if err != nil {
			return fmt.Errorf("JSON marshalling error: %w", err)
		}

		compressedOut, err := common.GZIPCompress(out)
		if err != nil {
			return fmt.Errorf("data compress error: %w", err)
		}

		_, err = client.R().
			SetHeader("Content-Type", "Content-Type: application/json").
			SetHeader("Content-Encoding", "gzip").
			SetHeader("Accept-Encoding", "gzip").
			SetBody(compressedOut).
			Post("http://" + baseURL + "/update")
		if err != nil {
			return fmt.Errorf("error executing request: %w", err)
		}
	}

	return nil
}

func SendMetricsURL(client *resty.Client, baseURL string, stat *Stat) error {
	copy := copyStat(stat)

	for _, item := range copy.Gauges {
		_, err := client.R().
			SetHeader("Content-Type", "Content-Type: text/plain").
			SetPathParams(map[string]string{
				"type":  item.Name,
				"value": strconv.FormatFloat(item.Value, 'f', -1, 64),
			}).
			Post("http://" + baseURL + "/update/gauge/{type}/{value}")
		if err != nil {
			return fmt.Errorf("error executing request: %w", err)
		}
	}

	for _, item := range copy.Counters {
		_, err := client.R().
			SetHeader("Content-Type", "Content-Type: text/plain").
			SetPathParams(map[string]string{
				"type":  item.Name,
				"value": strconv.FormatInt(item.Value, 10),
			}).
			Post("http://" + baseURL + "/update/counter/{type}/{value}")
		if err != nil {
			return fmt.Errorf("error executing request: %w", err)
		}
	}

	return nil
}

func copyStat(stat *Stat) common.MetricsBatch {
	copy := common.MetricsBatch{
		Counters: make([]common.CounterItem, 0, len(stat.Data.Counters)),
		Gauges:   make([]common.GaugeItem, 0, len(stat.Data.Counters)),
	}

	stat.Lock.RLock()
	for _, item := range stat.Data.Gauges {
		copy.Gauges = append(copy.Gauges, common.GaugeItem{Name: item.Name, Value: item.Value})
	}
	for _, item := range stat.Data.Counters {
		copy.Counters = append(copy.Counters, common.CounterItem{Name: item.Name, Value: item.Value})
	}
	stat.Lock.RUnlock()

	return copy
}
