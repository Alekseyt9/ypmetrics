package main

import (
	goflag "flag"

	"github.com/Alekseyt9/ypmetrics/internal/client/run"
	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

func ParseFlags(cfg *run.Config) {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	address := flag.StringP("address", "a", "localhost:8080", "address and port to run server")
	reportInterval := flag.IntP("reportInterval", "r", 10, "frequency of sending metrics to the server")
	pollInterval := flag.IntP("pollInterval", "p", 2, "frequency of polling metrics from the runtime package")
	key := flag.StringP("key", "k", "", "key for SHA256 signing")
	rateLimit := flag.IntP("rateLimit", "l", 5, "upper limit on the number of outgoing requests")

	flag.Parse()

	cfg.Address = *address
	cfg.ReportInterval = *reportInterval
	cfg.PollInterval = *pollInterval
	cfg.HashKey = *key
	cfg.RateLimit = *rateLimit
}

func SetEnv(cfg *run.Config) {
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
}
