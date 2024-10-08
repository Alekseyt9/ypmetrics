// Package handlers provides the implementation of request handlers for the server.
package handlers

import (
	"github.com/Alekseyt9/ypmetrics/internal/server/log"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
)

// MetricsHandler represents a structure to manage the storage and handler settings.
type MetricsHandler struct {
	log      log.Logger
	store    storage.Storage
	settings HandlerSettings
}

// HandlerSettings contains settings for the handler.
type HandlerSettings struct {
	SaveFile        string // Path to the file for saving data.
	DatabaseDSN     string // Database connection string.
	HashKey         string // Key for hashing data.
	StoreToFileSync bool   // Indicates whether to synchronize data to a file.
}

// NewMetricsHandler creates a new Handler with the provided storage and settings.
func NewMetricsHandler(store storage.Storage, settings HandlerSettings, log log.Logger) *MetricsHandler {
	return &MetricsHandler{
		store:    store,
		settings: settings,
		log:      log,
	}
}

// StoreToFile saves the data to a file if StoreToFileSync is enabled.
func (h *MetricsHandler) StoreToFile() {
	if h.settings.StoreToFileSync {
		if memStore, ok := h.store.(*storage.MemStorage); ok {
			err := memStore.SaveToFile(h.settings.SaveFile)
			if err != nil {
				h.log.Error("Error save to file", "filepath", h.settings.SaveFile)
			}
		}
	}
}
