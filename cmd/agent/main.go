package main

import (
	"github.com/Alekseyt9/ypmetrics/internal/client/run"
)

func main() {
	cfg := &run.Config{}
	ParseFlags(cfg)
	SetEnv(cfg)
	run.Run(cfg)
}
