// Package main contains the main entry point and flag parsing for the agent.
package main

import (
	"os"

	"github.com/Alekseyt9/ypmetrics/internal/client/run"
	"github.com/Alekseyt9/ypmetrics/internal/common/version"
)

// Global variables for build information.
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// main is the entry point for the agent application.
// It parses command-line flags and environment variables, then starts the agent with the configured settings.
func main() {
	cfg := &run.Config{}
	ParseFlags(cfg)
	SetEnv(cfg)

	vInfo := version.Info{
		Version: buildVersion,
		Date:    buildDate,
		Commit:  buildCommit,
	}
	vInfo.Print(os.Stdout)

	run.Run(cfg)
}
