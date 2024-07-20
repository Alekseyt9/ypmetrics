package common

// CounterItem represents a counter metric with a name and value.
type CounterItem struct {
	Name  string
	Value int64
}

// GaugeItem represents a gauge metric with a name and value.
type GaugeItem struct {
	Name  string
	Value float64
}

// MetricItems holds slices of counter and gauge items.
type MetricItems struct {
	Counters []CounterItem
	Gauges   []GaugeItem
}

// ToMetricItems converts a MetricsSlice to a MetricItems structure.
func (s MetricsSlice) ToMetricItems() MetricItems {
	res := MetricItems{
		Counters: make([]CounterItem, 0),
		Gauges:   make([]GaugeItem, 0),
	}

	for _, item := range s {
		switch item.MType {
		case "gauge":
			res.Gauges = append(res.Gauges, GaugeItem{Name: item.ID, Value: *item.Value})
		case "counter":
			res.Counters = append(res.Counters, CounterItem{Name: item.ID, Value: *item.Delta})
		}
	}

	return res
}

// ToMetricsSlice converts a MetricItems structure back to a MetricsSlice.
func (s MetricItems) ToMetricsSlice() MetricsSlice {
	res := make([]Metrics, 0)

	for _, item := range s.Counters {
		res = append(res,
			Metrics{ID: item.Name, MType: "counter", Delta: &item.Value}) //nolint:gosec,exportloopref //version 1.22.2
	}

	for _, item := range s.Gauges {
		res = append(res,
			Metrics{ID: item.Name, MType: "gauge", Value: &item.Value}) //nolint:gosec,exportloopref //version 1.22.2
	}

	return res
}
