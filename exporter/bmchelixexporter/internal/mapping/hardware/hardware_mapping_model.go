package hardware

import "regexp"

var (
	ScopeName = "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/bmchelixexporter"
)

// EntityMapping represents the rules to map entityTypeId, entityId, and entityName based on metric names or attributes.
type EntityMapping struct {
	EntityTypeId         string            // The type of the entity (e.g., "battery", "fan", etc.)
	MetricPatterns       []*regexp.Regexp  // Precompiled regex patterns to identify this entity type
	RequiredAttrsForId   []string          // The list of attributes to form the entityId (e.g., ["id", "instance"])
	RequiredAttrsForName []string          // The list of attributes to form the entityName (e.g., ["name"])
	MetricAttributes     map[string]string // Specific metric attributes that identify this entity type (e.g., {"hw.type": "battery"})
}

// PredefinedMappings stores the mappings for known hardware types.
var PredefinedMappings = []EntityMapping{
	{
		EntityTypeId:         "connector",
		MetricPatterns:       compilePatterns([]string{`.*\.connector\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{},
	},
	{
		EntityTypeId:         "agent",
		MetricPatterns:       compilePatterns([]string{`.*\.agent\.`}),
		RequiredAttrsForId:   []string{"host.id"},
		RequiredAttrsForName: []string{"service.name"},
		MetricAttributes:     map[string]string{},
	},
	{
		EntityTypeId:         "host",
		MetricPatterns:       compilePatterns([]string{`.*\.host\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "host"},
	},
	{
		EntityTypeId:         "battery",
		MetricPatterns:       compilePatterns([]string{`hw\.battery\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "battery"},
	},
	{
		EntityTypeId:         "blade",
		MetricPatterns:       compilePatterns([]string{`hw\.blade\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "blade"},
	},
	{
		EntityTypeId:         "cpu",
		MetricPatterns:       compilePatterns([]string{`hw\.cpu\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "cpu"},
	},
	{
		EntityTypeId:         "disk_controller",
		MetricPatterns:       compilePatterns([]string{`hw\.disk_controller\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "disk_controller"},
	},
	{
		EntityTypeId:         "enclosure",
		MetricPatterns:       compilePatterns([]string{`hw\.enclosure\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "enclosure"},
	},
	{
		EntityTypeId:         "fan",
		MetricPatterns:       compilePatterns([]string{`hw\.fan\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "fan"},
	},
	{
		EntityTypeId:         "gpu",
		MetricPatterns:       compilePatterns([]string{`hw\.gpu\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "gpu"},
	},
	{
		EntityTypeId:         "led",
		MetricPatterns:       compilePatterns([]string{`hw\.led\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "led"},
	},
	{
		EntityTypeId:         "logical_disk",
		MetricPatterns:       compilePatterns([]string{`hw\.logical_disk\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "logical_disk"},
	},
	{
		EntityTypeId:         "lun",
		MetricPatterns:       compilePatterns([]string{`hw\.lun\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "lun"},
	},
	{
		EntityTypeId:         "memory",
		MetricPatterns:       compilePatterns([]string{`hw\.memory\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "memory"},
	},
	{
		EntityTypeId:         "network",
		MetricPatterns:       compilePatterns([]string{`hw\.network\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "network"},
	},
	{
		EntityTypeId:         "other_device",
		MetricPatterns:       compilePatterns([]string{`hw\.other_device\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "other_device"},
	},
	{
		EntityTypeId:         "physical_disk",
		MetricPatterns:       compilePatterns([]string{`hw\.physical_disk\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "physical_disk"},
	},
	{
		EntityTypeId:         "power_supply",
		MetricPatterns:       compilePatterns([]string{`hw\.power_supply\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "power_supply"},
	},
	{
		EntityTypeId:         "robotics",
		MetricPatterns:       compilePatterns([]string{`hw\.robotics\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "robotics"},
	},
	{
		EntityTypeId:         "tape_drive",
		MetricPatterns:       compilePatterns([]string{`hw\.tape_drive\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "tape_drive"},
	},
	{
		EntityTypeId:         "temperature",
		MetricPatterns:       compilePatterns([]string{`hw\.temperature\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "temperature"},
	},
	{
		EntityTypeId:         "vm",
		MetricPatterns:       compilePatterns([]string{`hw\.vm\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "vm"},
	},
	{
		EntityTypeId:         "voltage",
		MetricPatterns:       compilePatterns([]string{`hw\.voltage\.`}),
		RequiredAttrsForId:   []string{"id"},
		RequiredAttrsForName: []string{"name"},
		MetricAttributes:     map[string]string{"hw.type": "voltage"},
	},
}

// compilePatterns is a helper to compile regex patterns
func compilePatterns(patterns []string) []*regexp.Regexp {
	var compiled []*regexp.Regexp
	for _, pattern := range patterns {
		r, err := regexp.Compile(pattern)
		if err == nil {
			compiled = append(compiled, r)
		}
	}
	return compiled
}

// knownHwMetricAttributeMap are the attributes that are known to be part of the metric name
// https://opentelemetry.io/docs/specs/semconv/system/hardware-metrics/
var KnownHwMetricAttributeMap = map[string]bool{"state": true, "direction": true, "hw.error.type": true, "limit_type": true, "task": true}

// knownHwMetricAttributeKeys are the keys of the knownHwMetricAttributeMap
var KnownHwMetricAttributeKeys = func() []string {
	keys := make([]string, 0, len(KnownHwMetricAttributeMap))
	for k := range KnownHwMetricAttributeMap {
		keys = append(keys, k)
	}
	return keys
}()