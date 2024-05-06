package logger

import (
	"log/slog"
	"os"
)

type slogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger() Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return &slogLogger{logger: logger}
}

func (z *slogLogger) Debug(msg string, keysAndValues ...interface{}) {
	z.logger.Debug(msg, keysAndValues...)
}

func (z *slogLogger) Info(msg string, keysAndValues ...interface{}) {
	z.logger.Info(msg, keysAndValues...)
}

func (z *slogLogger) Warn(msg string, keysAndValues ...interface{}) {
	z.logger.Warn(msg, keysAndValues...)
}

func (z *slogLogger) Error(msg string, keysAndValues ...interface{}) {
	z.logger.Error(msg, keysAndValues...)
}
