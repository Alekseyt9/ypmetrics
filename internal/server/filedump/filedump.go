package filedump

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"syscall"

	"github.com/Alekseyt9/ypmetrics/internal/common"
)

type FileDump struct {
	CounterData map[string]int64   `json:"counter"`
	GaugeData   map[string]float64 `json:"gauge"`
}

type Controller struct {
	mutex sync.Mutex
}

func NewController() *Controller {
	return &Controller{}
}

func (dc *Controller) Save(dump *FileDump, fname string) error {
	data, err := json.MarshalIndent(dump, "", "   ")
	if err != nil {
		return err
	}

	rc := common.NewRetryControllerStd(isRetriableError)
	err = rc.Do(func() error {
		dc.mutex.Lock()
		defer dc.mutex.Unlock()
		return os.WriteFile(fname, data, 0666) //nolint:gosec //for passing tests
	})
	if err != nil {
		return err
	}
	return nil
}

func (dc *Controller) Load(dump *FileDump, fname string) error {
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
		var errno syscall.Errno
		if errors.As(pathErr.Err, &errno) {
			switch errno { //nolint:exhaustive //-
			case syscall.EACCES, syscall.EAGAIN, syscall.EBUSY:
				return true
			default:
				return false
			}
		}
	}
	return false
}
