package logger

import (
	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return &zapLogger{logger: logger}
}

func anyToZapFields(keysAndValues ...interface{}) []zap.Field {
	var fields []zap.Field
	for i := 0; i < len(keysAndValues); i += 2 {
		key, val := keysAndValues[i], keysAndValues[i+1]
		if i+1 < len(keysAndValues) {
			fields = append(fields, zap.Any(key.(string), val))
		}
	}
	return fields
}

func (z *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	z.logger.Debug(msg, anyToZapFields(keysAndValues...)...)
}

func (z *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	z.logger.Info(msg, anyToZapFields(keysAndValues...)...)
}

func (z *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	z.logger.Warn(msg, anyToZapFields(keysAndValues...)...)
}

func (z *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	z.logger.Error(msg, anyToZapFields(keysAndValues...)...)
}
