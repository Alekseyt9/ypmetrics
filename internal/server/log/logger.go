package log

import (
	"log/slog"
	"os"
)

type Logger interface {
	Debug(template string, keysAndValues ...interface{})
	Info(template string, keysAndValues ...interface{})
	Warn(template string, keysAndValues ...interface{})
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
