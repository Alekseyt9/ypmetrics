// Package run provides the main execution logic for the server, including configuration and routing.
package run

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cr "github.com/Alekseyt9/ypmetrics/internal/common/crypto"
	"github.com/Alekseyt9/ypmetrics/internal/server/config"
	"github.com/Alekseyt9/ypmetrics/internal/server/handlers"
	"github.com/Alekseyt9/ypmetrics/internal/server/log"
	"github.com/Alekseyt9/ypmetrics/internal/server/middleware/compress"
	"github.com/Alekseyt9/ypmetrics/internal/server/middleware/crypto"
	"github.com/Alekseyt9/ypmetrics/internal/server/middleware/hash"
	"github.com/Alekseyt9/ypmetrics/internal/server/middleware/ipcheck"
	"github.com/Alekseyt9/ypmetrics/internal/server/middleware/logger"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	readTimeout  = time.Minute
	writeTimeout = time.Minute
	idleTimeout  = time.Minute
)

// Router sets up the router with the necessary routes and middleware.
// Parameters:
//   - store: the storage to use for handling metrics
//   - log: the logger instance
//   - cfg: the configuration settings for the server
//
// Returns a chi.Router with the configured routes and middleware.
func Router(store storage.Storage, log log.Logger, cfg *config.Config) chi.Router {
	h := initHandler(store, cfg, log)
	r := chi.NewRouter()

	setupMiddleware(r, log, cfg)
	setupUpdateRoutes(r, h)
	setupValueRoutes(r, h)
	setupOtherRoutes(r, h)

	return r
}

// setupOtherRoutes configures additional routes for the router.
// Parameters:
//   - r: the chi router to configure
//   - h: the metrics handler to handle the routes
func setupOtherRoutes(r *chi.Mux, h *handlers.MetricsHandler) {
	r.Get("/ping", h.HandlePing)
	r.Get("/", h.HandleGetAll)
}

// setupValueRoutes configures value routes for the router.
// Parameters:
//   - r: the chi router to configure
//   - h: the metrics handler to handle the routes
func setupValueRoutes(r *chi.Mux, h *handlers.MetricsHandler) {
	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.HandleValueJSON)
		r.Get("/gauge/{name}", h.HandleGetGauge)
		r.Get("/counter/{name}", h.HandleGetCounter)
	})
}

// setupUpdateRoutes configures update routes for the router.
// Parameters:
//   - r: the chi router to configure
//   - h: the metrics handler to handle the routes
func setupUpdateRoutes(r *chi.Mux, h *handlers.MetricsHandler) {
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
}

// initHandler initializes the metrics handler with the given storage and configuration.
// Parameters:
//   - store: the storage to use for handling metrics
//   - cfg: the configuration settings for the server
//
// Returns a MetricsHandler instance.
func initHandler(store storage.Storage, cfg *config.Config, log log.Logger) *handlers.MetricsHandler {
	hs := handlers.HandlerSettings{}
	if cfg.DataBaseDSN != nil {
		hs.DatabaseDSN = *cfg.DataBaseDSN
	}
	if cfg.HashKey != nil {
		hs.HashKey = *cfg.HashKey
	}

	if cfg.StoreInterval == nil {
		hs.StoreToFileSync = true
		hs.SaveFile = *cfg.FileStoragePath
	}

	return handlers.NewMetricsHandler(store, hs, log)
}

// setupMiddleware configures middleware for the router.
// Parameters:
//   - r: the chi router to configure
//   - log: the logger instance
//   - cfg: the configuration settings for the server
func setupMiddleware(r *chi.Mux, log log.Logger, cfg *config.Config) {
	r.Use(func(next http.Handler) http.Handler {
		return logger.WithLogging(next, log)
	})

	if cfg.TrustedSubnet != nil {
		r.Use(func(next http.Handler) http.Handler {
			return ipcheck.WithIPCheck(next, *cfg.TrustedSubnet)
		})
	}

	if cfg.HashKey != nil {
		r.Use(func(next http.Handler) http.Handler {
			return hash.WithHash(next, *cfg.HashKey)
		})
	}

	if cfg.CryptoKeyFile != nil {
		key, err := cr.LoadPrivateKey(*cfg.CryptoKeyFile)
		if err != nil {
			log.Error("LoadPrivateKey", err)
		}
		r.Use(func(next http.Handler) http.Handler {
			return crypto.WithCrypto(next, log, key)
		})
	}

	r.Use(func(next http.Handler) http.Handler {
		return compress.WithCompress(next, log)
	})

	r.Mount("/debug", middleware.Profiler())
}

// Run starts the server with the given configuration.
// It initializes the storage, restores data if necessary, and starts the server.
// Parameters:
//   - cfg: the configuration settings for the server
//
// Returns an error if the server fails to start.
func Run(cfg *config.Config) error {
	var store storage.Storage
	logger := log.NewSlogLogger()

	if cfg.DataBaseDSN != nil {
		innerStore, err := storage.NewDBStorage(*cfg.DataBaseDSN)
		if err != nil {
			return err
		}

		dbstore := storage.NewDBRetryStorage(innerStore)
		store = dbstore
	} else {
		store = storage.NewMemStorage()
	}

	if *cfg.Restore && cfg.FileStoragePath != nil {
		if memStore, ok := store.(*storage.MemStorage); ok {
			err := memStore.LoadFromFile(*cfg.FileStoragePath)
			if err != nil {
				logger.Info("Load from file", "error", err)
			}
		}
	}

	err := serverStart(store, cfg, logger)
	if err != nil {
		return err
	}

	err = grpcServerStart(store, cfg, logger)
	if err != nil {
		return err
	}

	return nil
}

func serverStart(store storage.Storage, cfg *config.Config, logger log.Logger) error {
	r := Router(store, logger, cfg)

	server := &http.Server{
		Addr:         *cfg.Address,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	if *cfg.StoreInterval > 0 {
		go func() {
			for {
				if memStore, ok := store.(*storage.MemStorage); ok {
					if err := memStore.SaveToFile(*cfg.FileStoragePath); err != nil {
						logger.Info("Save to dump", "error", err)
					}
					time.Sleep(time.Duration(*cfg.StoreInterval) * time.Second)
				}
			}
		}()
	}

	finalize(store, server, cfg, logger)

	logger.Info("Running server on ", "address", cfg.Address)
	return server.ListenAndServe()
}

func finalize(store storage.Storage, server *http.Server, cfg *config.Config, logger log.Logger) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigs

		if memStore, ok := store.(*storage.MemStorage); ok {
			if err := memStore.SaveToFile(*cfg.FileStoragePath); err != nil {
				logger.Info("Save to dump", "error", err)
			}
		}

		if err := server.Shutdown(context.Background()); err != nil {
			logger.Info("Server shutdown", "error", err)
		}
	}()
}
