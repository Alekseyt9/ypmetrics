package run

import (
	"time"

	"github.com/Alekseyt9/ypmetrics/internal/client/services"
	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval *int   `env:"REPORT_INTERVAL"`
	PollInterval   *int   `env:"POLL_INTERVAL"`
}

func Run() {
	pollInterval := *FlagPollInterval
	reportInterval := *FlagReportInterval

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
			services.SendMetrics(client, *FlagAddr, gMap, cMap)
			counter = 0
		}

		interval = interval + 1
		time.Sleep(1 * time.Second)
	}
}

func SetEnv() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	if cfg.Address != "" {
		*FlagAddr = cfg.Address
	}
	if cfg.PollInterval != nil {
		*FlagPollInterval = *cfg.PollInterval
	}
	if cfg.ReportInterval != nil {
		*FlagReportInterval = *cfg.ReportInterval
	}
}
