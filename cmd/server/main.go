package main

import (
	"github.com/Alekseyt9/ypmetrics/internal/server/run"
)

func main() {
	run.ParseFlags()
	run.SetEnv()

	if err := run.Run(); err != nil {
		panic(err)
	}
}
