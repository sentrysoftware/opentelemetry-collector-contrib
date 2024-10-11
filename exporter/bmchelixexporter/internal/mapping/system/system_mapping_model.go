package system

import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/bmchelixexporter/internal/mapping"

var (
	ScopeName = "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/bmchelixexporter"
)

// SystemMappingModelProvider implements the MappingProvider interface for system metrics.
type SystemMappingModelProvider struct {
	mappings        []mapping.EntityMapping
	stateAttributes []string
}

// Ensure SystemMappingModelProvider implements MappingProvider.
var _ mapping.MappingModelProvider = (*SystemMappingModelProvider)(nil)

// PredefinedMappings stores the mappings for known system metrics.
var PredefinedMappings = []mapping.EntityMapping{
	{
		EntityTypeId:         "system_cpu",
		MetricPatterns:       mapping.CompilePatterns([]string{`system\.cpu\.`}),
		RequiredAttrsForId:   []string{"system.cpu.logical_number"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{},
	},
	{
		EntityTypeId:         "system_memory",
		MetricPatterns:       mapping.CompilePatterns([]string{`system\.memory\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"id"},
		MetricAttributes:     map[string]string{},
	},
	{
		EntityTypeId:         "system_paging",
		MetricPatterns:       mapping.CompilePatterns([]string{`system\.paging\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"id"},
		MetricAttributes:     map[string]string{},
	},
	{
		EntityTypeId:         "system_disk",
		MetricPatterns:       mapping.CompilePatterns([]string{`system\.disk\.`}),
		RequiredAttrsForId:   []string{"id", "system.device"},
		RequiredAttrsForName: []string{"id", "system.device"},
		MetricAttributes:     map[string]string{},
	},
	{
		EntityTypeId:         "system_filesystem",
		MetricPatterns:       mapping.CompilePatterns([]string{`system\.filesystem\.`}),
		RequiredAttrsForId:   []string{"id", "system.device"},
		RequiredAttrsForName: []string{"id", "system.filesystem.mountpoint", "system.filesystem.volumeName"},
		MetricAttributes:     map[string]string{},
	},
	{
		EntityTypeId:         "system_network",
		MetricPatterns:       mapping.CompilePatterns([]string{`system\.network\.`}),
		RequiredAttrsForId:   []string{"id", "system.device"},
		RequiredAttrsForName: []string{"id", "system.device"},
		MetricAttributes:     map[string]string{},
	},
	{
		EntityTypeId:         "system_process",
		MetricPatterns:       mapping.CompilePatterns([]string{`process\.`}),
		RequiredAttrsForId:   []string{"process.id"},
		RequiredAttrsForName: []string{"process.name"},
		MetricAttributes:     map[string]string{},
	},
	{
		EntityTypeId:         "system_service",
		MetricPatterns:       mapping.CompilePatterns([]string{`system\.service\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"id", "system.service.name"},
		MetricAttributes:     map[string]string{},
	},
}

// KnownSystemMetricAttributeKeys are the attribute keys that are known to be part of the metric name
var KnownSystemMetricAttributeKeys = []string{
	"cpu.mode", "system.cpu.state", "system.memory.state", "system.paging.state", 
	"disk.io.direction", "network.io.direction", "system.filesystem.state", "system.process.status",
	 "system.filesystem.state", "state",
}

// NewSystemMappingModelProvider creates a new SystemMappingModelProvider
func NewSystemMappingModelProvider() *SystemMappingModelProvider {
	return &SystemMappingModelProvider{
		mappings:        PredefinedMappings,
		stateAttributes: KnownSystemMetricAttributeKeys,
	}
}

// GetMappingModel returns the system mapping model
func (smm *SystemMappingModelProvider) GetMappingModel() mapping.MappingModel {
	return mapping.MappingModel{
		EntityMappings:  smm.mappings,
		StateAttributes: smm.stateAttributes,
	}
}
