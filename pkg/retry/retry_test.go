package retry

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestController_Do_Success(t *testing.T) {
	needRetry := func(err error) bool {
		return false
	}
	controller := NewControllerStd(needRetry)

	err := controller.Do(func() error {
		return nil
	})

	require.NoError(t, err, "expected no error")
}

func TestController_Do_RetrySuccess(t *testing.T) {
	needRetry := func(err error) bool {
		return true
	}
	controller := NewControllerStd(needRetry)
	attempts := 0

	err := controller.Do(func() error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary error")
		}
		return nil
	})

	require.NoError(t, err, "expected no error")
	assert.Equal(t, 2, attempts, "expected 2 attempts")
}

func TestController_Do_MaxRetries(t *testing.T) {
	needRetry := func(err error) bool {
		return true
	}
	controller := NewControllerStd(needRetry)
	attempts := 0

	err := controller.Do(func() error {
		attempts++
		return errors.New("temporary error")
	})

	require.Error(t, err, "expected an error")
	expectedAttempts := maxRetries + 1
	assert.Equal(t, expectedAttempts, attempts, "expected attempts to match maxRetries + 1")
	expectedErrMsg := "the attempts are over, the last error: temporary error"
	assert.Equal(t, expectedErrMsg, err.Error(), "expected error message to match")
}

func TestController_Do_NoRetryOnSpecificError(t *testing.T) {
	needRetry := func(err error) bool {
		return err.Error() != "no retry"
	}
	controller := NewControllerStd(needRetry)

	err := controller.Do(func() error {
		return errors.New("no retry")
	})

	require.Error(t, err, "expected an error")
	assert.Equal(t, "no retry", err.Error(), "expected error message to match")
}
