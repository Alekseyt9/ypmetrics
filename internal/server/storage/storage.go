package storage

type MemStorage struct {
	counterData map[string]int64
	gaugeData   map[string]float64
}

type Storage interface {
	GetCounter(name string) int64
	SetCounter(name string, value int64)
	GetGauge(name string) float64
	SetGauge(name string, value float64)
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counterData: make(map[string]int64),
		gaugeData:   make(map[string]float64),
	}
}

func (store *MemStorage) GetCounter(name string) int64 {
	return store.counterData[name]
}

func (store *MemStorage) SetCounter(name string, value int64) {
	if v, ok := store.counterData[name]; ok {
		store.counterData[name] = v + value
		return
	}
	store.counterData[name] = value
}

func (store *MemStorage) GetGauge(name string) float64 {
	return store.gaugeData[name]
}

func (store *MemStorage) SetGauge(name string, value float64) {
	store.gaugeData[name] = value
}
