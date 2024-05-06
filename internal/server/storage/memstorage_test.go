package storage_test

import (
	"context"
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestGaugeStorage(t *testing.T) {
	tests := []struct {
		name   string
		set    float64
		want   float64
		metric string
	}{
		{
			name:   "test1",
			set:    1.0,
			want:   1.0,
			metric: "m1",
		},
		{
			name:   "test2",
			set:    -1.0,
			want:   -1.0,
			metric: "m1",
		},
		{
			name:   "test3",
			set:    0,
			want:   0,
			metric: "m2",
		},
		{
			name:   "test4",
			set:    -0.0000000000000001,
			want:   -0.0000000000000001,
			metric: "m2",
		},
	}

	store := storage.NewMemStorage()
	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.SetGauge(ctx, test.metric, test.set)
			v, err := store.GetGauge(ctx, test.metric)
			assert.NoError(t, err)
			assert.InDelta(t, test.want, v, 0.01)
		})
	}
}

func TestCounterStorage(t *testing.T) {
	tests := []struct {
		name   string
		set    int64
		want   int64
		metric string
	}{
		{
			name:   "test1",
			set:    1,
			want:   1,
			metric: "m1",
		},
		{
			name:   "test2",
			set:    2,
			want:   3,
			metric: "m1",
		},
		{
			name:   "test3",
			set:    0,
			want:   0,
			metric: "m2",
		},
		{
			name:   "test4",
			set:    0,
			want:   0,
			metric: "m2",
		},
	}

	store := storage.NewMemStorage()
	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.SetCounter(ctx, test.metric, test.set)
			v, err := store.GetCounter(ctx, test.metric)
			assert.NoError(t, err)
			assert.InDelta(t, test.want, v, 0.01)
		})
	}
}
