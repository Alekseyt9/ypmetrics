package common

import (
	"fmt"
	"time"
)

type RetryController struct {
	Retries   int
	Delays    []time.Duration
	NeedRetry func(error) bool
}

func NewRetryControllerStd(needRetry func(error) bool) *RetryController {
	return NewRetryController(
		3, []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}, needRetry)
}

func NewRetryController(retries int, delays []time.Duration, needRetry func(error) bool) *RetryController {
	return &RetryController{
		Retries:   retries,
		Delays:    delays,
		NeedRetry: needRetry,
	}
}

func (rc *RetryController) Do(f func() error) error {
	attempt := 0
	for {
		err := f()
		if err == nil {
			return nil
		}

		if !rc.NeedRetry(err) {
			return err
		}

		if attempt >= rc.Retries {
			return fmt.Errorf("the attempts are over, the last error: %w", err)
		}
		time.Sleep(rc.Delays[attempt])
		attempt++
	}
}
