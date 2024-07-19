// Package common provides common structures and functions for working with metrics.
package common

// Metrics represents the structure of a metric with its main attributes.
type Metrics struct {
	ID    string   `json:"id"`              // ID of the metric
	MType string   `json:"type"`            // Type of the metric: gauge or counter
	Delta *int64   `json:"delta,omitempty"` // Metric value for counter type
	Value *float64 `json:"value,omitempty"` // Metric value for gauge type
}

// MetricsSlice represents a slice of Metrics structures.
//
//easyjson:json
type MetricsSlice []Metrics
