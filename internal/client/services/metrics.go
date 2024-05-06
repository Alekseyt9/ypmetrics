package services

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"

	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
)

type Stat struct {
	CounterMap map[string]int64
	GaugeMap   map[string]float64
	MapLock    sync.RWMutex
}

func UpdateMetrics(stat *Stat, counter int64) {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	stat.MapLock.Lock()
	defer stat.MapLock.Unlock()

	stat.GaugeMap["Alloc"] = float64(ms.Alloc)
	stat.GaugeMap["BuckHashSys"] = float64(ms.BuckHashSys)
	stat.GaugeMap["Frees"] = float64(ms.Frees)
	stat.GaugeMap["BuckHashSys"] = float64(ms.BuckHashSys)
	stat.GaugeMap["GCCPUFraction"] = float64(ms.GCCPUFraction)
	stat.GaugeMap["GCSys"] = float64(ms.GCSys)
	stat.GaugeMap["HeapAlloc"] = float64(ms.HeapAlloc)
	stat.GaugeMap["HeapIdle"] = float64(ms.HeapIdle)
	stat.GaugeMap["HeapInuse"] = float64(ms.HeapInuse)
	stat.GaugeMap["HeapObjects"] = float64(ms.HeapObjects)
	stat.GaugeMap["HeapReleased"] = float64(ms.HeapReleased)
	stat.GaugeMap["HeapSys"] = float64(ms.HeapSys)
	stat.GaugeMap["LastGC"] = float64(ms.LastGC)
	stat.GaugeMap["Lookups"] = float64(ms.Lookups)
	stat.GaugeMap["MCacheInuse"] = float64(ms.MCacheInuse)
	stat.GaugeMap["MCacheSys"] = float64(ms.MCacheSys)
	stat.GaugeMap["MSpanInuse"] = float64(ms.MSpanInuse)
	stat.GaugeMap["MSpanSys"] = float64(ms.MSpanSys)
	stat.GaugeMap["Mallocs"] = float64(ms.Mallocs)
	stat.GaugeMap["NextGC"] = float64(ms.NextGC)
	stat.GaugeMap["NumForcedGC"] = float64(ms.NumForcedGC)
	stat.GaugeMap["NumGC"] = float64(ms.NumGC)
	stat.GaugeMap["OtherSys"] = float64(ms.OtherSys)
	stat.GaugeMap["PauseTotalNs"] = float64(ms.PauseTotalNs)
	stat.GaugeMap["StackInuse"] = float64(ms.StackInuse)
	stat.GaugeMap["StackSys"] = float64(ms.StackSys)
	stat.GaugeMap["Sys"] = float64(ms.Sys)
	stat.GaugeMap["TotalAlloc"] = float64(ms.TotalAlloc)
	stat.GaugeMap["RandomValue"] = rand.Float64() //nolint:gosec //rand хватает
	stat.CounterMap["PollCount"] = counter
}

func SendMetricsJSON(client *resty.Client, baseURL string, stat *Stat) error {
	gaugeMap, counterMap := copyMaps(stat)

	for k, v := range gaugeMap {
		data := common.Metrics{
			ID:    k,
			MType: "gauge",
			Value: &v,
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

	for k, v := range counterMap {
		data := common.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &v,
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
	gaugeMap, counterMap := copyMaps(stat)

	for k, v := range gaugeMap {
		_, err := client.R().
			SetHeader("Content-Type", "Content-Type: text/plain").
			SetPathParams(map[string]string{
				"type":  k,
				"value": strconv.FormatFloat(v, 'f', -1, 64),
			}).
			Post("http://" + baseURL + "/update/gauge/{type}/{value}")
		if err != nil {
			return fmt.Errorf("error executing request: %w", err)
		}
	}

	for k, v := range counterMap {
		_, err := client.R().
			SetHeader("Content-Type", "Content-Type: text/plain").
			SetPathParams(map[string]string{
				"type":  k,
				"value": strconv.FormatInt(v, 10),
			}).
			Post("http://" + baseURL + "/update/counter/{type}/{value}")
		if err != nil {
			return fmt.Errorf("error executing request: %w", err)
		}
	}

	return nil
}

func copyMaps(stat *Stat) (map[string]float64, map[string]int64) {
	gaugeMap := make(map[string]float64)
	counterMap := make(map[string]int64)
	stat.MapLock.RLock()
	for k, v := range stat.GaugeMap {
		gaugeMap[k] = v
	}
	for k, v := range stat.CounterMap {
		counterMap[k] = v
	}
	stat.MapLock.RUnlock()
	return gaugeMap, counterMap
}
