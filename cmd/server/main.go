// Package main contains the main entry point and flag parsing for the server.
package main

import (
	"log"
	_ "net/http/pprof"
	"os"

	"github.com/Alekseyt9/ypmetrics/internal/common/version"
	"github.com/Alekseyt9/ypmetrics/internal/server/config"
	"github.com/Alekseyt9/ypmetrics/internal/server/run"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Global variables for build information.
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// main is the entry point for the server application.
// It parses command-line flags and environment variables, then starts the server with the configured settings.
func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	vInfo := version.Info{
		Version: buildVersion,
		Date:    buildDate,
		Commit:  buildCommit,
	}
	vInfo.Print(os.Stdout)

	if err := run.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
