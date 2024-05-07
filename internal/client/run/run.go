package run

import (
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/Alekseyt9/ypmetrics/internal/client/services"
	"github.com/Alekseyt9/ypmetrics/internal/common"
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
		Data: &common.MetricsBatch{
			Counters: make([]common.CounterItem, 1),
			Gauges:   make([]common.GaugeItem, 10),
		},
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
			err := services.SendMetricsBatch(client, cfg.Address, stat)
			if err != nil {
				log.Print(err)
			}
			atomic.StoreInt64(&counter, 0)
			time.Sleep(time.Duration(reportInterval) * time.Second)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
}
