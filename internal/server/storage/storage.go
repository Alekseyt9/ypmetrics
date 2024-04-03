package storage

type MemStorage struct {
	counterData map[string]int64
	gaugeData   map[string]float64
}

type NameValueGauge struct {
	Name  string
	Value float64
}

type NameValueCounter struct {
	Name  string
	Value int64
}

type Storage interface {
	GetCounter(name string) (int64, bool)
	SetCounter(name string, value int64)
	GetCounterAll() []NameValueCounter

	GetGauge(name string) (float64, bool)
	SetGauge(name string, value float64)
	GetGaugeAll() []NameValueGauge
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counterData: make(map[string]int64),
		gaugeData:   make(map[string]float64),
	}
}

func (store *MemStorage) GetCounter(name string) (int64, bool) {
	v, ok := store.counterData[name]
	return v, ok
}

func (store *MemStorage) SetCounter(name string, value int64) {
	if v, ok := store.counterData[name]; ok {
		store.counterData[name] = v + value
		return
	}
	store.counterData[name] = value
}

func (store *MemStorage) GetCounterAll() []NameValueCounter {
	result := make([]NameValueCounter, 0, len(store.counterData))
	for name, value := range store.counterData {
		result = append(result, NameValueCounter{Name: name, Value: value})
	}
	return result
}

func (store *MemStorage) GetGauge(name string) (float64, bool) {
	v, ok := store.gaugeData[name]
	return v, ok
}

func (store *MemStorage) SetGauge(name string, value float64) {
	store.gaugeData[name] = value
}

func (store *MemStorage) GetGaugeAll() []NameValueGauge {
	result := make([]NameValueGauge, 0, len(store.gaugeData))
	for name, value := range store.gaugeData {
		result = append(result, NameValueGauge{Name: name, Value: value})
	}
	return result
}
