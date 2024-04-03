package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

func (suite *TestSuite) TestRouterPost() {
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
		statusCode := testReguestPost(suite.T(), suite.ts, v.url)
		assert.Equal(suite.T(), v.status, statusCode, v.url)
	}
}
