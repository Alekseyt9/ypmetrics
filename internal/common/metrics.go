package common

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type GaugeItem struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type CounterItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type MetricsBatch struct {
	Counters []CounterItem `json:"counters"`
	Gauges   []GaugeItem   `json:"gauges"`
}
