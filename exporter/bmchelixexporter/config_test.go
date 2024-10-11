package bmchelixexporter

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/confmap/confmaptest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/bmchelixexporter/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	tests := []struct {
		id           component.ID
		expected     component.Config
		errorMessage string
	}{

		{
			id: component.NewIDWithName(metadata.Type, "helix1"),
			expected: &Config{
				Endpoint:    "https://helix1:8080",
				ApiKey:      "api_key",
				Timeout:     10 * time.Second,
				RetryConfig: configretry.BackOffConfig{},
				ResourceToTelemetryConfig: resourcetotelemetry.Settings{
					Enabled: true,
				},
			},
		},
		{
			id: component.NewIDWithName(metadata.Type, "helix2"),
			expected: &Config{
				Endpoint: "https://helix2:8080",
				ApiKey:   "api_key",
				Timeout:  20 * time.Second,
				RetryConfig: configretry.BackOffConfig{
					Enabled:             true,
					InitialInterval:     5 * time.Second,
					RandomizationFactor: 0.5,
					Multiplier:          8.0,
					MaxInterval:         1 * time.Minute,
					MaxElapsedTime:      8 * time.Minute,
				},
				ResourceToTelemetryConfig: resourcetotelemetry.Settings{
					Enabled: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, sub.Unmarshal(cfg))

			assert.NoError(t, component.ValidateConfig(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		err    string
	}{
		{
			name: "valid_config",
			config: &Config{
				Endpoint: "https://helix:8080",
				ApiKey:   "api_key",
				Timeout:  10 * time.Second,
			},
		},
		{
			name: "invalid_config1",
			config: &Config{
				ApiKey: "api_key",
			},
			err: "endpoint is required",
		},
		{
			name: "invalid_config2",
			config: &Config{
				Endpoint: "https://helix:8080",
			},
			err: "api key is required",
		},
		{
			name: "invalid_config3",
			config: &Config{
				Endpoint: "https://helix:8080",
				ApiKey:   "api_key",
				Timeout:  -1,
			},
			err: "timeout must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != "" {
				err := tt.config.Validate()
				assert.Error(t, err)
				assert.Equal(t, tt.err, err.Error())
			} else {
				assert.NoError(t, tt.config.Validate())
			}
		})
	}
}

func TestValidateMappingModelConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		err    string
	}{
		{
			name: "valid_config",
			config: &Config{
				Endpoint: "https://helix:8080",
				ApiKey:   "api_key",
				Timeout:  10 * time.Second,
				MappingModel: &MappingModelConfig{
					Force: false,
					EntityMappings: []EntityMappingConfig{
						{
							EntityTypeId:         "system_memory",
							MetricPatterns:       []string{`system\.memory\.`},
							RequiredAttrsForId:   []string{"id"},
							RequiredAttrsForName: []string{"id"},
							MetricAttributes:     map[string]string{},
						},
					},
					StateAttributes: []string{"state"},
				},
			},
		},
		{
			name: "invalid_config1",
			config: &Config{
				Endpoint: "https://helix:8080",
				ApiKey:   "api_key",
				Timeout:  10 * time.Second,
				MappingModel: &MappingModelConfig{
					Force: false,
					EntityMappings: []EntityMappingConfig{
						{
							MetricPatterns:       []string{`system\.memory\.`},
							RequiredAttrsForId:   []string{"id"},
							RequiredAttrsForName: []string{"id"},
							MetricAttributes:     map[string]string{},
						},
					},
					StateAttributes: []string{"state"},
				},
			},
			err: "entity_type_id is required",
		},
		{
			name: "invalid_config2",
			config: &Config{
				Endpoint: "https://helix:8080",
				ApiKey:   "api_key",
				Timeout:  10 * time.Second,
				MappingModel: &MappingModelConfig{
					Force: false,
					EntityMappings: []EntityMappingConfig{
						{
							EntityTypeId:         "system_memory",
							MetricPatterns:       []string{`system\.memory\.`},
							RequiredAttrsForName: []string{"id"},
							MetricAttributes:     map[string]string{},
						},
					},
					StateAttributes: []string{"state"},
				},
			},
			err: "required_attributes_for_id is required",
		},
		{
			name: "invalid_config3",
			config: &Config{
				Endpoint: "https://helix:8080",
				ApiKey:   "api_key",
				Timeout:  10 * time.Second,
				MappingModel: &MappingModelConfig{
					Force: false,
					EntityMappings: []EntityMappingConfig{
						{
							EntityTypeId:         "system_memory",
							MetricPatterns:       []string{`system\.memory\.`},
							RequiredAttrsForId:   []string{"id"},
							MetricAttributes:     map[string]string{},
						},
					},
					StateAttributes: []string{"state"},
				},
			},
			err: "required_attributes_for_name is required",
		},
		{
			name: "invalid_config4",
			config: &Config{
				Endpoint: "https://helix:8080",
				ApiKey:   "api_key",
				Timeout:  10 * time.Second,
				MappingModel: &MappingModelConfig{
					Force: false,
					EntityMappings: []EntityMappingConfig{
						{
							EntityTypeId:         "system_memory",
							RequiredAttrsForId:   []string{"id"},
							RequiredAttrsForName:   []string{"id"},
						},
					},
					StateAttributes: []string{"state"},
				},
			},
			err: "metric_patterns or metric_attributes is required",
		},
		{
			name: "invalid_config5",
			config: &Config{
				Endpoint: "https://helix:8080",
				ApiKey:   "api_key",
				Timeout:  10 * time.Second,
				MappingModel: &MappingModelConfig{
					Force: false,
					EntityMappings: []EntityMappingConfig{
						{
							EntityTypeId:         "system_memory",
							MetricPatterns:       []string{`(`}, // Incorrect pattern: unbalanced parentheses
							RequiredAttrsForId:   []string{"id"},
							RequiredAttrsForName: []string{"id"},
							MetricAttributes:     map[string]string{},
						},
					},
					StateAttributes: []string{"state"},
				},
			},
			err: "\"(\", error parsing regexp: missing closing ): `(`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != "" {
				err := tt.config.Validate()
				assert.Error(t, err)
				assert.Equal(t, tt.err, err.Error())
			} else {
				assert.NoError(t, tt.config.Validate())
			}
		})
	}
}