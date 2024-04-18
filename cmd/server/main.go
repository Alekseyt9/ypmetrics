package main

import (
	"github.com/Alekseyt9/ypmetrics/internal/server/run"
)

func main() {
	cfg := &run.Config{}
	ParseFlags(cfg)
	SetEnv(cfg)

	if err := run.Run(cfg); err != nil {
		panic(err)
	}
}
