package mapping

import (
	"fmt"
	"regexp"
	"strings"

	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

var (
	ScopeName = "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/bmchelixexporter"
)

type MappingResolver struct {
	Logger *zap.Logger
	MappingModels map[string]MappingModel
}

// ResolveMetricEntity maps the entityTypeId, entityId, and entityName based on the mapping model.
func (mr *MappingResolver) ResolveMetricEntity(metricName string, metricAttrs map[string]string) (string, string, string, error) {
	var entityTypeId, entityId, entityName string

	resolved := false
	for id, mappingModel := range mr.MappingModels {

		if resolved {
			mr.Logger.Debug("Resolved entity through mapping model", zap.String("id", id),  zap.String("metricName", metricName), zap.String("entityTypeId", entityTypeId), zap.String("entityId", entityId), zap.String("entityName", entityName))
			break
		}

		// Check for a matching entity mapping
		for _, mapping := range mappingModel.EntityMappings {
			// Check if the metric name starts with any of the defined prefixes
			if matchesMetricName(metricName, mapping.MetricPatterns) || matchesMetricAttributes(metricAttrs, mapping.MetricAttributes) {
				entityTypeId = mapping.EntityTypeId

				// Build the entityId from the specified attributes
				entityId = buildValueFromAttrs(mapping.RequiredAttrsForId, metricAttrs)
				if (entityId == "") {
					mr.Logger.Warn("Failed to build entityId", zap.String("metricName", metricName), zap.Any("metricAttrs", metricAttrs))
					continue
				}
				// Build the entityName from the specified attributes
				entityName = buildValueFromAttrs(mapping.RequiredAttrsForName, metricAttrs)
				if (entityName == "") {
					mr.Logger.Warn("Failed to build entityName", zap.String("metricName", metricName), zap.Any("metricAttrs", metricAttrs))
					continue
				}
				resolved = true
			}
		}

	}

	// Fallback to default if no mapping was found
	if entityTypeId == "" || entityId == "" || entityName == "" {
		mr.Logger.Warn("No mapping found for metric", zap.String("metricName", metricName), zap.Any("metricAttrs", metricAttrs))
		return "", "", "", fmt.Errorf("no mapping found for metric %s with attributes %v", metricName, metricAttrs) 
	}

	return entityTypeId, entityId, entityName, nil
}

// matchesMetricName checks if the metric name matches any of the compiled regex patterns.
func matchesMetricName(metricName string, patterns []*regexp.Regexp) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(metricName) {
			return true
		}
	}
	return false
}

// Helper to check if the metric attributes match any of the defined attributes.
func matchesMetricAttributes(metricAttrs map[string]string, requiredAttrs map[string]string) bool {
	if (len(requiredAttrs) == 0) {
		return false
	}
	for key, value := range requiredAttrs {
		if metricAttrs[key] != value {
			return false
		}
	}
	return true
}

// Helper to build the entityId or entityName from specified attributes.
func buildValueFromAttrs(requiredAttrs []string, metricAttrs map[string]string) string {
	var elements []string

	// Build from specified attributes
	for _, attr := range requiredAttrs {
		if val, ok := metricAttrs[attr]; ok {
			elements = append(elements, val)
		}
	}

	// Join the elements into a single string
	return strings.Join(elements, "-")
}

// ResolveMetricName resolves the metric name based on the mapping model.
func (mr *MappingResolver) ResolveMetricName(labels map[string]string, metric pmetric.Metric, dp pmetric.NumberDataPoint) {
	// Update the metric name to include the labels
	// For example, if the original metric name is "hw.status" and the labels are "state=ok", the new metric name will be "hw.status.ok"
	// This is to ensure that the metric name is unique in BMC Helix

	if labels["metricName"] != "" {
		return
	}

	metricName := metric.Name()

	// Append known hw metric attributes to the metric name
	// E.g. "hw.status" -> "hw.status.ok"
	for _, mappingModel := range mr.MappingModels {
		for _, attribute := range mappingModel.StateAttributes {
			metricName = appendAttributeIfExists(metricName, dp, attribute)
		}
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