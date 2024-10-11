package bmchelixexporter

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/bmchelixexporter/internal/mapping/hardware"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	conventions "go.opentelemetry.io/collector/semconv/v1.27.0"
	"go.uber.org/zap"
)

// BmcHelixMetric represents the structure of the payload that will be sent to BMC Helix
type BmcHelixMetric struct {
	Labels  map[string]string `json:"labels"`
	Samples []BmcHelixSample  `json:"samples"`
}

// BmcHelixSample represents the individual sample for a metric
type BmcHelixSample struct {
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

type BmcHelixMetricsProducer struct {
	osHostname string
	logger *zap.Logger
	mappingResolver *hardware.HardwareMappingResolver
}

// ProduceHelixPayload takes the OpenTelemetry metrics and converts them into the BMC Helix metric format
func (mp *BmcHelixMetricsProducer) ProduceHelixPayload(metrics pmetric.Metrics) ([]BmcHelixMetric, error) {
	var helixMetrics []BmcHelixMetric

	// Iterate through each resource metrics
	rmetrics := metrics.ResourceMetrics()
	for i := 0; i < rmetrics.Len(); i++ {
		resourceMetric := rmetrics.At(i)
		resource := resourceMetric.Resource()

		mp.logger.Debug("Resource", zap.Any("resource", resource))
		// Extract resource-level attributes (e.g., "host.name", "service.instance.id")
		resourceAttrs := extractResourceAttributes(resource)

		mp.logger.Debug("Resource attributes", zap.Any("attributes", resourceAttrs))
		// Iterate through each scope metric within the resource
		scopeMetrics := resourceMetric.ScopeMetrics()
		for j := 0; j < scopeMetrics.Len(); j++ {
			scopeMetric := scopeMetrics.At(j)

			// Iterate through each individual metric
			metrics := scopeMetric.Metrics()
			for k := 0; k < metrics.Len(); k++ {
				metric := metrics.At(k)

				// Create the payload for each metric
				helixMetric, err := mp.createHelixMetric(metric, resourceAttrs)
				if err != nil {
					mp.logger.Warn("Failed to create Helix metric", zap.Error(err))
					continue
				}
				helixMetrics = append(helixMetrics, helixMetric)
			}
		}
	}
	return helixMetrics, nil
}

// createHelixMetric converts a single OpenTelemetry metric into a BmcHelixMetric payload
func (mp *BmcHelixMetricsProducer) createHelixMetric(metric pmetric.Metric, resourceAttrs map[string]string) (BmcHelixMetric, error) {
	labels := make(map[string]string)
	labels["source"] = "OTEL"

	// Add resource attributes as labels
	for k, v := range resourceAttrs {
		labels[k] = v
	}

	// Set the metric unit
	labels["unit"] = metric.Unit()

	// Set the host type
	labels["hostType"] = "server"

	// Indicates the monitor in the hierarchy that is mapped to the device
	labels["isDeviceMappingEnabled"] = "true"

	// Samples to hold the metric values
	var samples []BmcHelixSample

	// Handle different types of metrics (sum, gauge, histogram, etc.)
	switch metric.Type() {
	case pmetric.MetricTypeSum:
		dataPoints := metric.Sum().DataPoints()
		for i := 0; i < dataPoints.Len(); i++ {
			dp := dataPoints.At(i)
			samples = mp.updateMetricInformation(samples, dp, labels, metric, resourceAttrs)

		}
	case pmetric.MetricTypeGauge:
		dataPoints := metric.Gauge().DataPoints()
		for i := 0; i < dataPoints.Len(); i++ {
			dp := dataPoints.At(i)
			samples = mp.updateMetricInformation(samples, dp, labels, metric, resourceAttrs)
		}
	}

	if len(samples) == 0 {
		return BmcHelixMetric{}, fmt.Errorf("no samples found for metric %s", metric.Name())
	}

	return BmcHelixMetric{
		Labels:  labels,
		Samples: samples,
	}, nil
}

// Updates the metric information for the BMC Helix payload
func (mp *BmcHelixMetricsProducer)  updateMetricInformation(samples []BmcHelixSample, dp pmetric.NumberDataPoint, labels map[string]string, metric pmetric.Metric, resourceAttrs map[string]string) []BmcHelixSample {

	// Update the entity information for the BMC Helix payload
	err := mp.updateEntityInformation(labels,  metric.Name(), resourceAttrs, dp.Attributes().AsRaw())
	if err != nil {
		return samples
	}

	// Update the metric name for the BMC Helix payload
	updateMetricName(labels, metric, dp)

	return append(samples, newSample(dp))
}

// Update the entity information for the BMC Helix payload
func (mp *BmcHelixMetricsProducer) updateEntityInformation(labels map[string]string, metricName string, resourceAttrs map[string]string, metricAttrs map[string]any) (error) {
    // If the entityTypeId and entityName are already set, return early
    if labels["entityTypeId"] != "" {
        return nil
    }

    // Try to get the hostname from resource attributes first
    hostname, found := resourceAttrs[conventions.AttributeHostName]
    if !found || hostname == "" {
        // Fallback to metric attributes if not found or empty in resource attributes
        if maybeHostname, ok := metricAttrs[conventions.AttributeHostName].(string); ok && maybeHostname != "" {
            hostname = maybeHostname
        } else {
            // Fallback to osHostname if hostname is not found in both places
            hostname = mp.osHostname
        }
    }

    // Add the hostname as a label (required for BMC Helix payload)
    labels["hostname"] = hostname

    // Convert metricAttrs from map[string]any to map[string]string for compatibility with the resolver
    stringMetricAttrs := make(map[string]string)
    for k, v := range metricAttrs {
        stringMetricAttrs[k] = fmt.Sprintf("%v", v)
		labels[k] = fmt.Sprintf("%v", v)
    }

	// Add the resource attributes to the metric attributes
	for k, v := range resourceAttrs {
		stringMetricAttrs[k] = v
	}

    // Use the mapping resolver to determine entityTypeId, entityId, and entityName
    entityTypeId, entityId, entityName, err := mp.mappingResolver.MapEntityTypeAndAttributes(metricName, stringMetricAttrs)
    if err != nil {
        mp.logger.Warn("Failed to map entity type and attributes", zap.String("metricName", metricName), zap.Error(err))
        return err
    }

    // Set the entityTypeId, entityId, and entityName in labels
    labels["entityTypeId"] = entityTypeId
    labels["entityId"] = fmt.Sprintf("%s:%s:%s:%s", labels["source"], labels["hostname"], entityId, entityName)
    labels["entityName"] = entityName
	return nil
}


// newSample creates a new BmcHelixSample from the OpenTelemetry data point
func newSample(dp pmetric.NumberDataPoint) BmcHelixSample {
    var value float64
    switch dp.ValueType() {
    case pmetric.NumberDataPointValueTypeDouble:
        value = dp.DoubleValue()
    case pmetric.NumberDataPointValueTypeInt:
        value = float64(dp.IntValue()) // convert int to float for consistency
    }
    
    return BmcHelixSample{
        Value:     value,
        Timestamp: dp.Timestamp().AsTime().Unix() * 1000,
    }
}

// Build the metric name for the BMC Helix payload
func updateMetricName(labels map[string]string, metric pmetric.Metric, dp pmetric.NumberDataPoint) {
	// Update the metric name to include the labels
	// For example, if the original metric name is "hw.status" and the labels are "state=ok", the new metric name will be "hw.status.ok"
	// This is to ensure that the metric name is unique in BMC Helix

	if labels["metricName"] != "" {
		return
	}

	metricName := metric.Name()

	// Append known hw metric attributes to the metric name
	// E.g. "hw.status" -> "hw.status.ok"
	for _, attribute := range hardware.KnownHwMetricAttributeKeys {
		metricName = appendAttributeIfExists(metricName, dp, attribute)
	}

	labels["metricName"] = metricName
}

// appendAttributeIfExists appends the attribute to the metric name if it exists
func appendAttributeIfExists(metricName string, dp pmetric.NumberDataPoint, attributeKey string) string {
	if attr, ok := dp.Attributes().Get(attributeKey); ok {
		metricName = fmt.Sprintf("%s.%s", metricName, attr.AsString())
	}
	return metricName
}

// extractResourceAttributes extracts the resource attributes from OpenTelemetry resource data
func extractResourceAttributes(resource pcommon.Resource) map[string]string {
	attributes := make(map[string]string)

	resource.Attributes().Range(func(k string, v pcommon.Value) bool {
		attributes[k] = v.AsString()
		return true
	})

	return attributes
}
