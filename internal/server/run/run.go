package run

import (
	"fmt"
	"net/http"

	handlers "github.com/Alekseyt9/ypmetrics/internal/server/handlers"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func Router(store storage.Storage) chi.Router {
	h := handlers.NewHandler(store)

	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/*", h.HandleIncorrectType)
		r.Route("/gauge", func(r chi.Router) {
			r.Post("/", h.HandleNotValue)
			r.Route("/{name}", func(r chi.Router) {
				r.Post("/{value}", h.HandleGauge)
			})
		})
		r.Route("/counter", func(r chi.Router) {
			r.Post("/", h.HandleNotValue)
			r.Route("/{name}", func(r chi.Router) {
				r.Post("/{value}", h.HandleCounter)
			})
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/gauge/{name}", h.HandleGetGauge)
		r.Get("/counter/{name}", h.HandleGetCounter)
	})

	r.Get("/", h.HandleGetAll)

	return r
}

func Run() error {
	store := storage.NewMemStorage()
	r := Router(store)

	fmt.Println("Running server on ", *FlagAddr)
	return http.ListenAndServe(*FlagAddr, r)
}

func SetEnv() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	if cfg.Address != "" {
		*FlagAddr = cfg.Address
	}
}
