package run

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/Alekseyt9/ypmetrics/internal/client/services"
	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/Alekseyt9/ypmetrics/pkg/retry"
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
		Data: &common.MetricItems{
			Counters: make([]common.CounterItem, 0),
			Gauges:   make([]common.GaugeItem, 0),
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
			retryCtr := retry.NewControllerStd(func(err error) bool {
				var netErr net.Error
				if (errors.As(err, &netErr) && netErr.Timeout()) ||
					strings.Contains(err.Error(), "EOF") ||
					strings.Contains(err.Error(), "connection reset by peer") {
					return true
				}
				return false
			})
			err := retryCtr.Do(func() error {
				return services.SendMetricsBatch(client, cfg.Address, stat)
			})
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
