package handlers

import (
	"github.com/Alekseyt9/ypmetrics/internal/server/logger"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
)

type Handler struct {
	store    storage.Storage
	settings HandlerSettings
	log      logger.Logger
}

type HandlerSettings struct {
	StoreToFileSync bool // сохранять сразу после изменения значений
	FilePath        string
	DatabaseDSN     string
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
