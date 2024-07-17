package handlers

import (
	"log/slog"

	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
)

type Handler struct {
	store    storage.Storage
	settings HandlerSettings
	log      *slog.Logger
}

type HandlerSettings struct {
	StoreToFileSync bool
	FilePath        string
	DatabaseDSN     string
	HashKey         string
}

func NewHandler(store storage.Storage, settings HandlerSettings) *Handler {
	return &Handler{
		store:    store,
		settings: settings,
	}
}

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
