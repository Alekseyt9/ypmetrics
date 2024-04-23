package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/mailru/easyjson"
)

func (h *Handler) HandleUpdateJSON(w http.ResponseWriter, r *http.Request) {
	сontentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(сontentType, "application/json") {
		http.Error(w, "incorrect Content-Type", http.StatusUnsupportedMediaType)
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
	}

	if strings.Contains(r.Header.Get("Content-Encoding"), "br") {
		body, err = common.BrotliDecompress(body)
		if err != nil {
			http.Error(w, "error decompress brodli", http.StatusBadRequest)
		}
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
		h.store.SetGauge(data.ID, *data.Value)
		v, b := h.store.GetGauge(data.ID)
		if b {
			restData.Value = &v
		} else {
			var z float64 = 0
			restData.Value = &z
		}
	case "counter":
		h.store.SetCounter(data.ID, *data.Delta)
		v, b := h.store.GetCounter(data.ID)
		if b {
			restData.Delta = &v
		} else {
			var z int64 = 0
			restData.Delta = &z
		}
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

func (h *Handler) HandleValueJSON(w http.ResponseWriter, r *http.Request) {
	сontentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(сontentType, "application/json") {
		http.Error(w, "incorrect Content-Type "+сontentType, http.StatusUnsupportedMediaType)
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
	}

	if strings.Contains(r.Header.Get("Content-Encoding"), "br") {
		body, err = common.BrotliDecompress(body)
		if err != nil {
			http.Error(w, "error decompress brodli", http.StatusBadRequest)
		}
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
		v, b := h.store.GetGauge(data.ID)
		if b {
			restData.Value = &v
		} else {
			var z float64 = 0
			restData.Value = &z
		}
	case "counter":
		v, b := h.store.GetCounter(data.ID)
		if b {
			restData.Delta = &v
		} else {
			var z int64 = 0
			restData.Delta = &z
		}
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
