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

type Controller struct {
	Retries   int
	Delays    []time.Duration
	NeedRetry func(error) bool
}

func NewControllerStd(needRetry func(error) bool) *Controller {
	return NewController(
		maxRetries, []time.Duration{retryDelay1, retryDelay2, retryDelay3}, needRetry)
}

func NewController(retries int, delays []time.Duration, needRetry func(error) bool) *Controller {
	return &Controller{
		Retries:   retries,
		Delays:    delays,
		NeedRetry: needRetry,
	}
}

func (rc *Controller) Do(f func() error) error {
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
