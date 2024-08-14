package services

import (
	"sync"

	"github.com/Alekseyt9/ypmetrics/internal/common/items"
)

// Stat holds metric data and a read-write mutex for synchronized access.
type Stat struct {
	Data *items.MetricItems
	Lock sync.RWMutex
}

// AddOrUpdateGauge adds a new gauge metric or updates an existing one.
// Parameters:
//   - name: the name of the gauge metric
//   - value: the value of the gauge metric
func (s *Stat) AddOrUpdateGauge(name string, value float64) {
	old, ok := s.FindGauge(name)
	if ok {
		old.Value = value
		return
	}
	s.Data.Gauges = append(s.Data.Gauges, items.GaugeItem{Name: name, Value: value})
}

// AddOrUpdateCounter adds a new counter metric or updates an existing one.
// Parameters:
//   - name: the name of the counter metric
//   - value: the value of the counter metric
func (s *Stat) AddOrUpdateCounter(name string, value int64) {
	old, ok := s.FindCounter(name)
	if ok {
		old.Value = value
		return
	}
	s.Data.Counters = append(s.Data.Counters, items.CounterItem{Name: name, Value: value})
}

// FindGauge searches for a gauge metric by name.
// Parameters:
//   - name: the name of the gauge metric
//
// Returns:
//   - a pointer to the found GaugeItem and a boolean indicating if it was found
func (s *Stat) FindGauge(name string) (*items.GaugeItem, bool) {
	for i := 0; i < len(s.Data.Gauges); i++ {
		if s.Data.Gauges[i].Name == name {
			return &s.Data.Gauges[i], true
		}
	}
	return nil, false
}

// FindCounter searches for a counter metric by name.
// Parameters:
//   - name: the name of the counter metric
//
// Returns:
//   - a pointer to the found CounterItem and a boolean indicating if it was found
func (s *Stat) FindCounter(name string) (*items.CounterItem, bool) {
	for i := 0; i < len(s.Data.Counters); i++ {
		if s.Data.Counters[i].Name == name {
			return &s.Data.Counters[i], true
		}
	}
	return nil, false
}
