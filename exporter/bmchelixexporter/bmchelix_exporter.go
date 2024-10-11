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

	var mappingModels map[string]mapping.MappingModel
	if be.config.MappingModel != nil {
		mappingModels = buildMappingModelsUsingConfig(be.config.MappingModel)
	} else {
		mappingModels = newDefaultMappingModelRegistry().GetMappingModels()
	}

	// Initialize and store the BmcHelixMetricsProducer
	be.metricsProducer = &BmcHelixMetricsProducer{
		osHostname: osHostname,
		logger:     be.logger,
		mappingResolver: &mapping.MappingResolver{
			Logger:        be.logger,
			MappingModels: mappingModels,
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

// buildMappingModelsUsingConfig builds the mapping models using the provided configuration
func buildMappingModelsUsingConfig(mappingModelConfig *MappingModelConfig) (map[string]mapping.MappingModel) {
	mappingModels := map[string]mapping.MappingModel{
		"fromConfig": {
			EntityMappings: buildEntityMappingsFromConfig(mappingModelConfig.EntityMappings),
			StateAttributes: mappingModelConfig.StateAttributes,
		},
	}

	// Early return if ReplaceDefault is true
	if mappingModelConfig.Force {
		return mappingModels
	}

	// Add the default mapping models
	for k, v := range newDefaultMappingModelRegistry().GetMappingModels() {
		mappingModels[k] = v
	}
	return mappingModels
}

// Helper to build entity mappings from the configuration
func buildEntityMappingsFromConfig(entityMappingsConfig []EntityMappingConfig) []mapping.EntityMapping {
	entityMappings := make([]mapping.EntityMapping, 0, len(entityMappingsConfig))
	for _, config := range entityMappingsConfig {
		entityMappings = append(entityMappings, mapping.EntityMapping{
			EntityTypeId:         config.EntityTypeId,
			MetricPatterns:       mapping.CompilePatterns(config.MetricPatterns),
			RequiredAttrsForId:   config.RequiredAttrsForId,
			RequiredAttrsForName: config.RequiredAttrsForName,
			MetricAttributes:     config.MetricAttributes,
		})
	}
	return entityMappings
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
