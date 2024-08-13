package items_test

import (
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/common/items"
	"github.com/stretchr/testify/assert"
)

func TestMetricsJSONMarshaling(t *testing.T) {
	m := items.Metrics{
		ID:    "test_id",
		MType: "gauge",
		Delta: nil,
		Value: new(float64),
	}
	*m.Value = 123.45

	data, err := m.MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{"id":"test_id","type":"gauge","value":123.45}`, string(data))

	var unmarshaledMetrics items.Metrics
	err = unmarshaledMetrics.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, m, unmarshaledMetrics)
}

func TestMetricsSliceJSONMarshaling(t *testing.T) {
	metricsSlice := items.MetricsSlice{
		{
			ID:    "test_id_1",
			MType: "gauge",
			Delta: nil,
			Value: new(float64),
		},
		{
			ID:    "test_id_2",
			MType: "counter",
			Delta: new(int64),
			Value: nil,
		},
	}
	*metricsSlice[0].Value = 123.45
	*metricsSlice[1].Delta = 42

	data, err := metricsSlice.MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `[{"id":"test_id_1","type":"gauge","value":123.45},{"id":"test_id_2","type":"counter","delta":42}]`, string(data))

	var unmarshaledMetricsSlice items.MetricsSlice
	err = unmarshaledMetricsSlice.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, metricsSlice, unmarshaledMetricsSlice)
}
