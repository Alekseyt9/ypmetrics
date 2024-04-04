package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) HandleGetGauge(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	v, ok := h.store.GetGauge(name)
	if ok {
		_, err := io.WriteString(w, strconv.FormatFloat(v, 'f', -1, 64))
		if err != nil {
			log.Fatalf("io.WriteString error %v", err)
		}
	} else {
		http.Error(w, "metric not found", http.StatusNotFound)
	}
}

func (h *Handler) HandleGetCounter(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	v, ok := h.store.GetCounter(name)
	if ok {
		_, err := io.WriteString(w, strconv.FormatInt(v, 10))
		if err != nil {
			log.Fatalf("io.WriteString error %v", err)
		}
	} else {
		http.Error(w, "metric not found", http.StatusNotFound)
	}
}

func (h *Handler) HandleGetAll(w http.ResponseWriter, _ *http.Request) {
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
		log.Fatalf("w.WriteHeader error %v", err)
	}

	for _, item := range h.store.GetGaugeAll() {
		li := fmt.Sprintf("<li>%s: %s</li>", item.Name, strconv.FormatFloat(item.Value, 'f', -1, 64))
		_, err = w.Write([]byte(li))
		if err != nil {
			log.Fatalf("w.Write error %v", err)
		}
	}

	for _, item := range h.store.GetCounterAll() {
		li := fmt.Sprintf("<li>%s: %d</li>", item.Name, item.Value)
		_, err = w.Write([]byte(li))
		if err != nil {
			log.Fatalf("w.Write error %v", err)
		}
	}

	_, err = w.Write([]byte(`
				</ul>
			</body>
			</html>
		`))
	if err != nil {
		log.Fatalf("w.Write error %v", err)
	}
}
