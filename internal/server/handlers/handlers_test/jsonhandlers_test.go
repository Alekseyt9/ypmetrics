package handlers_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type jsonTestStruct struct {
	mType        string
	ID           string
	valueGauge   float64
	wantGauge    float64
	valueCounter int64
	wantCounter  int64
}

func testReguestJSON(t *testing.T, ts *httptest.Server, test *jsonTestStruct) {

	data := common.Metrics{
		ID:    test.ID,
		MType: test.mType,
	}
	switch test.mType {
	case "gauge":
		data.Value = &test.valueGauge
	case "counter":
		data.Delta = &test.valueCounter
	}

	jsonData, err := easyjson.Marshal(data)
	require.NoError(t, err)

	reqp, err := http.NewRequest(http.MethodPost, ts.URL+"/update/", bytes.NewReader(jsonData))
	require.NoError(t, err)
	reqp.Header.Set("Content-Type", "application/json")
	respp, err := ts.Client().Do(reqp)
	require.NoError(t, err)
	defer respp.Body.Close()
	_, err = io.Copy(io.Discard, respp.Body)
	require.NoError(t, err)

	data1 := common.Metrics{
		ID:    test.ID,
		MType: test.mType,
	}
	jsonData, err = easyjson.Marshal(data1)
	require.NoError(t, err)

	reqg, err := http.NewRequest(http.MethodPost, ts.URL+"/value/", bytes.NewReader(jsonData))
	require.NoError(t, err)
	reqg.Header.Set("Content-Type", "application/json")
	respg, err := ts.Client().Do(reqg)
	require.NoError(t, err)
	defer respg.Body.Close()
	bodyBytes, err := io.ReadAll(respg.Body)
	require.NoError(t, err)

	var vData common.Metrics
	log.Printf("!!! %v", string(bodyBytes))
	err = easyjson.Unmarshal(bodyBytes, &vData)
	require.NoError(t, err)

	switch test.mType {
	case "gauge":
		assert.Equal(t, test.wantGauge, *vData.Value, "gauge")
	case "counter":
		assert.Equal(t, test.wantCounter, *vData.Delta, "counter")
	}
}

func (suite *TestSuite) TestRouterJSON() {
	tests := []jsonTestStruct{
		{
			mType:      "gauge",
			ID:         "mj",
			valueGauge: 1,
			wantGauge:  1,
		},
		{
			mType:      "gauge",
			ID:         "mj",
			valueGauge: -0.1,
			wantGauge:  -0.1,
		},
		{
			mType:        "counter",
			ID:           "mj",
			valueCounter: 1,
			wantCounter:  1,
		},
		{
			mType:        "counter",
			ID:           "mj",
			valueCounter: 1,
			wantCounter:  2,
		},
	}

	for _, v := range tests {
		testReguestJSON(suite.T(), suite.ts, &v)
	}
}
