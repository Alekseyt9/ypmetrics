package main

import (
	"fmt"
	"net/http"

	handlers "github.com/Alekseyt9/ypmetrics/internal/server/handlers"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

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

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}
