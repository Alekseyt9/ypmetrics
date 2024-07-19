// Package main contains the main entry point and flag parsing for the agent.
package main

import (
	"github.com/Alekseyt9/ypmetrics/internal/client/run"
)

// main is the entry point for the agent application.
// It parses command-line flags and environment variables, then starts the agent with the configured settings.
func main() {
	cfg := &run.Config{}
	ParseFlags(cfg)
	SetEnv(cfg)
	run.Run(cfg)
}
