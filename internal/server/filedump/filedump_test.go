package filedump

import (
	"errors"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndLoad(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "filedump_test_*.json")
	require.NoError(t, err, "Failed to create temporary file")
	defer os.Remove(tmpfile.Name())

	originalDump := &FileDump{
		CounterData: map[string]int64{
			"counter1": 100,
			"counter2": 200,
		},
		GaugeData: map[string]float64{
			"gauge1": 1.1,
			"gauge2": 2.2,
		},
	}

	controller := NewController()
	err = controller.Save(originalDump, tmpfile.Name())
	require.NoError(t, err, "Failed to save data to file")
	loadedDump := &FileDump{}

	err = controller.Load(loadedDump, tmpfile.Name())
	require.NoError(t, err, "Failed to load data from file")
	assert.Equal(t, originalDump, loadedDump, "Loaded data does not match original data")
}

func TestRetriableError(t *testing.T) {
	err := &os.PathError{
		Err: syscall.EAGAIN,
	}

	assert.True(t, isRetriableError(err), "Expected error to be retriable")
	err = &os.PathError{
		Err: syscall.ENOENT,
	}
	assert.False(t, isRetriableError(err), "Expected error to be non-retriable")
}

func TestNonRetriableError(t *testing.T) {
	err := errors.New("non-retriable error")
	assert.False(t, isRetriableError(err), "Expected error to be non-retriable")
}
