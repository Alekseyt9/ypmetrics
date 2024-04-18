package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// абстрагируемся от реализации
type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

func WithLogging(h http.Handler, log Logger) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)
		log.Infof(
			"uri %s, method %s, duration %v, status %s, size %v",
			r.RequestURI, r.Method, duration,
			responseData.status, responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}

func NewLogger() Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()
	return sugar
}
