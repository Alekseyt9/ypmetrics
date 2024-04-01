package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

func updateMetrics(gMap map[string]float64, cMap map[string]int64, counter int64) {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)
	gMap["Alloc"] = float64(ms.Alloc)
	gMap["BuckHashSys"] = float64(ms.BuckHashSys)
	gMap["Frees"] = float64(ms.Frees)
	gMap["BuckHashSys"] = float64(ms.BuckHashSys)
	gMap["GCCPUFraction"] = float64(ms.GCCPUFraction)
	gMap["GCSys"] = float64(ms.GCSys)
	gMap["HeapAlloc"] = float64(ms.HeapAlloc)
	gMap["HeapIdle"] = float64(ms.HeapIdle)
	gMap["HeapInuse"] = float64(ms.HeapInuse)
	gMap["HeapObjects"] = float64(ms.HeapObjects)
	gMap["HeapReleased"] = float64(ms.HeapReleased)
	gMap["HeapSys"] = float64(ms.HeapSys)
	gMap["LastGC"] = float64(ms.LastGC)
	gMap["Lookups"] = float64(ms.Lookups)
	gMap["MCacheInuse"] = float64(ms.MCacheInuse)
	gMap["MCacheSys"] = float64(ms.MCacheSys)
	gMap["MSpanInuse"] = float64(ms.MSpanInuse)
	gMap["MSpanSys"] = float64(ms.MSpanSys)
	gMap["Mallocs"] = float64(ms.Mallocs)
	gMap["NextGC"] = float64(ms.NextGC)
	gMap["NumForcedGC"] = float64(ms.NumForcedGC)
	gMap["NumGC"] = float64(ms.NumGC)
	gMap["OtherSys"] = float64(ms.OtherSys)
	gMap["PauseTotalNs"] = float64(ms.PauseTotalNs)
	gMap["StackInuse"] = float64(ms.StackInuse)
	gMap["StackSys"] = float64(ms.StackSys)
	gMap["Sys"] = float64(ms.Sys)
	gMap["TotalAlloc"] = float64(ms.TotalAlloc)
	gMap["RandomValue"] = rand.Float64()
	cMap["PollCount"] = int64(counter)

	fmt.Println("updateMetrics")
}

func sendMetrics(client *resty.Client, gMap map[string]float64, cMap map[string]int64) {
	for k, v := range gMap {
		client.R().
			SetHeader("Content-Type", "Content-Type: text/plain").
			SetPathParams(map[string]string{
				"type":  k,
				"value": strconv.FormatFloat(v, 'f', -1, 64),
			}).
			Post("http://localhost:8080/update/gauge/{type}/{value}")
	}

	for k, v := range cMap {
		client.R().
			SetHeader("Content-Type", "Content-Type: text/plain").
			SetPathParams(map[string]string{
				"type":  k,
				"value": strconv.FormatInt(v, 10),
			}).
			Post("http://localhost:8080/update/counter/{type}/{value}")
	}

	fmt.Println("sendMetrics")
}

func main() {
	pollInterval := 2
	reportInterval := 10
	var interval int64 = 0
	var counter int64 = 0
	gMap := make(map[string]float64)
	cMap := make(map[string]int64)
	client := resty.New()

	for {
		if interval%int64(pollInterval) == 0 {
			updateMetrics(gMap, cMap, counter)
			counter = counter + 1
		}

		if interval%int64(reportInterval) == 0 {
			sendMetrics(client, gMap, cMap)
		}

		interval = interval + 1
		time.Sleep(1 * time.Second)
	}
}
