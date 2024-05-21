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
	"github.com/Alekseyt9/ypmetrics/pkg/workerpool"
	"github.com/go-resty/resty/v2"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	HashKey        string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func Run(cfg *Config) {
	pollInterval := cfg.PollInterval
	var counter int64

	stat := &services.Stat{
		Data: &common.MetricItems{
			Counters: make([]common.CounterItem, 0),
			Gauges:   make([]common.GaugeItem, 0),
		},
	}

	go func() {
		for {
			err := services.UpdateMetrics1(stat, counter)
			if err != nil {
				log.Print(err)
			}
			atomic.AddInt64(&counter, 1)
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			err := services.UpdateMetrics2(stat, counter)
			if err != nil {
				log.Print(err)
			}
			atomic.AddInt64(&counter, 1)
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}()

	workerPool := workerpool.New(cfg.RateLimit)
	workerPool.Run()

	runSend(cfg, workerPool, stat, &counter)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		workerPool.Close()
		os.Exit(0)
	}()

	select {}
}

func runSend(cfg *Config, workerPool *workerpool.WorkerPool, stat *services.Stat, counter *int64) {
	client := resty.New()
	reportInterval := cfg.ReportInterval

	go func() {
		retryCtr := retry.NewControllerStd(func(err error) bool {
			var netErr net.Error
			if (errors.As(err, &netErr) && netErr.Timeout()) ||
				strings.Contains(err.Error(), "EOF") ||
				strings.Contains(err.Error(), "connection reset by peer") {
				return true
			}
			return false
		})
		sendOpts := &services.SendOptions{
			BaseURL: cfg.Address,
			HashKey: cfg.HashKey,
		}

		for {
			workerPool.AddTask(func() {
				err := retryCtr.Do(func() error {
					return services.SendMetricsBatch(client, stat, sendOpts)
				})
				if err != nil {
					log.Print(err)
				}
				atomic.StoreInt64(counter, 0)
			})

			time.Sleep(time.Duration(reportInterval) * time.Second)
		}
	}()
}
