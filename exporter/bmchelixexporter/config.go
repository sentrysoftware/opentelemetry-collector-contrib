package bmchelixexporter

import (
	"errors"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry"
	"go.opentelemetry.io/collector/config/configretry"
)

// Config struct is used to store the configuration of the exporter
type Config struct {
	Endpoint                  string                       `mapstructure:"endpoint"`
	Api                       ApiConfig                    `mapstructure:"api"`
	Timeout                   time.Duration                `mapstructure:"timeout"`
	RetryConfig               configretry.BackOffConfig    `mapstructure:"retry_on_failure"`
	ResourceToTelemetryConfig resourcetotelemetry.Settings `mapstructure:"resource_to_telemetry_conversion"`
}

// ApiConfig struct is used to store the configuration of the BMC Helix API
type ApiConfig struct {
	AccessKey       string `mapstructure:"access_key"`
	AccessSecretKey string `mapstructure:"access_secret_key"`
	TenantId        string `mapstructure:"tenant_id"`
}

// validate function is used to validate the configuration
func (c *Config) Validate() error {
	if c.Endpoint == "" {
		return errors.New("endpoint is required")
	}
	if c.Api.AccessKey == "" {
		return errors.New("api access key is required")
	}
	if c.Api.AccessSecretKey == "" {
		return errors.New("api access secret key is required")
	}
	if c.Api.TenantId == "" {
		return errors.New("api tenant id is required")
	}
	if c.Timeout < 0 {
		return errors.New("timeout must be a positive integer")
	}
	return nil
}
