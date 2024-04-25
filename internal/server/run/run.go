package run

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alekseyt9/ypmetrics/internal/server/compress"
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
	Address         string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func Router(store storage.Storage, log logger.Logger, cfg *Config) chi.Router {
	hs := handlers.HandlerSettings{}
	if cfg.StoreInterval == 0 {
		hs.StoreToFileSync = true
		hs.FilePath = cfg.FileStoragePath
	}

	h := handlers.NewHandler(store, hs)
	r := chi.NewRouter()

	// тк использую библиотеку chi - подключаю middleware стандартным способом
	r.Use(func(next http.Handler) http.Handler {
		return logger.WithLogging(next, log)
	})
	r.Use(compress.WithCompress)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.HandleUpdateJSON)
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
		r.Post("/", h.HandleValueJSON)
		r.Get("/gauge/{name}", h.HandleGetGauge)
		r.Get("/counter/{name}", h.HandleGetCounter)
	})

	r.Get("/", h.HandleGetAll)

	return r
}

func Run(cfg *Config) error {
	store := storage.NewMemStorage()
	logger := logger.NewSlogLogger()

	if cfg.Restore && cfg.FileStoragePath != "" {
		err := store.LoadFromFile(cfg.FileStoragePath)
		if err != nil {
			logger.Info("Load from file", "error", err)
		}
	}

	r := Router(store, logger, cfg)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	if cfg.StoreInterval > 0 {
		go func() {
			for {
				if err := store.SaveToFile(cfg.FileStoragePath); err != nil {
					logger.Info("Save to dump", "error", err)
				}
				time.Sleep(time.Duration(cfg.StoreInterval) * time.Second)
			}
		}()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs

		if err := store.SaveToFile(cfg.FileStoragePath); err != nil {
			logger.Info("Save to dump", "error", err)
		}

		if err := server.Shutdown(context.Background()); err != nil {
			logger.Info("Server shutdown", "error", err)
		}
	}()

	logger.Info("Running server on ", "address", cfg.Address)
	return server.ListenAndServe()
}
