package handlers

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Alekseyt9/ypmetrics/internal/common/hash"
	"github.com/Alekseyt9/ypmetrics/internal/common/items"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/mailru/easyjson"
	"golang.org/x/net/context"
)

// HandleUpdateJSON handles the update of metrics via JSON data.
func (h *MetricsHandler) HandleUpdateJSON(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "incorrect Content-Type", http.StatusUnsupportedMediaType)
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
	}

	ha := r.Header.Get("HashSHA256")
	if ha != "" {
		if h.settings.HashKey != "" {
			if ha != hash.HashSHA256(body, []byte(h.settings.HashKey)) {
				http.Error(w, "hash check error", http.StatusBadRequest)
			}
		} else {
			h.log.Error("hash key not specified")
		}
	}

	var data items.Metrics
	err = easyjson.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "error unmarshaling JSON", http.StatusBadRequest)
	}

	resData := h.setMetrics(w, r, &data)
	h.StoreToFile()

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

func (h *MetricsHandler) setMetrics(w http.ResponseWriter, r *http.Request, data *items.Metrics) *items.Metrics {
	var resData = &items.Metrics{
		MType: data.MType,
		ID:    data.ID,
	}

	switch data.MType {
	case "gauge":
		err := h.store.SetGauge(r.Context(), data.ID, *data.Value)
		if err != nil {
			http.Error(w, "error SetGauge", http.StatusBadRequest)
		}

		var v float64
		v, err = h.store.GetGauge(r.Context(), data.ID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(w, "metric not found", http.StatusNotFound)
			}
			http.Error(w, "error GetGauge", http.StatusBadRequest)
		}
		resData.Value = &v

	case "counter":
		err := h.store.SetCounter(r.Context(), data.ID, *data.Delta)
		if err != nil {
			http.Error(w, "error SetCounter", http.StatusBadRequest)
		}

		var v int64
		v, err = h.store.GetCounter(r.Context(), data.ID)
		if err != nil {
			http.Error(w, "error GetCounter", http.StatusBadRequest)
		}
		resData.Delta = &v
	}

	return resData
}

// HandleValueJSON handles the retrieval of a single metric value via JSON request.
func (h *MetricsHandler) HandleValueJSON(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "incorrect Content-Type", http.StatusUnsupportedMediaType)
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
	}

	ha := r.Header.Get("HashSHA256")
	if ha != "" && h.settings.HashKey != "" {
		if ha != hash.HashSHA256(body, []byte(h.settings.HashKey)) {
			http.Error(w, "hash check error", http.StatusBadRequest)
		}
	}

	var data items.Metrics
	err = easyjson.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "error unmarshaling JSON", http.StatusBadRequest)
	}

	resData := h.getMetrics(r.Context(), data, w)

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

func (h *MetricsHandler) getMetrics(ctx context.Context, data items.Metrics, w http.ResponseWriter) *items.Metrics {
	var resData = &items.Metrics{
		MType: data.MType,
		ID:    data.ID,
	}

	switch data.MType {
	case "gauge":
		var v float64
		v, err := h.store.GetGauge(ctx, data.ID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				v = 0
			} else {
				http.Error(w, "error GetGauge", http.StatusBadRequest)
			}
		}
		resData.Value = &v

	case "counter":
		var v int64
		v, err := h.store.GetCounter(ctx, data.ID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				v = 0
			} else {
				http.Error(w, "error GetCounter", http.StatusBadRequest)
			}
		}
		resData.Delta = &v
	}
	return resData
}

// HandleUpdateJSONBatch handles the batch update of metrics via JSON data.
func (h *MetricsHandler) HandleUpdateBatchJSON(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "incorrect Content-Type", http.StatusUnsupportedMediaType)
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
	}

	ha := r.Header.Get("HashSHA256")
	if ha != "" && h.settings.HashKey != "" {
		if ha != hash.HashSHA256(body, []byte(h.settings.HashKey)) {
			http.Error(w, "hash check error", http.StatusBadRequest)
		}
	}

	var data items.MetricsSlice
	err = easyjson.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "error unmarshaling JSON", http.StatusBadRequest)
	}

	metricsItems := data.ToMetricItems()

	err = h.store.SetCounters(r.Context(), metricsItems.Counters)
	if err != nil {
		http.Error(w, "error SetCounters", http.StatusBadRequest)
	}

	err = h.store.SetGauges(r.Context(), metricsItems.Gauges)
	if err != nil {
		http.Error(w, "error SetGauges", http.StatusBadRequest)
	}

	resData := items.MetricItems{
		Counters: make([]items.CounterItem, 0),
		Gauges:   make([]items.GaugeItem, 0),
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

	out, err := easyjson.Marshal(resData.ToMetricsSlice())
	if err != nil {
		http.Error(w, "error marshaling JSON", http.StatusBadRequest)
	}
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "error write body", http.StatusBadRequest)
	}
}
