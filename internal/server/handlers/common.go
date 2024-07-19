// Package handlers provides the implementation of request handlers for the server.
package handlers

import (
	"log/slog"

	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
)

// Handler represents a structure to manage the storage and handler settings.
type Handler struct {
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

// NewHandler creates a new Handler with the provided storage and settings.
func NewHandler(store storage.Storage, settings HandlerSettings) *Handler {
	return &Handler{
		store:    store,
		settings: settings,
	}
}

// StoreToFile saves the data to a file if StoreToFileSync is enabled.
func (h *Handler) StoreToFile() {
	if h.settings.StoreToFileSync {
		if memStore, ok := h.store.(*storage.MemStorage); ok {
			err := memStore.SaveToFile(h.settings.FilePath)
			if err != nil {
				h.log.Error("Error save to file", "filepath", h.settings.FilePath)
			}
		}
	}
}
