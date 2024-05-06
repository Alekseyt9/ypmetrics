package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/mailru/easyjson"
)

func (h *Handler) HandleUpdateJSON(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "incorrect Content-Type", http.StatusUnsupportedMediaType)
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
	}

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		body, err = common.GZIPDecompress(body)
		if err != nil {
			http.Error(w, "error decompress gzip", http.StatusBadRequest)
		}
	}

	var data common.Metrics
	err = easyjson.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "error unmarshaling JSON", http.StatusBadRequest)
	}

	var restData = common.Metrics{
		MType: data.MType,
		ID:    data.ID,
	}

	switch data.MType {
	case "gauge":
		err := h.store.SetGauge(r.Context(), data.ID, *data.Value)
		if err != nil {
			http.Error(w, "error SetGauge", http.StatusBadRequest)
		}

		v, err := h.store.GetGauge(r.Context(), data.ID)
		if err != nil {
			http.Error(w, "error GetGauge", http.StatusBadRequest)
		}
		restData.Value = &v

	case "counter":
		err := h.store.SetCounter(r.Context(), data.ID, *data.Delta)
		if err != nil {
			http.Error(w, "error SetCounter", http.StatusBadRequest)
		}

		v, err := h.store.GetCounter(r.Context(), data.ID)
		if err != nil {
			http.Error(w, "error GetCounter", http.StatusBadRequest)
		}
		restData.Delta = &v
	}

	h.StoreToFile()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	out, err := easyjson.Marshal(restData)
	if err != nil {
		http.Error(w, "error marshaling JSON", http.StatusBadRequest)
	}
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "error write body", http.StatusBadRequest)
	}
}

func (h *Handler) HandleValueJSON(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "incorrect Content-Type", http.StatusUnsupportedMediaType)
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
	}

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		body, err = common.GZIPDecompress(body)
		if err != nil {
			http.Error(w, "error decompress GZIP", http.StatusBadRequest)
		}
	}

	var data common.Metrics
	err = easyjson.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "error unmarshaling JSON", http.StatusBadRequest)
	}

	var restData = common.Metrics{
		MType: data.MType,
		ID:    data.ID,
	}

	switch data.MType {
	case "gauge":
		v, err := h.store.GetGauge(r.Context(), data.ID)
		if err != nil {
			http.Error(w, "error GetGauge", http.StatusBadRequest)
		}
		restData.Value = &v

	case "counter":
		v, err := h.store.GetCounter(r.Context(), data.ID)
		if err != nil {
			http.Error(w, "error GetCounter", http.StatusBadRequest)
		}
		restData.Delta = &v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	out, err := easyjson.Marshal(restData)
	if err != nil {
		http.Error(w, "error marshaling JSON", http.StatusBadRequest)
	}
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "error write body", http.StatusBadRequest)
	}
}
