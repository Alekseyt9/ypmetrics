package utils

import (
	"net/http"
	"strings"
)

type MetricInfo struct {
	Name  string
	Value string
}

type URLParseError struct {
	Message string
	Status  int
}

func (e *URLParseError) Error() string {
	return e.Message
}

func ParseURL(url string, prefix string) (MetricInfo, error) {
	trimPath := strings.TrimPrefix(url, prefix)

	if trimPath == "" {
		return MetricInfo{}, &URLParseError{
			Message: "metric name is needed",
			Status:  http.StatusNotFound,
		}
	}

	parts := strings.Split(trimPath, "/")
	res := MetricInfo{}
	if len(parts) == 2 {
		res.Name = parts[0]
		res.Value = parts[1]
		return res, nil
	}

	return MetricInfo{}, &URLParseError{
		Message: "Invalid URL format",
		Status:  http.StatusBadRequest,
	}
}
