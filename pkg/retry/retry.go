// Package retry provides mechanisms for retrying operations with configurable delays and retry logic.
package retry

import (
	"fmt"
	"time"
)

const (
	maxRetries  = 3
	retryDelay1 = 1 * time.Second
	retryDelay2 = 3 * time.Second
	retryDelay3 = 5 * time.Second
)

// Controller handles the retry logic for operations.
type Controller struct {
	Retries   int
	Delays    []time.Duration
	NeedRetry func(error) bool
}

// NewControllerStd creates a new Controller with standard retry settings.
// Parameters:
//   - needRetry: a function to determine if an error should trigger a retry
//
// Returns a pointer to a Controller.
func NewControllerStd(needRetry func(error) bool) *Controller {
	return NewController(
		maxRetries, []time.Duration{retryDelay1, retryDelay2, retryDelay3}, needRetry)
}

// NewController creates a new Controller with custom retry settings.
// Parameters:
//   - retries: maximum number of retries
//   - delays: slice of durations between retries
//   - needRetry: a function to determine if an error should trigger a retry
//
// Returns a pointer to a Controller.
func NewController(retries int, delays []time.Duration, needRetry func(error) bool) *Controller {
	return &Controller{
		Retries:   retries,
		Delays:    delays,
		NeedRetry: needRetry,
	}
}

// Do executes the given function with retry logic.
// Parameters:
//   - f: the function to execute
//
// Returns an error if the function ultimately fails after the specified retries.
func (rc *Controller) Do(f func() error) error {
	for attempt := 0; attempt <= rc.Retries; attempt++ {
		err := f()
		if err == nil {
			return nil
		}

		if !rc.NeedRetry(err) {
			return err
		}

		if attempt == rc.Retries {
			return fmt.Errorf("the attempts are over, the last error: %w", err)
		}

		time.Sleep(rc.Delays[attempt])
	}
	return nil
}
