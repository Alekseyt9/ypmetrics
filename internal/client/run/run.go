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

	data := services.NewMetricsData()

	go func() {
		for {
			err := services.UpdateMetrics(data, counter)
			if err != nil {
				log.Print(err)
			}
			atomic.AddInt64(&counter, 1)
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}()

	workerPool := workerpool.New(cfg.RateLimit)
	runSend(cfg, workerPool, data, &counter)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		workerPool.Close()
		os.Exit(0)
	}()

	select {}
}

func runSend(cfg *Config, workerPool *workerpool.WorkerPool, data *services.MetricsData, counter *int64) {
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
			// Отправляем разные наборы метрик
			workerPool.AddTask(func() {
				err := retryCtr.Do(func() error {
					return services.SendMetricsBatch(client, data.StatRuntime, sendOpts)
				})
				if err != nil {
					log.Print(err)
				}
				atomic.StoreInt64(counter, 0)
			})

			workerPool.AddTask(func() {
				err := retryCtr.Do(func() error {
					return services.SendMetricsBatch(client, data.StatGopsutil, sendOpts)
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
