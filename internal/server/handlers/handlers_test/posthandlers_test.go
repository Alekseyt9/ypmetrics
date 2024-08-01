package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/server/config"
	"github.com/Alekseyt9/ypmetrics/internal/server/log"
	"github.com/Alekseyt9/ypmetrics/internal/server/run"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/stretchr/testify/require"
)

func testReguestPost(t *testing.T, ts *httptest.Server, path string) int {
	req, err := http.NewRequest(http.MethodPost, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	require.NoError(t, err)

	return resp.StatusCode
}

func (suite *TestSuite) TestRouterPost() {
	tests := []struct {
		url    string
		status int
	}{
		/*
			{
				url:    "/update",
				status: http.StatusBadRequest,
			},*/ // not actual after iter 7
		{
			url:    "/update/unknown/",
			status: http.StatusBadRequest,
		},
		{
			url:    "/update/gauge/",
			status: http.StatusNotFound,
		},
		{
			url:    "/update/gauge/m1/1",
			status: http.StatusOK,
		},
		{
			url:    "/update/counter/m1/1",
			status: http.StatusOK,
		},
	}

	for _, v := range tests {
		statusCode := testReguestPost(suite.T(), suite.ts, v.url)
		suite.Equal(v.status, statusCode, v.url)
	}
}

func BenchmarkHandlePost(b *testing.B) {
	store := storage.NewMemStorage()
	logger := log.NewNoOpLogger()
	cfg := &config.Config{}
	ts := httptest.NewServer(run.Router(store, logger, cfg))

	b.Run("gauge_update", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, err := http.NewRequest(http.MethodPost, ts.URL+"/update/gauge/g1/1", nil)
			require.NoError(b, err)
			req.Header.Set("Content-Type", "text/plain")
			resp, err := ts.Client().Do(req)
			require.NoError(b, err)
			defer resp.Body.Close()
		}
	})

	b.Run("counter_update", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, err := http.NewRequest(http.MethodPost, ts.URL+"/update/counter/c1/1", nil)
			require.NoError(b, err)
			req.Header.Set("Content-Type", "text/plain")
			resp, err := ts.Client().Do(req)
			require.NoError(b, err)
			defer resp.Body.Close()
		}
	})
}
