package dump

import (
	"encoding/json"
	"os"
)

type Dump struct {
	CounterData map[string]int64   `json:"counter"`
	GaugeData   map[string]float64 `json:"gauge"`
}

func (dump Dump) Save(fname string) error {
	data, err := json.MarshalIndent(dump, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(fname, data, 0666) //nolint:gosec //чтобы прошли тесты
}

func (dump *Dump) Load(fname string) error {
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, dump); err != nil {
		return err
	}
	return nil
}
