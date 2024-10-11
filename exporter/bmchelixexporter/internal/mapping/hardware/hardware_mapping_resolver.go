package hardware

import (
	"fmt"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

// PredefinedMappings is a list of predefined mappings for hardware metrics.
type HardwareMappingResolver struct {
	logger *zap.Logger
}

// New instance of HardwareMappingResolver
func NewHardwareMappingResolver(logger *zap.Logger) *HardwareMappingResolver {
	return &HardwareMappingResolver{
		logger: logger,
	}
}

// MapEntityTypeAndAttributes maps the entityTypeId, entityId, and entityName based on the hardware mappings.
func (mr *HardwareMappingResolver) MapEntityTypeAndAttributes(metricName string, metricAttrs map[string]string) (string, string, string, error) {
	var entityTypeId, entityId, entityName string

	// Check for a matching entity mapping
	for _, mapping := range PredefinedMappings {
		// Check if the metric name starts with any of the defined prefixes
		if matchesMetricName(metricName, mapping.MetricPatterns) || matchesMetricAttributes(metricAttrs, mapping.MetricAttributes) {
			entityTypeId = mapping.EntityTypeId

			// Build the entityId from the specified attributes
			entityId = buildValueFromAttrs(mapping.RequiredAttrsForId, metricAttrs)
			if (entityId == "") {
				mr.logger.Warn("Failed to build entityId", zap.String("metricName", metricName), zap.Any("metricAttrs", metricAttrs))
				continue
			}
			// Build the entityName from the specified attributes
			entityName = buildValueFromAttrs(mapping.RequiredAttrsForName, metricAttrs)
			if (entityName == "") {
				mr.logger.Warn("Failed to build entityName", zap.String("metricName", metricName), zap.Any("metricAttrs", metricAttrs))
				continue
			}
			break
		}
	}

	// Fallback to default if no mapping was found
	if entityTypeId == "" || entityId == "" || entityName == "" {
		mr.logger.Warn("No mapping found for metric", zap.String("metricName", metricName), zap.Any("metricAttrs", metricAttrs))
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
