package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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
