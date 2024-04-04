package main

import (
	goflag "flag"

	"github.com/Alekseyt9/ypmetrics/internal/client/run"
	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

var FlagAddr *string = flag.StringP("address", "a", "localhost:8080", "address and port to connect to server")
var FlagReportInterval *int = flag.IntP("reportInterval", "r", 10, "frequency of sending metrics to the server")
var FlagPollInterval *int = flag.IntP("pollInterval", "p", 2, "frequency of polling metrics from the runtime package")

func ParseFlags() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
}

func SetEnv() {
	var cfg run.Config
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
