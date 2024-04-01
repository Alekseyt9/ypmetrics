package main

import (
	"fmt"
	"net/http"
	"strconv"

	utils "github.com/Alekseyt9/ypmetrics/internal/server/utils"
)

func handleGauge(w http.ResponseWriter, r *http.Request) {
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

func handleCounter(w http.ResponseWriter, r *http.Request) {
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

func checkPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST method only", http.StatusMethodNotAllowed)
	}
}

func handleIncorrectType(w http.ResponseWriter, r *http.Request) {
	checkPost(w, r)
	http.Error(w, "incorrect metric type", http.StatusBadRequest)
}

func main() {
	http.HandleFunc("/update/gauge/", handleGauge)
	http.HandleFunc("/update/counter/", handleCounter)
	http.HandleFunc("/update/", handleIncorrectType)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
