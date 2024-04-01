package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Alekseyt9/ypmetrics/internal/server/utils"
)

func HandleGauge(w http.ResponseWriter, r *http.Request) {
	checkPost(w, r)

	metricInfo, err := utils.ParseURL(r.URL.Path, "/update/gauge/")
	if err != nil {
		pErr := err.(*utils.URLParseError)
		http.Error(w, pErr.Message, pErr.Status)
	}

	gaugeValue, err := strconv.ParseFloat(metricInfo.Value, 64)
	if err != nil {
		http.Error(w, "incorrect metric value", http.StatusBadRequest)
	}

	fmt.Printf("gauge added: %s-%f", metricInfo.Name, gaugeValue)
	w.WriteHeader(http.StatusOK)
}

func HandleCounter(w http.ResponseWriter, r *http.Request) {
	checkPost(w, r)

	metricInfo, err := utils.ParseURL(r.URL.Path, "/update/counter/")
	if err != nil {
		pErr := err.(*utils.URLParseError)
		http.Error(w, pErr.Message, pErr.Status)
	}

	counterValue, err := strconv.ParseInt(metricInfo.Value, 10, 64)
	if err != nil {
		http.Error(w, "incorrect metric value", http.StatusBadRequest)
	}

	fmt.Printf("counter added: %s-%d", metricInfo.Name, counterValue)
	w.WriteHeader(http.StatusOK)
}

func HandleIncorrectType(w http.ResponseWriter, r *http.Request) {
	checkPost(w, r)
	http.Error(w, "incorrect metric type", http.StatusBadRequest)
}

func checkPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST method only", http.StatusMethodNotAllowed)
	}
}
