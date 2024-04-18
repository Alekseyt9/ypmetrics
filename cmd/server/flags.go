package main

import (
	goflag "flag"

	"github.com/Alekseyt9/ypmetrics/internal/server/run"
	"github.com/caarlos0/env"
	flag "github.com/spf13/pflag"
)

func ParseFlags(cfg *run.Config) {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	address := flag.StringP("address", "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	cfg.Address = *address
}

func SetEnv(cfg *run.Config) {
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
}
