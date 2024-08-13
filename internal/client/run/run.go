// Package run provides the main execution logic for the client, including configuration and metric reporting.
package run

import (
	"crypto/rsa"
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
	"github.com/Alekseyt9/ypmetrics/internal/common/crypto"
	"github.com/Alekseyt9/ypmetrics/pkg/retry"
	"github.com/Alekseyt9/ypmetrics/pkg/workerpool"
	"github.com/go-resty/resty/v2"
)

// Config holds the configuration settings for the client.
type Config struct {
	HashKey        string `env:"KEY"`             // Key for hashing
	Address        string `env:"ADDRESS"`         // Server address
	ReportInterval int    `env:"REPORT_INTERVAL"` // Interval for reporting metrics
	PollInterval   int    `env:"POLL_INTERVAL"`   // Interval for polling metrics
	RateLimit      int    `env:"RATE_LIMIT"`      // Rate limit for sending metrics
	CryptoKeyFile  string `env:"CRYPTO_KEY"`      // Key for RSA cypering
}

// Run starts the client with the given configuration.
// It initializes metric polling, worker pool, and signal handling for graceful shutdown.
// Parameters:
//   - cfg: the configuration settings for the client
func Run(cfg *Config) {
	var counter int64
	data := initMetricsData()
	startMetricsPolling(data, cfg, &counter)
	workerPool := initWorkerPool(cfg)
	runMetricsSender(cfg, workerPool, data, &counter)
	handleSysSignals(workerPool)
}

// startMetricsPolling begins the polling of metrics at a regular interval.
// Parameters:
//   - data: the metrics data to update
//   - cfg: the configuration settings for the client
//   - counter: a counter for the number of polling iterations
func startMetricsPolling(data *services.MetricsData, cfg *Config, counter *int64) {
	go func() {
		for {
			err := services.UpdateMetrics(data, *counter)
			if err != nil {
				log.Print(err)
			}
			atomic.AddInt64(counter, 1)
			time.Sleep(time.Duration(cfg.PollInterval) * time.Second)
		}
	}()
}

// initMetricsData initializes and returns a new MetricsData instance.
func initMetricsData() *services.MetricsData {
	return services.NewMetricsData()
}

// initWorkerPool initializes and returns a new WorkerPool instance with the given configuration.
// Parameters:
//   - cfg: the configuration settings for the client
func initWorkerPool(cfg *Config) *workerpool.WorkerPool {
	return workerpool.New(cfg.RateLimit)
}

// handleSysSignals sets up signal handling for graceful shutdown of the worker pool.
// Parameters:
//   - wp: the worker pool to close on receiving a shutdown signal
func handleSysSignals(wp *workerpool.WorkerPool) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		wp.Close()
		os.Exit(0)
	}()

	select {}
}

// runMetricsSender starts sending metrics batches at a regular interval using a worker pool.
// Parameters:
//   - cfg: the configuration settings for the client
//   - workerPool: the worker pool to manage metric sending tasks
//   - data: the metrics data to send
//   - counter: a counter for the number of sending iterations
func runMetricsSender(cfg *Config,
	workerPool *workerpool.WorkerPool,
	data *services.MetricsData,
	counter *int64) {
	client := resty.New()
	reportInterval := cfg.ReportInterval

	var pKey *rsa.PublicKey
	var err error
	if cfg.CryptoKeyFile != "" {
		pKey, err = crypto.LoadPublicKey(cfg.CryptoKeyFile)
		if err != nil {
			log.Print(err)
		}
	}

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
			BaseURL:   cfg.Address,
			HashKey:   cfg.HashKey,
			CryptoKey: pKey,
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
