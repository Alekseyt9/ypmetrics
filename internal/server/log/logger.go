// Package log provides utilities for logging within the server.
package log

import (
	"log/slog"
	"os"
)

// Logger defines an interface for logging messages at various levels of severity.
type Logger interface {
	// Debug logs a message at the Debug level.
	// Parameters:
	//   - template: the message template to log
	//   - keysAndValues: optional key-value pairs to include with the message
	Debug(template string, keysAndValues ...interface{})

	// Info logs a message at the Info level.
	// Parameters:
	//   - template: the message template to log
	//   - keysAndValues: optional key-value pairs to include with the message
	Info(template string, keysAndValues ...interface{})

	// Warn logs a message at the Warn level.
	// Parameters:
	//   - template: the message template to log
	//   - keysAndValues: optional key-value pairs to include with the message
	Warn(template string, keysAndValues ...interface{})

	// Error logs a message at the Error level.
	// Parameters:
	//   - template: the message template to log
	//   - keysAndValues: optional key-value pairs to include with the message
	Error(template string, keysAndValues ...interface{})
}

func NewSlogLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger
}

type NoOpLogger struct{}

func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

func (l *NoOpLogger) Debug(template string, keysAndValues ...interface{}) {}
func (l *NoOpLogger) Info(template string, keysAndValues ...interface{})  {}
func (l *NoOpLogger) Warn(template string, keysAndValues ...interface{})  {}
func (l *NoOpLogger) Error(template string, keysAndValues ...interface{}) {}
