package storage_test

import (
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.SetGauge(test.metric, test.set)
			v, ok := store.GetGauge(test.metric)
			assert.True(t, ok)
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.SetCounter(test.metric, test.set)
			v, ok := store.GetCounter(test.metric)
			assert.True(t, ok)
			assert.InDelta(t, test.want, v, 0.01)
		})
	}
}
