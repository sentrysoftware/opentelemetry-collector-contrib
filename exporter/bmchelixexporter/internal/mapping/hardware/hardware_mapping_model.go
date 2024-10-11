package hardware

import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/bmchelixexporter/internal/mapping"

var (
	ScopeName = "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/bmchelixexporter"
)

// HardwareMappingModelProvider implements the MappingProvider interface for hardware metrics.
type HardwareMappingModelProvider struct {
	mappings        []mapping.EntityMapping
	stateAttributes []string
}

// Ensure HardwareMappingModelProvider implements MappingProvider.
var _ mapping.MappingModelProvider = (*HardwareMappingModelProvider)(nil)

// PredefinedMappings stores the mappings for known hardware types.
var PredefinedMappings = []mapping.EntityMapping{
	{
		EntityTypeId:         "connector",
		MetricPatterns:       mapping.CompilePatterns([]string{`.*\.connector\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{},
	},
	{
		EntityTypeId:         "agent",
		MetricPatterns:       mapping.CompilePatterns([]string{`.*\.agent\.`}),
		RequiredAttrsForId:   []string{"host.id"},
		RequiredAttrsForName: []string{"service.name"},
		MetricAttributes:     map[string]string{},
	},
	{
		EntityTypeId:         "host",
		MetricPatterns:       mapping.CompilePatterns([]string{`.*\.host\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "host"},
	},
	{
		EntityTypeId:         "battery",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.battery\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "battery"},
	},
	{
		EntityTypeId:         "blade",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.blade\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "blade"},
	},
	{
		EntityTypeId:         "cpu",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.cpu\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "cpu"},
	},
	{
		EntityTypeId:         "disk_controller",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.disk_controller\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "disk_controller"},
	},
	{
		EntityTypeId:         "enclosure",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.enclosure\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "enclosure"},
	},
	{
		EntityTypeId:         "fan",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.fan\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "fan"},
	},
	{
		EntityTypeId:         "gpu",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.gpu\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "gpu"},
	},
	{
		EntityTypeId:         "led",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.led\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "led"},
	},
	{
		EntityTypeId:         "logical_disk",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.logical_disk\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "logical_disk"},
	},
	{
		EntityTypeId:         "lun",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.lun\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "lun"},
	},
	{
		EntityTypeId:         "memory",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.memory\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "memory"},
	},
	{
		EntityTypeId:         "network",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.network\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "network"},
	},
	{
		EntityTypeId:         "other_device",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.other_device\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "other_device"},
	},
	{
		EntityTypeId:         "physical_disk",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.physical_disk\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "physical_disk"},
	},
	{
		EntityTypeId:         "power_supply",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.power_supply\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "power_supply"},
	},
	{
		EntityTypeId:         "robotics",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.robotics\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "robotics"},
	},
	{
		EntityTypeId:         "tape_drive",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.tape_drive\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "tape_drive"},
	},
	{
		EntityTypeId:         "temperature",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.temperature`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "temperature"},
	},
	{
		EntityTypeId:         "vm",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.vm\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "vm"},
	},
	{
		EntityTypeId:         "voltage",
		MetricPatterns:       mapping.CompilePatterns([]string{`hw\.voltage`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "voltage"},
	},
}

// KnownHwMetricAttributeKeys are the attribute keys that are known to be part of the metric name
// https://opentelemetry.io/docs/specs/semconv/system/hardware-metrics/
var KnownHwMetricAttributeKeys = []string{"state", "direction", "hw.error.type", "limit_type", "task"}

// NewHardwareMappingModelProvider creates a new HardwareMappingModelProvider
func NewHardwareMappingModelProvider() *HardwareMappingModelProvider {
	return &HardwareMappingModelProvider{
		mappings:        PredefinedMappings,
		stateAttributes: KnownHwMetricAttributeKeys,
	}
}

// GetMappingModel returns the hardware mapping model
func (hmm *HardwareMappingModelProvider) GetMappingModel() mapping.MappingModel {
	return mapping.MappingModel{
		EntityMappings:  hmm.mappings,
		StateAttributes: hmm.stateAttributes,
	}
}