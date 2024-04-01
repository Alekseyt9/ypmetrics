package handlers

import (
	"net/http"
	"strconv"

	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func HandleGauge(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")

		gaugeValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "incorrect metric value", http.StatusBadRequest)
		}

		store.SetGauge(name, gaugeValue)
		w.WriteHeader(http.StatusOK)
	}
}

func HandleCounter(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")

		counterValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "incorrect metric value", http.StatusBadRequest)
		}

		store.SetCounter(name, counterValue)
		w.WriteHeader(http.StatusOK)
	}
}

func HandleIncorrectType(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "incorrect metric type", http.StatusBadRequest)
}

func HandleNotValue(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "value is missing", http.StatusNotFound)
}
