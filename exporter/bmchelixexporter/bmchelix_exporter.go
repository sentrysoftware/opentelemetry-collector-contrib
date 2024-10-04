package bmchelixexporter

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type bmcHelixExporter struct {
	config            *Config
	logger            *zap.Logger
	version           string
	telemetrySettings component.TelemetrySettings
}

func newBmcHelixExporter(config *Config, createSettings exporter.Settings) (*bmcHelixExporter, error) {
	if config == nil {
		return nil, errors.New("nil config")
	}

	return &bmcHelixExporter{
		config:            config,
		version:           createSettings.BuildInfo.Version,
		logger:            createSettings.Logger,
		telemetrySettings: createSettings.TelemetrySettings,
	}, nil
}

func (be *bmcHelixExporter) pushMetrics(ctx context.Context, md pmetric.Metrics) error {
	return nil
}

func (be *bmcHelixExporter) start(ctx context.Context, host component.Host) error {

	be.logger.Info("BMC Helix Exporter is starting ...")

	// write the whole config
	be.logger.Info("BMC Helix Exporter config", zap.Any("config", be.config))

	return nil

}
