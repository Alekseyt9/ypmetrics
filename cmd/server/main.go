package main

import (
	"fmt"
	"net/http"

	handlers "github.com/Alekseyt9/ypmetrics/internal/server/handlers"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func Router(store storage.Storage) chi.Router {
	r := chi.NewRouter()

	r.Route("/update", func(r chi.Router) {
		r.Post("/*", handlers.HandleIncorrectType)
		r.Route("/gauge", func(r chi.Router) {
			r.Post("/", handlers.HandleNotValue)
			r.Route("/{name}", func(r chi.Router) {
				r.Post("/{value}", handlers.HandleGauge(store))
			})
		})
		r.Route("/counter", func(r chi.Router) {
			r.Post("/", handlers.HandleNotValue)
			r.Route("/{name}", func(r chi.Router) {
				r.Post("/{value}", handlers.HandleCounter(store))
			})
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/gauge/{name}", handlers.HandleGetGauge(store))
		r.Get("/counter/{name}", handlers.HandleGetCounter(store))
	})

	r.Get("/", handlers.HandleGetAll(store))

	return r
}

func run() error {
	store := storage.NewMemStorage()
	r := Router(store)

	fmt.Println("Running server on ", *flagAddr)
	return http.ListenAndServe(*flagAddr, r)
}

func setEnv() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	if cfg.Address != "" {
		*flagAddr = cfg.Address
	}
}

func main() {
	parseFlags()
	setEnv()

	if err := run(); err != nil {
		panic(err)
	}
}
