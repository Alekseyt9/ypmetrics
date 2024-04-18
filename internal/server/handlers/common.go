package handlers

import "github.com/Alekseyt9/ypmetrics/internal/server/storage"

type Handler struct {
	store storage.Storage
}

func NewHandler(store storage.Storage) *Handler {
	return &Handler{
		store: store,
	}
}
