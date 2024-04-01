package main

import (
	"net/http"

	handlers "github.com/Alekseyt9/ypmetrics/internal/server/handlers"
)

func main() {
	http.HandleFunc("/update/gauge/", handlers.HandleGauge)
	http.HandleFunc("/update/counter/", handlers.HandleCounter)
	http.HandleFunc("/update/", handlers.HandleIncorrectType)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
