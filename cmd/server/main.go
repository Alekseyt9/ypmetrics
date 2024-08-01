// Package main contains the main entry point and flag parsing for the server.
package main

import (
	_ "net/http/pprof"

	"github.com/Alekseyt9/ypmetrics/internal/server/config"
	"github.com/Alekseyt9/ypmetrics/internal/server/run"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// main is the entry point for the server application.
// It parses command-line flags and environment variables, then starts the server with the configured settings.
func main() {
	cfg := &config.Config{}
	config.ParseFlags(cfg)
	config.SetEnv(cfg)

	if err := run.Run(cfg); err != nil {
		panic(err)
	}
}
