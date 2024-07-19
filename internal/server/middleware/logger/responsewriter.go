// Package logger defines a custom response writer used within the logging middleware
// to capture and record response details such as status code and response size.
package logger

import (
	"net/http"
)

// ResponseWriter interface is redeclared here for clarity, outlining methods that
// manipulate the HTTP response.
type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write writes data to the wrapped ResponseWriter and updates the size tracked in responseData.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader writes the HTTP status code to the wrapped ResponseWriter and records it in responseData.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
