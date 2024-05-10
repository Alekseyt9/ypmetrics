package filedump

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"syscall"

	"github.com/Alekseyt9/ypmetrics/internal/common"
)

var mutex sync.Mutex

type FileDump struct {
	CounterData map[string]int64   `json:"counter"`
	GaugeData   map[string]float64 `json:"gauge"`
}

func (dump FileDump) Save(fname string) error {
	data, err := json.MarshalIndent(dump, "", "   ")
	if err != nil {
		return err
	}

	rc := common.NewRetryControllerStd(isRetriableError)
	err = rc.Do(func() error {
		mutex.Lock()
		defer mutex.Unlock()
		return os.WriteFile(fname, data, 0666) //nolint:gosec //to pass tests
	})
	if err != nil {
		return err
	}
	return nil
}

func (dump *FileDump) Load(fname string) error {
	rc := common.NewRetryControllerStd(isRetriableError)
	err := rc.Do(func() error {
		data, err := os.ReadFile(fname)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(data, dump); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func isRetriableError(err error) bool {
	if err == nil {
		return false
	}

	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		if errno, ok := pathErr.Err.(syscall.Errno); ok {
			switch errno {
			case syscall.EACCES, syscall.EAGAIN, syscall.EBUSY:
				return true
			}
		}
	}
	return false
}
