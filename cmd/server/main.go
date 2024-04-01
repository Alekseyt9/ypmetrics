package main

import (
	"log"
	"net/http"

	handlers "github.com/Alekseyt9/ypmetrics/internal/server/handlers"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func main() {

	store := storage.NewMemStorage()

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

	log.Fatal(http.ListenAndServe(":8080", r))
}
