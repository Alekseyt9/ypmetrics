package run

import (
	"net/http"
	"time"

	handlers "github.com/Alekseyt9/ypmetrics/internal/server/handlers"
	"github.com/Alekseyt9/ypmetrics/internal/server/logger"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 15 * time.Second
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func Router(store storage.Storage, log logger.Logger) chi.Router {
	h := handlers.NewHandler(store)

	r := chi.NewRouter()

	// тк использую библиотеку chi - подключаю middleware стандартным способом
	r.Use(func(next http.Handler) http.Handler {
		return logger.WithLogging(next, log)
	})

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

func Run(cfg *Config) error {
	store := storage.NewMemStorage()
	logger := logger.NewLogger()

	r := Router(store, logger)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	logger.Infof("Running server on %s", cfg.Address)
	return server.ListenAndServe()
}
