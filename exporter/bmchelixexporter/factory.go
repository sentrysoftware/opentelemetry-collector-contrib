package bmchelixexporter

import (
	"context"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/bmchelixexporter/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Create bmchelixexporter factory
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithMetrics(createMetricsExporter, metadata.MetricsStability),
	)
}

// CreateDefaultConfig creates the default configuration for BMC Helix exporter.
func createDefaultConfig() component.Config {
	return &Config{
		Timeout:                   10 * time.Second,
		ResourceToTelemetryConfig: resourcetotelemetry.Settings{Enabled: true},
	}
}

// CreateMetricsExporter creates a BMC Helix exporter.
func createMetricsExporter(
	ctx context.Context,
	set exporter.Settings,
	config component.Config,
) (exporter.Metrics, error) {
	cfg := config.(*Config)
	bmcHelixExp, err := newBmcHelixExporter(cfg, set)

	if err != nil {
		return nil, err
	}

	me, err := exporterhelper.NewMetricsExporter(
		ctx,
		set,
		cfg,
		bmcHelixExp.pushMetrics,
		exporterhelper.WithTimeout(exporterhelper.TimeoutConfig{Timeout: 0}),
		exporterhelper.WithRetry(cfg.RetryConfig),
		exporterhelper.WithStart(bmcHelixExp.start),
	)

	if err != nil {
		return nil, err
	}

	return resourcetotelemetry.WrapMetricsExporter(cfg.ResourceToTelemetryConfig, me), nil
}
