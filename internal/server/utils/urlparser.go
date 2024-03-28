package utils

import (
	"net/http"
	"strings"
)

type MetricInfo struct {
	Name  string
	Value string
}

type UrlParseError struct {
	Message string
	Status  int
}

func (e *UrlParseError) Error() string {
	return e.Message
}

func ParseUrl(url string, prefix string) (MetricInfo, error) {
	trimPath := strings.TrimPrefix(url, prefix)

	if trimPath == "" {
		return MetricInfo{}, &UrlParseError{
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

	return MetricInfo{}, &UrlParseError{
		Message: "Invalid URL format",
		Status:  http.StatusBadRequest,
	}
}
