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

func testReguestPost(t *testing.T, ts *httptest.Server, path string) int {
	req, err := http.NewRequest("POST", ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	require.NoError(t, err)

	return resp.StatusCode
}

func TestRouterPost(t *testing.T) {
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
		statusCode := testReguestPost(t, ts, v.url)
		assert.Equal(t, v.status, statusCode, v.url)
	}
}

func testReguestGet(t *testing.T, ts *httptest.Server, test *getTestStruct) {
	reqp, err := http.NewRequest("POST", ts.URL+test.posturl, nil)
	require.NoError(t, err)
	respp, err := ts.Client().Do(reqp)
	require.NoError(t, err)
	defer respp.Body.Close()
	_, err = io.Copy(io.Discard, respp.Body)
	require.NoError(t, err)

	reqg, err := http.NewRequest("GET", ts.URL+test.geturl, nil)
	require.NoError(t, err)
	respg, err := ts.Client().Do(reqg)
	require.NoError(t, err)
	defer respg.Body.Close()
	bodyBytes, err := io.ReadAll(respg.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, bodyString, test.want)
}

type getTestStruct struct {
	posturl string
	geturl  string
	want    string
}

func TestRouterGet(t *testing.T) {
	store := storage.NewMemStorage()
	ts := httptest.NewServer(Router(store))
	defer ts.Close()

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
		testReguestGet(t, ts, &v)
	}

}
