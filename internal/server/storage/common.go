package storage

import "context"

type NameValueGauge struct {
	Name  string
	Value float64
}

type NameValueCounter struct {
	Name  string
	Value int64
}

type Storage interface {
	GetCounter(ctx context.Context, name string) (int64, error)
	SetCounter(ctx context.Context, name string, value int64) error
	GetCounterAll(ctx context.Context) ([]NameValueCounter, error)

	GetGauge(ctx context.Context, name string) (float64, error)
	SetGauge(ctx context.Context, name string, value float64) error
	GetGaugeAll(ctx context.Context) ([]NameValueGauge, error)
}
