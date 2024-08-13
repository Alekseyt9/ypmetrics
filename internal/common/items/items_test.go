package items_test

import (
	"reflect"
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/common/items"
)

func TestMetricsSlice_ToMetricItems(t *testing.T) {
	gaugeValue := 123.45
	counterValue := int64(678)

	input := items.MetricsSlice{
		{ID: "gauge1", MType: "gauge", Value: &gaugeValue},
		{ID: "counter1", MType: "counter", Delta: &counterValue},
	}

	expected := items.MetricItems{
		Counters: []items.CounterItem{{Name: "counter1", Value: 678}},
		Gauges:   []items.GaugeItem{{Name: "gauge1", Value: 123.45}},
	}

	output := input.ToMetricItems()

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("ToMetricItems() = %v, want %v", output, expected)
	}
}

func TestMetricItems_ToMetricsSlice(t *testing.T) {
	gaugeValue := 123.45
	counterValue := int64(678)

	input := items.MetricItems{
		Counters: []items.CounterItem{{Name: "counter1", Value: 678}},
		Gauges:   []items.GaugeItem{{Name: "gauge1", Value: 123.45}},
	}

	expected := items.MetricsSlice{
		{ID: "counter1", MType: "counter", Delta: &counterValue},
		{ID: "gauge1", MType: "gauge", Value: &gaugeValue},
	}

	output := input.ToMetricsSlice()

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("ToMetricsSlice() = %v, want %v", output, expected)
	}
}
