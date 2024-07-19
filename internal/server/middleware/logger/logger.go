// Package logger provides middleware that logs HTTP requests and responses,
// including the method, URI, status, response size, and the duration of the request.
package logger

import (
	"net/http"
	"time"

	"github.com/Alekseyt9/ypmetrics/internal/server/log"
)

// WithLogging returns a middleware handler that logs details about each HTTP request and its response.
// It logs the URI, method, duration, status, and response size.
func WithLogging(h http.Handler, log log.Logger) http.Handler {
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
