package bmchelixexporter

import (
	"context"
	"errors"
	"os"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/bmchelixexporter/internal/mapping"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/bmchelixexporter/internal/mapping/hardware"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/bmchelixexporter/internal/mapping/system"
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
		mappingResolver: &mapping.MappingResolver{
			Logger:        be.logger,
			MappingModels: newDefaultMappingModelRegistry().GetMappingModels(),
		},
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

// newDefaultMappingModelRegistry creates a new mapping model registry with the default mapping models
func newDefaultMappingModelRegistry() *mapping.MappingModelRegistry {
	// Create the mapping model registry
	mappingModelRegistry := mapping.NewMappingModelRegistry()

	// Get the hardware mapping model
	hardwareMappingModel := hardware.NewHardwareMappingModelProvider().GetMappingModel()

	// Get the system mapping model
	systemMappingModel := system.NewSystemMappingModelProvider().GetMappingModel()

	// Register the hardware mapping model
	mappingModelRegistry.RegisterMappingModel("hardware", hardwareMappingModel)

	// Register the system mapping model
	mappingModelRegistry.RegisterMappingModel("system", systemMappingModel)

	return mappingModelRegistry
}
