package storage

import (
	"sync"

	"github.com/Alekseyt9/ypmetrics/internal/server/filedump"
)

type MemStorage struct {
	counterData map[string]int64
	counterLock sync.RWMutex

	gaugeData map[string]float64
	gaugeLock sync.RWMutex
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

	LoadFromDump(dump *filedump.FileDump)
	SaveToDump(dump *filedump.FileDump)

	SaveToFile(filePath string) error
	LoadFromFile(filePath string) error
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counterData: make(map[string]int64),
		gaugeData:   make(map[string]float64),
	}
}

func (store *MemStorage) GetCounter(name string) (int64, bool) {
	store.counterLock.RLock()
	defer store.counterLock.RUnlock()
	v, ok := store.counterData[name]
	return v, ok
}

func (store *MemStorage) SetCounter(name string, value int64) {
	store.counterLock.Lock()
	defer store.counterLock.Unlock()

	if v, ok := store.counterData[name]; ok {
		store.counterData[name] = v + value
		return
	}
	store.counterData[name] = value
}

func (store *MemStorage) GetCounterAll() []NameValueCounter {
	store.counterLock.RLock()
	defer store.counterLock.RUnlock()

	result := make([]NameValueCounter, 0, len(store.counterData))
	for name, value := range store.counterData {
		result = append(result, NameValueCounter{Name: name, Value: value})
	}
	return result
}

func (store *MemStorage) GetGauge(name string) (float64, bool) {
	store.gaugeLock.RLock()
	defer store.gaugeLock.RUnlock()

	v, ok := store.gaugeData[name]
	return v, ok
}

func (store *MemStorage) SetGauge(name string, value float64) {
	store.gaugeLock.Lock()
	defer store.gaugeLock.Unlock()

	store.gaugeData[name] = value
}

func (store *MemStorage) GetGaugeAll() []NameValueGauge {
	store.gaugeLock.RLock()
	defer store.gaugeLock.RUnlock()

	result := make([]NameValueGauge, 0, len(store.gaugeData))
	for name, value := range store.gaugeData {
		result = append(result, NameValueGauge{Name: name, Value: value})
	}
	return result
}

func (store *MemStorage) LoadFromDump(dump *filedump.FileDump) {
	store.gaugeLock.Lock()
	defer store.gaugeLock.Unlock()

	for k, v := range dump.GaugeData {
		store.gaugeData[k] = v
	}
	for k, v := range dump.CounterData {
		store.counterData[k] = v
	}
}

func (store *MemStorage) SaveToDump(dump *filedump.FileDump) {
	store.gaugeLock.RLock()
	defer store.gaugeLock.RUnlock()

	for k, v := range store.gaugeData {
		dump.GaugeData[k] = v
	}
	for k, v := range store.counterData {
		dump.CounterData[k] = v
	}
}

func (store *MemStorage) SaveToFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	dump := &filedump.FileDump{
		CounterData: make(map[string]int64),
		GaugeData:   make(map[string]float64),
	}
	store.SaveToDump(dump)
	return dump.Save(filePath)
}

func (store *MemStorage) LoadFromFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	dump := &filedump.FileDump{}
	err := dump.Load(filePath)
	if err != nil {
		return err
	}
	store.LoadFromDump(dump)
	return nil
}
