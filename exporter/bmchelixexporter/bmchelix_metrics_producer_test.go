package bmchelixexporter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

// Mock data generation for testing
func generateMockMetrics(dpCreator func(metric pmetric.Metric) pmetric.NumberDataPoint) pmetric.Metrics {
	metrics := pmetric.NewMetrics()
	rm := metrics.ResourceMetrics().AppendEmpty()
	il := rm.ScopeMetrics().AppendEmpty().Metrics()
	metric := il.AppendEmpty()
	metric.SetName("test_metric")
	metric.SetDescription("This is a test metric")
	metric.SetUnit("ms")
	dp := dpCreator(metric)
	dp.Attributes().PutStr("entityId", "test-entity")
	dp.Attributes().PutStr("entityTypeId", "test-entity-type-id")
	dp.Attributes().PutStr("entityName", "test-entity-Name")
	dp.SetTimestamp(1634236000000000) // Example timestamp
	dp.SetDoubleValue(42.0)
	return metrics
}

// Test for the ProduceHelixPayload method
func TestProduceHelixPayload(t *testing.T) {
	t.Parallel()

	sample := BmcHelixSample{
		Value:     42,
		Timestamp: 1634236000,
	}

	metric := BmcHelixMetric{
		Labels: map[string]string{
			"isDeviceMappingEnabled": "true",
			"entityTypeId":           "test-entity-type-id",
			"entityName":             "test-entity-Name",
			"source":                 "OTEL",
			"unit":                   "ms",
			"hostType":               "server",
			"metricName":             "test_metric",
			"hostname":               "test-hostname",
			"entityId":               "OTEL:test-hostname:test-entity:test-entity-Name",
		},
		Samples: []BmcHelixSample{sample},
	}

	expectedPayload := []BmcHelixMetric{metric}

	producer := &BmcHelixMetricsProducer{
		osHostname: "test-hostname",
		logger:     zap.NewExample(),
	}

	tests := []struct {
		name                string
		generateMockMetrics func() pmetric.Metrics
		expectedPayload     []BmcHelixMetric
	}{
		{
			name: "SetGauge",
			generateMockMetrics: func() pmetric.Metrics {
				return generateMockMetrics(func(metric pmetric.Metric) pmetric.NumberDataPoint {
					return metric.SetEmptyGauge().DataPoints().AppendEmpty()
				})
			},
			expectedPayload: expectedPayload,
		},
		{
			name: "SetSum",
			generateMockMetrics: func() pmetric.Metrics {
				return generateMockMetrics(func(metric pmetric.Metric) pmetric.NumberDataPoint {
					return metric.SetEmptySum().DataPoints().AppendEmpty()
				})
			},
			expectedPayload: expectedPayload,
		},
		{
			name: "emptyPayload",
			generateMockMetrics: func() pmetric.Metrics {
				return pmetric.NewMetrics()
			},
			expectedPayload: []BmcHelixMetric{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMetrics := tt.generateMockMetrics()
			payload, err := producer.ProduceHelixPayload(mockMetrics)
			assert.NoError(t, err, "Expected no error during payload production")
			assert.NotNil(t, payload, "Payload should not be nil")

			assert.Equal(t, tt.expectedPayload, payload, "Payload should match the expected payload")
		})
	}
}
