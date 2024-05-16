package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) HandleGetGauge(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	v, err := h.store.GetGauge(r.Context(), name)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "metric not found", http.StatusNotFound)
		}
		http.Error(w, "error GetGauge", http.StatusBadRequest)
	}

	_, err = io.WriteString(w, strconv.FormatFloat(v, 'f', -1, 64))
	if err != nil {
		http.Error(w, "io.WriteString error", http.StatusBadRequest)
	}
}

func (h *Handler) HandleGetCounter(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	v, err := h.store.GetCounter(r.Context(), name)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "metric not found", http.StatusNotFound)
		}
		http.Error(w, "error GetGauge", http.StatusBadRequest)
	}

	_, err = io.WriteString(w, strconv.FormatInt(v, 10))
	if err != nil {
		http.Error(w, "io.WriteString error", http.StatusBadRequest)
	}
}

func (h *Handler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Metrics list</title>
			</head>
			<body>
				<ul>
		`))
	if err != nil {
		http.Error(w, "w.WriteHeader error", http.StatusBadRequest)
	}

	colGauge, err := h.store.GetGauges(r.Context())
	if err != nil {
		http.Error(w, "error GetGaugeAll", http.StatusBadRequest)
	}
	for _, item := range colGauge {
		li := fmt.Sprintf("<li>%s: %s</li>", item.Name, strconv.FormatFloat(item.Value, 'f', -1, 64))
		_, err = w.Write([]byte(li))
		if err != nil {
			http.Error(w, "w.Write error", http.StatusBadRequest)
		}
	}

	colCounter, err := h.store.GetCounters(r.Context())
	if err != nil {
		http.Error(w, "error GetCounterAll", http.StatusBadRequest)
	}
	for _, item := range colCounter {
		li := fmt.Sprintf("<li>%s: %d</li>", item.Name, item.Value)
		_, err = w.Write([]byte(li))
		if err != nil {
			http.Error(w, "w.Write error", http.StatusBadRequest)
		}
	}

	_, err = w.Write([]byte(`
				</ul>
			</body>
			</html>
		`))
	if err != nil {
		http.Error(w, "w.Write error", http.StatusBadRequest)
	}
}

func (h *Handler) HandlePing(w http.ResponseWriter, r *http.Request) {
	err := h.store.Ping(r.Context())
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
