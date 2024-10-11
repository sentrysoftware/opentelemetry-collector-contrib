package bmchelixexporter

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry"
	"go.opentelemetry.io/collector/config/configretry"
)

// Config struct is used to store the configuration of the exporter
type Config struct {
	Endpoint                  string                       `mapstructure:"endpoint"`
	ApiKey                    string                       `mapstructure:"api_key"`
	Timeout                   time.Duration                `mapstructure:"timeout"`
	RetryConfig               configretry.BackOffConfig    `mapstructure:"retry_on_failure"`
	ResourceToTelemetryConfig resourcetotelemetry.Settings `mapstructure:"resource_to_telemetry_conversion"`
	MappingModel			  *MappingModelConfig           `mapstructure:"mapping_model"`
}

// EntityMappingConfig struct is used to store the configuration of the BMC Helix entity mapping
type MappingModelConfig struct {
	Force bool 				 `mapstructure:"force"`
	EntityMappings []EntityMappingConfig `mapstructure:"entity_mappings"`
	StateAttributes []string             `mapstructure:"state_attributes"`
}

type EntityMappingConfig struct {
	EntityTypeId         string            `mapstructure:"entity_type_id"`
	MetricPatterns       []string          `mapstructure:"metric_patterns"`
	RequiredAttrsForId   []string          `mapstructure:"required_attributes_for_id"`
	RequiredAttrsForName []string          `mapstructure:"required_attributes_for_name"`
	MetricAttributes     map[string]string `mapstructure:"metric_attributes"`
}

// validate function is used to validate the configuration
func (c *Config) Validate() error {
	if c.Endpoint == "" {
		return errors.New("endpoint is required")
	}
	if c.ApiKey == "" {
		return errors.New("api key is required")
	}
	if c.Timeout < 0 {
		return errors.New("timeout must be a positive integer")
	}

	// Validate EntityMappingConfig
	if (c.MappingModel != nil) {
		

	for _, entityMapping := range c.MappingModel.EntityMappings {
		if entityMapping.EntityTypeId == "" {
			return errors.New("entity_type_id is required")
		}
		if len(entityMapping.MetricPatterns) == 0 && len(entityMapping.MetricAttributes) == 0 {
			return errors.New("metric_patterns or metric_attributes is required")
		}

		if (len(entityMapping.MetricPatterns) > 0) {
			for _, pattern := range entityMapping.MetricPatterns {
				_, err := regexp.Compile(pattern)
				if err != nil {
					return fmt.Errorf("%q, %w", pattern, err)
				}
			}
		}

		if len(entityMapping.RequiredAttrsForId) == 0 {
			return errors.New("required_attributes_for_id is required")
		}
		if len(entityMapping.RequiredAttrsForName) == 0 {
			return errors.New("required_attributes_for_name is required")
		}
	}
}
	return nil
}
