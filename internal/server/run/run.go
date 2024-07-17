package run

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlers "github.com/Alekseyt9/ypmetrics/internal/server/handlers"
	"github.com/Alekseyt9/ypmetrics/internal/server/middleware/compress"
	"github.com/Alekseyt9/ypmetrics/internal/server/middleware/hash"
	"github.com/Alekseyt9/ypmetrics/internal/server/middleware/logger"
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
	DataBaseDSN     string `env:"DATABASE_DSN"`
	HashKey         string `env:"KEY"`
}

func Router(store storage.Storage, log *slog.Logger, cfg *Config) chi.Router {
	hs := handlers.HandlerSettings{
		DatabaseDSN: cfg.DataBaseDSN,
		HashKey:     cfg.HashKey,
	}
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
	r.Use(func(next http.Handler) http.Handler {
		return hash.WithHash(next, cfg.HashKey)
	})

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
	r.Post("/updates/", h.HandleUpdateBatchJSON)

	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.HandleValueJSON)
		r.Get("/gauge/{name}", h.HandleGetGauge)
		r.Get("/counter/{name}", h.HandleGetCounter)
	})

	r.Get("/ping", h.HandlePing)
	r.Get("/", h.HandleGetAll)

	return r
}

func Run(cfg *Config) error {
	var store storage.Storage
	logger := logger.NewSlogLogger()

	if cfg.DataBaseDSN != "" {
		innerStore, err := storage.NewDBStorage(cfg.DataBaseDSN)
		if err != nil {
			return err
		}

		dbstore := storage.NewDBRetryStorage(innerStore)
		store = dbstore
	} else {
		store = storage.NewMemStorage()
	}

	if cfg.Restore && cfg.FileStoragePath != "" {
		if memStore, ok := store.(*storage.MemStorage); ok {
			err := memStore.LoadFromFile(cfg.FileStoragePath)
			if err != nil {
				logger.Info("Load from file", "error", err)
			}
		}
	}

	return serverStart(store, cfg, logger)
}

func serverStart(store storage.Storage, cfg *Config, logger *slog.Logger) error {
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
				if memStore, ok := store.(*storage.MemStorage); ok {
					if err := memStore.SaveToFile(cfg.FileStoragePath); err != nil {
						logger.Info("Save to dump", "error", err)
					}
					time.Sleep(time.Duration(cfg.StoreInterval) * time.Second)
				}
			}
		}()
	}

	finalize(store, server, cfg, logger)

	logger.Info("Running server on ", "address", cfg.Address)
	return server.ListenAndServe()
}

func finalize(store storage.Storage, server *http.Server, cfg *Config, logger *slog.Logger) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs

		if memStore, ok := store.(*storage.MemStorage); ok {
			if err := memStore.SaveToFile(cfg.FileStoragePath); err != nil {
				logger.Info("Save to dump", "error", err)
			}
		}

		if err := server.Shutdown(context.Background()); err != nil {
			logger.Info("Server shutdown", "error", err)
		}
	}()
}
