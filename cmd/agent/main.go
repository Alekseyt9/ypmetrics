package main

import (
	"time"

	services "github.com/Alekseyt9/ypmetrics/internal/client/services"
	"github.com/go-resty/resty/v2"
)

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
			services.UpdateMetrics(gMap, cMap, counter)
			counter = counter + 1
		}

		if interval%int64(reportInterval) == 0 {
			services.SendMetrics(client, gMap, cMap)
			counter = 0
		}

		interval = interval + 1
		time.Sleep(1 * time.Second)
	}
}
