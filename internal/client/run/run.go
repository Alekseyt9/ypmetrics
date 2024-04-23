package run

import (
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/Alekseyt9/ypmetrics/internal/client/services"
	"github.com/go-resty/resty/v2"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func Run(cfg *Config) {
	pollInterval := cfg.PollInterval
	reportInterval := cfg.ReportInterval
	var counter int64

	stat := &services.Stat{
		CounterMap: make(map[string]int64),
		GaugeMap:   make(map[string]float64),
	}
	client := resty.New()

	go func() {
		for {
			services.UpdateMetrics(stat, counter)
			atomic.AddInt64(&counter, 1)
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			services.SendMetricsJSON(client, cfg.Address, stat)
			atomic.StoreInt64(&counter, 0)
			time.Sleep(time.Duration(reportInterval) * time.Second)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
}
