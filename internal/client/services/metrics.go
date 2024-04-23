package services

import (
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"sync"

	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
)

type Stat struct {
	CounterMap  map[string]int64
	CounterLock sync.RWMutex
	GaugeMap    map[string]float64
	GaugeLock   sync.RWMutex
}

func UpdateMetrics(stat *Stat, counter int64) {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	stat.GaugeLock.Lock()
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
	stat.GaugeLock.Unlock()

	stat.CounterLock.Lock()
	stat.CounterMap["PollCount"] = counter
	stat.CounterLock.Unlock()
}

func SendMetricsJSON(client *resty.Client, baseURL string, stat *Stat) {
	stat.GaugeLock.RLock()
	for k, v := range stat.GaugeMap {
		data := common.Metrics{
			ID:    k,
			MType: "gauge",
			Value: &v,
		}
		out, err := easyjson.Marshal(data)
		if err != nil {
			panic(err)
		}

		compressedOut, err := common.BrotliCompress(out)
		if err != nil {
			log.Printf("Ошибка при сжатии: %v", err)
		}

		_, err = client.R().
			SetHeader("Content-Type", "Content-Type: application/json").
			SetHeader("Content-Encoding", "br").
			SetHeader("Accept-Encoding", "br").
			SetBody(compressedOut).
			Post("http://" + baseURL + "/update")

		if err != nil {
			log.Printf("Ошибка при выполнении запроса: %v", err)
		}
	}
	stat.GaugeLock.RUnlock()

	stat.CounterLock.RLock()
	for k, v := range stat.CounterMap {
		data := common.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &v,
		}
		out, err := easyjson.Marshal(data)
		if err != nil {
			panic(err)
		}

		compressedOut, err := common.BrotliCompress(out)
		if err != nil {
			log.Printf("Ошибка при сжатии: %v", err)
		}

		_, err = client.R().
			SetHeader("Content-Type", "Content-Type: application/json").
			SetHeader("Content-Encoding", "br").
			SetHeader("Accept-Encoding", "br").
			SetBody(compressedOut).
			Post("http://" + baseURL + "/update")
		if err != nil {
			log.Printf("Ошибка при выполнении запроса: %v", err)
		}
	}
	stat.CounterLock.RUnlock()
}

func SendMetricsURL(client *resty.Client, baseURL string, stat *Stat) {
	stat.GaugeLock.RLock()
	for k, v := range stat.GaugeMap {
		_, err := client.R().
			SetHeader("Content-Type", "Content-Type: text/plain").
			SetPathParams(map[string]string{
				"type":  k,
				"value": strconv.FormatFloat(v, 'f', -1, 64),
			}).
			Post("http://" + baseURL + "/update/gauge/{type}/{value}")
		if err != nil {
			log.Printf("Ошибка при выполнении запроса: %v", err)
		}
	}
	stat.GaugeLock.RUnlock()

	stat.CounterLock.RLock()
	for k, v := range stat.CounterMap {
		_, err := client.R().
			SetHeader("Content-Type", "Content-Type: text/plain").
			SetPathParams(map[string]string{
				"type":  k,
				"value": strconv.FormatInt(v, 10),
			}).
			Post("http://" + baseURL + "/update/counter/{type}/{value}")
		if err != nil {
			log.Printf("Ошибка при выполнении запроса: %v", err)
		}
	}
	stat.CounterLock.RUnlock()
}
