package bmchelixexporter

import (
	"context"
	"errors"
	"os"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type BmcHelixExporter struct {
	config            *Config
	logger            *zap.Logger
	version           string
	telemetrySettings component.TelemetrySettings
	metricsProducer   *BmcHelixMetricsProducer
	metricsConsumer   *BmcHelixMetricsConsumer
}

func newBmcHelixExporter(config *Config, createSettings exporter.Settings) (*BmcHelixExporter, error) {
	if config == nil {
		return nil, errors.New("nil config")
	}

	return &BmcHelixExporter{
		config:            config,
		version:           createSettings.BuildInfo.Version,
		logger:            createSettings.Logger,
		telemetrySettings: createSettings.TelemetrySettings,
	}, nil
}

// pushMetrics is invoked by the OpenTelemetry Collector to push metrics to the BMC Helix
func (be *BmcHelixExporter) pushMetrics(ctx context.Context, md pmetric.Metrics) error {

	// Push the metrics to the BMC Helix
	be.logger.Info("Building BMC Helix payload")
	helixMetrics, err := be.metricsProducer.ProduceHelixPayload(md)
	if err != nil {
		be.logger.Error("Failed to build BMC Helix payload", zap.Error(err))
		return err
	}

	be.logger.Info("Sending BMC Helix payload")
	err = be.metricsConsumer.SendHelixPayload(ctx, helixMetrics)
	if err != nil {
		be.logger.Error("Failed to send BMC Helix payload", zap.Error(err))
		return err
	}

	return nil
}

// start is invoked during service start
func (be *BmcHelixExporter) start(ctx context.Context, host component.Host) error {

	be.logger.Info("Starting BMC Helix Exporter")

	// Get the hostname reported by the kernel
	osHostname, err := os.Hostname()
	if err != nil {
		be.logger.Warn("Failed to get OS hostname", zap.Error(err))
		return err
	}

	// Initialize and store the BmcHelixMetricsProducer
	be.metricsProducer = &BmcHelixMetricsProducer{
		osHostname: osHostname,
		logger:     be.logger,
	}

	// Initialize and store the BmcHelixMetricsConsumer
	be.metricsConsumer = &BmcHelixMetricsConsumer{
		apiKey:   be.config.ApiKey,
		endpoint: be.config.Endpoint,
		logger:   be.logger,
		timeout:  be.config.Timeout,
	}

	be.logger.Info("Initialized BMC Helix Exporter")
	return nil

}
