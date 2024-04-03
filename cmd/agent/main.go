package main

import (
	"github.com/Alekseyt9/ypmetrics/internal/client/run"
)

func main() {
	run.ParseFlags()
	run.SetEnv()
	run.Run()
}
