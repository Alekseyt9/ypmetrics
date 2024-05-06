package main

import (
	goflag "flag"

	"github.com/Alekseyt9/ypmetrics/internal/server/run"
	"github.com/caarlos0/env"
	flag "github.com/spf13/pflag"
)

func ParseFlags(cfg *run.Config) {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	address := flag.StringP("address", "a", "localhost:8080", "Address and port to run server")
	storeInterval := flag.IntP("store-interval", "i", 300,
		"Time interval in seconds, based on which the current state of the server is displayed on disk")
	fileStoragePath := flag.StringP("file-storage-path", "f", "/tmp/metrics-db.json",
		"Full name of the file where the current values are saved")
	restore := flag.BoolP("restore", "r", true,
		"Load previously saved values from the file when the server starts")
	flag.Parse()

	cfg.Address = *address
	cfg.StoreInterval = *storeInterval
	cfg.FileStoragePath = *fileStoragePath
	cfg.Restore = *restore
}

func SetEnv(cfg *run.Config) {
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
}
