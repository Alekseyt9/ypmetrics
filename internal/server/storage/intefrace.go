package storage

import (
	"context"
	"errors"

	"github.com/Alekseyt9/ypmetrics/internal/common"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	GetCounter(ctx context.Context, name string) (int64, error)
	SetCounter(ctx context.Context, name string, value int64) error
	GetCounters(ctx context.Context) ([]common.CounterItem, error)
	SetCounters(ctx context.Context, items []common.CounterItem) error

	GetGauge(ctx context.Context, name string) (float64, error)
	SetGauge(ctx context.Context, name string, value float64) error
	GetGauges(ctx context.Context) ([]common.GaugeItem, error)
	SetGauges(ctx context.Context, items []common.GaugeItem) error

	Ping(ctx context.Context) error
}
