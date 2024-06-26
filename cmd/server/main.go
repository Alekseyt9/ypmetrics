package main

import (
	"github.com/Alekseyt9/ypmetrics/internal/server/run"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := &run.Config{}
	ParseFlags(cfg)
	SetEnv(cfg)

	if err := run.Run(cfg); err != nil {
		panic(err)
	}
}
