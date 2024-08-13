// Package storage defines the interface for storage backends and common storage errors.
package storage

import (
	"context"
	"errors"

	"github.com/Alekseyt9/ypmetrics/internal/common/items"
)

// ErrNotFound is returned when a requested item is not found in the storage.
var ErrNotFound = errors.New("not found")

// Storage defines the interface for metric storage backends.
type Storage interface {
	// GetCounter retrieves the value of the specified counter metric.
	// Parameters:
	//   - ctx: the context for the operation
	//   - name: the name of the counter metric
	// Returns the value of the counter metric and an error if the operation fails.
	GetCounter(ctx context.Context, name string) (int64, error)

	// SetCounter sets the value of the specified counter metric.
	// Parameters:
	//   - ctx: the context for the operation
	//   - name: the name of the counter metric
	//   - value: the value to set for the counter metric
	// Returns an error if the operation fails.
	SetCounter(ctx context.Context, name string, value int64) error

	// GetCounters retrieves all counter metrics.
	// Parameters:
	//   - ctx: the context for the operation
	// Returns a slice of CounterItem and an error if the operation fails.
	GetCounters(ctx context.Context) ([]items.CounterItem, error)

	// SetCounters sets multiple counter metrics.
	// Parameters:
	//   - ctx: the context for the operation
	//   - items: a slice of CounterItem to set
	// Returns an error if the operation fails.
	SetCounters(ctx context.Context, items []items.CounterItem) error

	// GetGauge retrieves the value of the specified gauge metric.
	// Parameters:
	//   - ctx: the context for the operation
	//   - name: the name of the gauge metric
	// Returns the value of the gauge metric and an error if the operation fails.
	GetGauge(ctx context.Context, name string) (float64, error)

	// SetGauge sets the value of the specified gauge metric.
	// Parameters:
	//   - ctx: the context for the operation
	//   - name: the name of the gauge metric
	//   - value: the value to set for the gauge metric
	// Returns an error if the operation fails.
	SetGauge(ctx context.Context, name string, value float64) error

	// GetGauges retrieves all gauge metrics.
	// Parameters:
	//   - ctx: the context for the operation
	// Returns a slice of GaugeItem and an error if the operation fails.
	GetGauges(ctx context.Context) ([]items.GaugeItem, error)

	// SetGauges sets multiple gauge metrics.
	// Parameters:
	//   - ctx: the context for the operation
	//   - items: a slice of GaugeItem to set
	// Returns an error if the operation fails.
	SetGauges(ctx context.Context, items []items.GaugeItem) error

	// Ping checks the connectivity or health of the storage backend.
	// Parameters:
	//   - ctx: the context for managing request lifetime
	// Returns an error if the ping fails.
	Ping(ctx context.Context) error
}
