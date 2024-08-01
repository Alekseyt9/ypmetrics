package log_test

import (
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/server/log"
	"github.com/stretchr/testify/assert"
)

func TestNoOpLogger(t *testing.T) {
	logger := log.NewNoOpLogger()

	assert.NotPanics(t, func() {
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")
	})
}

func TestSlogLogger(t *testing.T) {
	logger := log.NewSlogLogger()

	assert.NotPanics(t, func() {
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")
	})
}
