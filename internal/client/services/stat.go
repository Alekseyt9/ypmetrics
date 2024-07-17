package services

import (
	"sync"

	"github.com/Alekseyt9/ypmetrics/internal/common"
)

type Stat struct {
	Data *common.MetricItems
	Lock sync.RWMutex
}

func (s *Stat) AddOrUpdateGauge(name string, value float64) {
	old, ok := s.FindGauge(name)
	if ok {
		old.Value = value
		return
	}
	s.Data.Gauges = append(s.Data.Gauges, common.GaugeItem{Name: name, Value: value})
}

func (s *Stat) AddOrUpdateCounter(name string, value int64) {
	old, ok := s.FindCounter(name)
	if ok {
		old.Value = value
		return
	}
	s.Data.Counters = append(s.Data.Counters, common.CounterItem{Name: name, Value: value})
}

func (s *Stat) FindGauge(name string) (*common.GaugeItem, bool) {
	for _, item := range s.Data.Gauges {
		if item.Name == name {
			return &item, true
		}
	}
	return nil, false
}

func (s *Stat) FindCounter(name string) (*common.CounterItem, bool) {
	for _, item := range s.Data.Counters {
		if item.Name == name {
			return &item, true
		}
	}
	return nil, false
}
