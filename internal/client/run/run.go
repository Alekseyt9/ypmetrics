package run

import (
	"time"

	"github.com/Alekseyt9/ypmetrics/internal/client/services"
	"github.com/go-resty/resty/v2"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval *int   `env:"REPORT_INTERVAL"`
	PollInterval   *int   `env:"POLL_INTERVAL"`
}

func Run(cfg *Config) {
	pollInterval := *cfg.PollInterval
	reportInterval := *cfg.ReportInterval

	var interval int64
	var counter int64
	gMap := make(map[string]float64)
	cMap := make(map[string]int64)
	client := resty.New()

	for {
		if interval%int64(pollInterval) == 0 {
			services.UpdateMetrics(gMap, cMap, counter)
			counter++
		}

		if interval%int64(reportInterval) == 0 {
			services.SendMetrics(client, cfg.Address, gMap, cMap)
			counter = 0
		}

		interval++
		time.Sleep(1 * time.Second)
	}
}
