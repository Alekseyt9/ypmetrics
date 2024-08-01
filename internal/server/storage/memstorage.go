package storage

import (
	"context"
	"sync"

	"github.com/Alekseyt9/ypmetrics/internal/common"
	"github.com/Alekseyt9/ypmetrics/internal/server/filedump"
)

type MemStorage struct {
	counterData map[string]int64
	gaugeData   map[string]float64
	counterLock sync.RWMutex
	gaugeLock   sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counterData: make(map[string]int64),
		gaugeData:   make(map[string]float64),
	}
}

func (store *MemStorage) GetCounter(_ context.Context, name string) (int64, error) {
	store.counterLock.RLock()
	defer store.counterLock.RUnlock()
	v, ok := store.counterData[name]
	if !ok {
		return 0, ErrNotFound
	}
	return v, nil
}

func (store *MemStorage) SetCounter(_ context.Context, name string, value int64) error {
	store.counterLock.Lock()
	defer store.counterLock.Unlock()

	if v, ok := store.counterData[name]; ok {
		store.counterData[name] = v + value
		return nil
	}
	store.counterData[name] = value
	return nil
}

func (store *MemStorage) GetCounters(_ context.Context) ([]common.CounterItem, error) {
	store.counterLock.RLock()
	defer store.counterLock.RUnlock()

	result := make([]common.CounterItem, 0, len(store.counterData))
	for name, value := range store.counterData {
		result = append(result, common.CounterItem{Name: name, Value: value})
	}
	return result, nil
}

func (store *MemStorage) SetCounters(_ context.Context, items []common.CounterItem) error {
	store.counterLock.Lock()
	defer store.counterLock.Unlock()

	for _, item := range items {
		store.counterData[item.Name] = item.Value
	}

	return nil
}

func (store *MemStorage) GetGauge(_ context.Context, name string) (float64, error) {
	store.gaugeLock.RLock()
	defer store.gaugeLock.RUnlock()

	v, ok := store.gaugeData[name]
	if !ok {
		return 0, ErrNotFound
	}
	return v, nil
}

func (store *MemStorage) SetGauge(_ context.Context, name string, value float64) error {
	store.gaugeLock.Lock()
	defer store.gaugeLock.Unlock()

	store.gaugeData[name] = value
	return nil
}

func (store *MemStorage) GetGauges(_ context.Context) ([]common.GaugeItem, error) {
	store.gaugeLock.RLock()
	defer store.gaugeLock.RUnlock()

	result := make([]common.GaugeItem, 0, len(store.gaugeData))
	for name, value := range store.gaugeData {
		result = append(result, common.GaugeItem{Name: name, Value: value})
	}
	return result, nil
}

func (store *MemStorage) SetGauges(_ context.Context, items []common.GaugeItem) error {
	store.gaugeLock.Lock()
	defer store.gaugeLock.Unlock()

	for _, item := range items {
		store.gaugeData[item.Name] = item.Value
	}

	return nil
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
	dc := filedump.NewController()
	return dc.Save(dump, filePath)
}

func (store *MemStorage) LoadFromFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	dump := &filedump.FileDump{}
	dc := filedump.NewController()
	err := dc.Load(dump, filePath)
	if err != nil {
		return err
	}
	store.LoadFromDump(dump)
	return nil
}

func (store *MemStorage) Ping(_ context.Context) error {
	return nil
}
