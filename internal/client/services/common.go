package services

import "sync"

type GaugeItem struct {
	name  string
	value float64
}

type CounterItem struct {
	name  string
	value int64
}

type Stat struct {
	Counters []CounterItem
	Gauges   []GaugeItem
	MapLock  sync.RWMutex
}

func (s *Stat) AddGauge(name string, value float64) {
	s.Gauges = append(s.Gauges, GaugeItem{name: name, value: value})
}

func (s *Stat) AddCounter(name string, value int64) {
	s.Counters = append(s.Counters, CounterItem{name: name, value: value})
}

func (s *Stat) FindGauge(name string) (*GaugeItem, bool) {
	for _, item := range s.Gauges {
		if item.name == name {
			return &item, true
		}
	}
	return nil, false
}

func (s *Stat) FindCounter(name string) (*CounterItem, bool) {
	for _, item := range s.Counters {
		if item.name == name {
			return &item, true
		}
	}
	return nil, false
}
