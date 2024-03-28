package storage

type MemStorage struct {
	counterData map[string]int64
	gaugeData   map[string]float64
}

type Storage interface {
	GetCounter(name string)
	SetCounter(name string, value int64)
	GetGauge(name string)
	SetGauge(name string, value int64)
}
