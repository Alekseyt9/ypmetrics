package logger

import (
	"net/http"
	"time"
)

// Абстрагируемся от реализации.
type Logger interface {
	Debug(template string, keysAndValues ...interface{})
	Info(template string, keysAndValues ...interface{})
	Warn(template string, keysAndValues ...interface{})
	Error(template string, keysAndValues ...interface{})
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
		log.Info("query: ",
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
