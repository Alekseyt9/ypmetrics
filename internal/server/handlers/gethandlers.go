package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func HandleGetGauge(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		v, ok := store.GetGauge(name)
		if ok {
			io.WriteString(w, strconv.FormatFloat(v, 'f', -1, 64))
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "metric not found", http.StatusNotFound)
		}

	}
}

func HandleGetCounter(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		v, ok := store.GetCounter(name)
		if ok {
			io.WriteString(w, strconv.FormatInt(v, 10))
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "metric not found", http.StatusNotFound)
		}
	}
}

func HandleGetAll(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Metrics list</title>
			</head>
			<body>
				<ul>
		`))

		for _, item := range store.GetGaugeAll() {
			li := fmt.Sprintf("<li>%s: %s</li>", item.Name, strconv.FormatFloat(item.Value, 'f', -1, 64))
			w.Write([]byte(li))
		}

		for _, item := range store.GetCounterAll() {
			li := fmt.Sprintf("<li>%s: %d</li>", item.Name, item.Value)
			w.Write([]byte(li))
		}

		w.Write([]byte(`
				</ul>
			</body>
			</html>
		`))
	}
}
