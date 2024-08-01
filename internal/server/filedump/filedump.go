// Package filedump provides functionality for saving and loading metric data to and from a file.
package filedump

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"syscall"

	"github.com/Alekseyt9/ypmetrics/pkg/retry"
)

// FileDump holds the data for counters and gauges.
type FileDump struct {
	CounterData map[string]int64   `json:"counter"`
	GaugeData   map[string]float64 `json:"gauge"`
}

// Controller manages the file operations for saving and loading FileDump.
type Controller struct {
	mutex sync.Mutex
}

// NewController creates a new Controller instance.
// Returns a pointer to a Controller.
func NewController() *Controller {
	return &Controller{}
}

// Save writes the given FileDump data to a file with the specified name.
// Parameters:
//   - dump: the FileDump data to save
//   - fname: the name of the file to write to
//
// Returns an error if the save operation fails.
func (dc *Controller) Save(dump *FileDump, fname string) error {
	data, err := json.MarshalIndent(dump, "", "   ")
	if err != nil {
		return err
	}

	rc := retry.NewControllerStd(isRetriableError)
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

// Load reads the FileDump data from a file with the specified name.
// Parameters:
//   - dump: the FileDump structure to populate with the loaded data
//   - fname: the name of the file to read from
//
// Returns an error if the load operation fails.
func (dc *Controller) Load(dump *FileDump, fname string) error {
	rc := retry.NewControllerStd(isRetriableError)
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
