package handlers

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
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
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(w, "metric not found", http.StatusNotFound)
			}
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

	var resData = common.Metrics{
		MType: data.MType,
		ID:    data.ID,
	}

	switch data.MType {
	case "gauge":
		v, err := h.store.GetGauge(r.Context(), data.ID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				v = 0
			} else {
				http.Error(w, "error GetGauge", http.StatusBadRequest)
			}
		}
		resData.Value = &v

	case "counter":
		v, err := h.store.GetCounter(r.Context(), data.ID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				v = 0
			} else {
				http.Error(w, "error GetCounter", http.StatusBadRequest)
			}
		}
		resData.Delta = &v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	out, err := easyjson.Marshal(resData)
	if err != nil {
		http.Error(w, "error marshaling JSON", http.StatusBadRequest)
	}
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "error write body", http.StatusBadRequest)
	}
}

func (h *Handler) HandleUpdateBatchJSON(w http.ResponseWriter, r *http.Request) {
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

	var data common.MetricsBatch
	err = easyjson.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "error unmarshaling JSON", http.StatusBadRequest)
	}

	err = h.store.SetCounters(r.Context(), data.Counters)
	if err != nil {
		http.Error(w, "error SetCounters", http.StatusBadRequest)
	}

	err = h.store.SetGauges(r.Context(), data.Gauges)
	if err != nil {
		http.Error(w, "error SetGauges", http.StatusBadRequest)
	}

	resData := common.MetricsBatch{
		Counters: make([]common.CounterItem, 1),
		Gauges:   make([]common.GaugeItem, 1),
	}

	resData.Counters, err = h.store.GetCounters(r.Context())
	if err != nil {
		http.Error(w, "error GetCounters", http.StatusBadRequest)
	}

	resData.Gauges, err = h.store.GetGauges(r.Context())
	if err != nil {
		http.Error(w, "error GetGauges", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	out, err := easyjson.Marshal(resData)
	if err != nil {
		http.Error(w, "error marshaling JSON", http.StatusBadRequest)
	}
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "error write body", http.StatusBadRequest)
	}
}
