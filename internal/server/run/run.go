package run

import (
	"context"
	"database/sql"
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
	_ "github.com/jackc/pgx/v5/stdlib"
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
}

func Router(store storage.Storage, log logger.Logger, cfg *Config) chi.Router {
	hs := handlers.HandlerSettings{
		DatabaseDSN: cfg.DataBaseDSN,
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
		db, err := sql.Open("pgx", cfg.DataBaseDSN)
		if err != nil {
			return err
		}

		dbstore := storage.NewDBRetryStorage(storage.NewDBStorage(db))
		//dbstore := storage.NewDBStorage(db)
		err = dbstore.Bootstrap(context.Background())
		if err != nil {
			logger.Error("Database has already initialized")
		}
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

	logger.Info("Running server on ", "address", cfg.Address)
	return server.ListenAndServe()
}
