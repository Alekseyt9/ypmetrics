package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) HandleGauge(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	gaugeValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		http.Error(w, "incorrect metric value", http.StatusBadRequest)
	}

	err = h.store.SetGauge(r.Context(), name, gaugeValue)
	if err != nil {
		http.Error(w, "error SetGauge", http.StatusBadRequest)
	}

	h.StoreToFile()
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleCounter(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	counterValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		http.Error(w, "incorrect metric value", http.StatusBadRequest)
	}

	err = h.store.SetCounter(r.Context(), name, counterValue)
	if err != nil {
		http.Error(w, "error SetCounter", http.StatusBadRequest)
	}

	h.StoreToFile()
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleIncorrectType(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "incorrect metric type", http.StatusBadRequest)
}

func (h *Handler) HandleNotValue(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "value is missing", http.StatusNotFound)
}
