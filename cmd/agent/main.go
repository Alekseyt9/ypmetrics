package main

import (
	"github.com/Alekseyt9/ypmetrics/internal/client/run"
)

func main() {
	cfg := ParseFlags()
	SetEnv(cfg)
	run.Run(cfg)
}
