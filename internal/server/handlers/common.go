// Package handlers provides the implementation of request handlers for the server.
package handlers

import (
	"log/slog"

	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
)

// MetricsHandler represents a structure to manage the storage and handler settings.
type MetricsHandler struct {
	store    storage.Storage
	settings HandlerSettings
	log      *slog.Logger
}

// HandlerSettings contains settings for the handler.
type HandlerSettings struct {
	StoreToFileSync bool   // Indicates whether to synchronize data to a file.
	FilePath        string // Path to the file for saving data.
	DatabaseDSN     string // Database connection string.
	HashKey         string // Key for hashing data.
}

// NewMetricsHandler creates a new Handler with the provided storage and settings.
func NewMetricsHandler(store storage.Storage, settings HandlerSettings) *MetricsHandler {
	return &MetricsHandler{
		store:    store,
		settings: settings,
	}
}

// StoreToFile saves the data to a file if StoreToFileSync is enabled.
func (h *MetricsHandler) StoreToFile() {
	if h.settings.StoreToFileSync {
		if memStore, ok := h.store.(*storage.MemStorage); ok {
			err := memStore.SaveToFile(h.settings.FilePath)
			if err != nil {
				h.log.Error("Error save to file", "filepath", h.settings.FilePath)
			}
		}
	}
}
