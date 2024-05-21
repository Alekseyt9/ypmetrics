package handlers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/server/middleware/logger"
	"github.com/Alekseyt9/ypmetrics/internal/server/run"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	ts *httptest.Server
}

func (suite *TestSuite) SetupSuite() {
	store := storage.NewMemStorage()
	logger := logger.NewSlogLogger()
	cfg := &run.Config{}
	suite.ts = httptest.NewServer(run.Router(store, logger, cfg))
}

func (suite *TestSuite) TearDownSuite() {
	suite.ts.Close()
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
