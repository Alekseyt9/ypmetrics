package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testReguest(t *testing.T, ts *httptest.Server, method, path string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	require.NoError(t, err)

	return resp
}

func TestRouter(t *testing.T) {
	store := storage.NewMemStorage()
	ts := httptest.NewServer(Router(store))
	defer ts.Close()

	tests := []struct {
		url    string
		status int
	}{
		{
			url:    "/update",
			status: http.StatusBadRequest,
		},
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
		resp := testReguest(t, ts, "POST", v.url)
		assert.Equal(t, v.status, resp.StatusCode, v.url)
	}

}
