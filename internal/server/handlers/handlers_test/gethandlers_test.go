package handlers_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/server/config"
	"github.com/Alekseyt9/ypmetrics/internal/server/log"
	"github.com/Alekseyt9/ypmetrics/internal/server/run"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testReguestGet(t *testing.T, ts *httptest.Server, test *getTestStruct) {
	reqp, err := http.NewRequest(http.MethodPost, ts.URL+test.posturl, nil)
	require.NoError(t, err)
	respp, err := ts.Client().Do(reqp)
	require.NoError(t, err)
	defer respp.Body.Close()
	_, err = io.Copy(io.Discard, respp.Body)
	require.NoError(t, err)

	reqg, err := http.NewRequest(http.MethodGet, ts.URL+test.geturl, nil)
	require.NoError(t, err)
	respg, err := ts.Client().Do(reqg)
	require.NoError(t, err)
	defer respg.Body.Close()
	bodyBytes, err := io.ReadAll(respg.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, test.want, bodyString)
}

type getTestStruct struct {
	posturl string
	geturl  string
	want    string
}

func (suite *TestSuite) TestRouterGet() {
	tests := []getTestStruct{
		{
			posturl: "/update/gauge/m1/1",
			geturl:  "/value/gauge/m1",
			want:    "1",
		},
		{
			posturl: "/update/gauge/m1/-0.1",
			geturl:  "/value/gauge/m1",
			want:    "-0.1",
		},
		{
			posturl: "/update/counter/m1/1",
			geturl:  "/value/counter/m1",
			want:    "1",
		},
		{
			posturl: "/update/counter/m1/1",
			geturl:  "/value/counter/m1",
			want:    "2",
		},
	}

	for _, v := range tests {
		testReguestGet(suite.T(), suite.ts, &v)
	}
}

func (suite *TestSuite) TestHandleGetAll() {
	tests := []getTestStruct{
		{
			posturl: "/update/gauge/m1/1.1",
			geturl:  "/",
		},
		{
			posturl: "/update/counter/c1/1",
			geturl:  "/",
		},
		{
			posturl: "/update/counter/c2/2",
			geturl:  "/",
		},
	}

	for _, v := range tests {
		reqp, err := http.NewRequest(http.MethodPost, suite.ts.URL+v.posturl, nil)
		require.NoError(suite.T(), err)
		respp, err := suite.ts.Client().Do(reqp)
		require.NoError(suite.T(), err)
		defer respp.Body.Close()
		_, err = io.Copy(io.Discard, respp.Body)
		require.NoError(suite.T(), err)
	}

	reqg, err := http.NewRequest(http.MethodGet, suite.ts.URL+"/", nil)
	require.NoError(suite.T(), err)
	respg, err := suite.ts.Client().Do(reqg)
	require.NoError(suite.T(), err)
	defer respg.Body.Close()
	bodyBytes, err := io.ReadAll(respg.Body)
	require.NoError(suite.T(), err)
	bodyString := string(bodyBytes)

	assert.Contains(suite.T(), bodyString, "m1: 1.1")
	assert.Contains(suite.T(), bodyString, "c1: 1")
	assert.Contains(suite.T(), bodyString, "c2: 2")
}

func (suite *TestSuite) TestHandlePing() {
	req, err := http.NewRequest(http.MethodGet, suite.ts.URL+"/ping", nil)
	require.NoError(suite.T(), err)
	resp, err := suite.ts.Client().Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func BenchmarkHandleGetAll(b *testing.B) {
	store := storage.NewMemStorage()
	logger := log.NewNoOpLogger()
	cfg := &config.Config{}
	ts := httptest.NewServer(run.Router(store, logger, cfg))

	ctx := context.Background()
	for i := 0; i < 1000; i++ {
		require.NoError(b, store.SetGauge(ctx, fmt.Sprintf("gauge%d", i), float64(i)))
		require.NoError(b, store.SetCounter(ctx, fmt.Sprintf("counter%d", i), int64(i)))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reqg, err := http.NewRequest(http.MethodGet, ts.URL+"/", nil)
		require.NoError(b, err)
		respg, err := ts.Client().Do(reqg)
		require.NoError(b, err)
		defer respg.Body.Close()
	}
}

func BenchmarkHandleGet(b *testing.B) {
	store := storage.NewMemStorage()
	logger := log.NewNoOpLogger()
	cfg := &config.Config{}
	ts := httptest.NewServer(run.Router(store, logger, cfg))

	ctx := context.Background()
	require.NoError(b, store.SetGauge(ctx, "g1", float64(1)))
	require.NoError(b, store.SetCounter(ctx, "c1", int64(1)))

	b.Run("gauge_value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reqg, err := http.NewRequest(http.MethodGet, ts.URL+"/value/gauge/g1", nil)
			require.NoError(b, err)
			respg, err := ts.Client().Do(reqg)
			require.NoError(b, err)
			defer respg.Body.Close()
		}
	})

	b.Run("counter_value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reqg, err := http.NewRequest(http.MethodGet, ts.URL+"/value/counter/c1", nil)
			require.NoError(b, err)
			respg, err := ts.Client().Do(reqg)
			require.NoError(b, err)
			defer respg.Body.Close()
		}
	})
}
