package storage_test

import (
	"context"
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/common/items"
	"github.com/Alekseyt9/ypmetrics/internal/server/config"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/stretchr/testify/suite"
)

type PGRetryRepositorySuite struct {
	suite.Suite
	store            storage.Storage
	connectionString string
}

func TestPGRetryRepositorySuite(t *testing.T) {
	cfg := &config.Config{}
	config.SetEnv(cfg)

	if cfg.DataBaseDSN == "" {
		t.Skip("Skipping repository tests. Set DATABASE_DSN to run them.")
	}
	suite.Run(t, &PGRepositorySuite{
		connectionString: cfg.DataBaseDSN,
	})
}

func (s *PGRetryRepositorySuite) SetupSuite() {
	store, err := storage.NewDBStorage(s.connectionString)
	if err != nil {
		s.T().Fatalf("Could not init postgres repository: %s", err.Error())
	}
	s.store = storage.NewDBRetryStorage(store)
}

func (s *PGRetryRepositorySuite) TestCounter() {
	ctx := context.Background()
	name := "test_counter"
	value := int64(10)

	err := s.store.SetCounter(ctx, name, value)
	s.Require().NoError(err)

	retrievedValue, err := s.store.GetCounter(ctx, name)
	s.Require().NoError(err)
	s.Equal(value, retrievedValue)

	items := []items.CounterItem{
		{Name: "test_counter1", Value: 20},
		{Name: "test_counter2", Value: 30},
	}
	err = s.store.SetCounters(ctx, items)
	s.Require().NoError(err)

	counters, err := s.store.GetCounters(ctx)
	s.Require().NoError(err)
	s.Equal(len(items)+1, len(counters))
}

func (s *PGRetryRepositorySuite) TestGauge() {
	ctx := context.Background()
	name := "test_gauge"
	value := float64(1.23)

	err := s.store.SetGauge(ctx, name, value)
	s.Require().NoError(err)

	retrievedValue, err := s.store.GetGauge(ctx, name)
	s.Require().NoError(err)
	s.Equal(value, retrievedValue)

	items := []items.GaugeItem{
		{Name: "test_gauge1", Value: 2.34},
		{Name: "test_gauge2", Value: 3.45},
	}
	err = s.store.SetGauges(ctx, items)
	s.Require().NoError(err)

	gauges, err := s.store.GetGauges(ctx)
	s.Require().NoError(err)
	s.Equal(len(items)+1, len(gauges))
}

func (s *PGRetryRepositorySuite) TestPing() {
	ctx := context.Background()
	err := s.store.Ping(ctx)
	s.Require().NoError(err)
}
