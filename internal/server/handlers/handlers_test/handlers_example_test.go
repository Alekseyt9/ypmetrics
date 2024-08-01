package handlers_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/Alekseyt9/ypmetrics/internal/server/handlers"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
)

func ExampleMetricsHandler_HandleUpdateJSON() {
	store := storage.NewMemStorage()

	router := chi.NewRouter()
	h := handlers.NewMetricsHandler(store, handlers.HandlerSettings{})
	router.Post("/update/", h.HandleUpdateJSON)

	ts := httptest.NewServer(router)
	defer ts.Close()

	vg := 1.1
	dataG := common.Metrics{
		ID:    "g",
		MType: "gauge",
		Value: &vg,
	}
	jsonDataG, err := easyjson.Marshal(dataG)
	if err != nil {
		fmt.Println("error marshaling gauge data:", err)
		return
	}

	reqp, err := http.NewRequest(http.MethodPost, ts.URL+"/update/", bytes.NewReader(jsonDataG))
	if err != nil {
		fmt.Println("error creating request:", err)
		return
	}
	reqp.Header.Set("Content-Type", "application/json")

	resp, err := ts.Client().Do(reqp)
	if err != nil {
		fmt.Println("error performing request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)

	var vc int64 = 1
	dataC := common.Metrics{
		ID:    "c",
		MType: "counter",
		Delta: &vc,
	}
	jsonDataC, err := easyjson.Marshal(dataC)
	if err != nil {
		fmt.Println("error marshaling counter data:", err)
		return
	}

	reqp, err = http.NewRequest(http.MethodPost, ts.URL+"/update/", bytes.NewReader(jsonDataC))
	if err != nil {
		fmt.Println("error creating request:", err)
		return
	}
	reqp.Header.Set("Content-Type", "application/json")

	resp, err = ts.Client().Do(reqp)
	if err != nil {
		fmt.Println("error performing request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)

	// Output:
	// Status Code: 200
	// Status Code: 200
}

func ExampleMetricsHandler_HandleValueJSON() {
	store := storage.NewMemStorage()

	router := chi.NewRouter()
	h := handlers.NewMetricsHandler(store, handlers.HandlerSettings{})
	router.Post("/value/", h.HandleValueJSON)

	ts := httptest.NewServer(router)
	defer ts.Close()

	vg := 1.1
	dataG := common.Metrics{
		ID:    "g",
		MType: "gauge",
		Value: &vg,
	}
	jsonDataG, err := easyjson.Marshal(dataG)
	if err != nil {
		fmt.Println("error marshaling gauge data:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/value/", bytes.NewReader(jsonDataG))
	if err != nil {
		fmt.Println("error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ts.Client().Do(req)
	if err != nil {
		fmt.Println("error performing request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)

	var vc int64 = 1
	dataC := common.Metrics{
		ID:    "c",
		MType: "counter",
		Delta: &vc,
	}
	jsonDataC, err := easyjson.Marshal(dataC)
	if err != nil {
		fmt.Println("error marshaling counter data:", err)
		return
	}

	req, err = http.NewRequest(http.MethodPost, ts.URL+"/value/", bytes.NewReader(jsonDataC))
	if err != nil {
		fmt.Println("error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = ts.Client().Do(req)
	if err != nil {
		fmt.Println("error performing request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)

	// Output:
	// Status Code: 200
	// Status Code: 200
}

func ExampleMetricsHandler_HandleUpdateBatchJSON() {
	store := storage.NewMemStorage()

	router := chi.NewRouter()
	h := handlers.NewMetricsHandler(store, handlers.HandlerSettings{})
	router.Post("/updates/", h.HandleUpdateBatchJSON)

	ts := httptest.NewServer(router)
	defer ts.Close()

	vg := 1.1
	vc := int64(1)
	data :=
		common.MetricsSlice{
			{
				ID:    "g",
				MType: "gauge",
				Value: &vg,
			},
			{
				ID:    "c",
				MType: "counter",
				Delta: &vc,
			},
		}

	jsonData, err := easyjson.Marshal(data)
	if err != nil {
		fmt.Println("error marshaling batch data:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/updates/", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println("error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ts.Client().Do(req)
	if err != nil {
		fmt.Println("error performing request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)

	// Output:
	// Status Code: 200
}

func ExampleMetricsHandler_HandleGetAll() {
	store := storage.NewMemStorage()

	router := chi.NewRouter()
	h := handlers.NewMetricsHandler(store, handlers.HandlerSettings{})
	router.Get("/all", h.HandleGetAll)

	ts := httptest.NewServer(router)
	defer ts.Close()

	ctx := context.Background()
	_ = store.SetGauge(ctx, "g1", 1.1)
	_ = store.SetGauge(ctx, "g2", 2.2)
	_ = store.SetCounter(ctx, "c1", 1)
	_ = store.SetCounter(ctx, "c2", 2)

	resp, err := http.Get(ts.URL + "/all")
	if err != nil {
		fmt.Println("error performing request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response:", err)
		return
	}

	responseString := strings.ReplaceAll(string(body), "\t", "")
	fmt.Println("Response Body:")
	fmt.Println(responseString)

	// Output:
	// Status Code: 200
	// Response Body:
	//
	//<!DOCTYPE html>
	//<html>
	//<head>
	//<title>Metrics list</title>
	//</head>
	//<body>
	//<ul>
	//<li>g1: 1.1</li><li>g2: 2.2</li><li>c1: 1</li><li>c2: 2</li>
	//</ul>
	//</body>
	//</html>
}
