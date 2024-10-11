package mapping

import "regexp"

type EntityMapping struct {
	EntityTypeId         string            // The type of the entity (e.g., "battery", "fan", etc.)
	MetricPatterns       []*regexp.Regexp  // Precompiled regex patterns to identify this entity type
	RequiredAttrsForId   []string          // The list of attributes to form the entityId (e.g., ["id", "instance"])
	RequiredAttrsForName []string          // The list of attributes to form the entityName (e.g., ["name"])
	MetricAttributes     map[string]string // Specific metric attributes that identify this entity type (e.g., {"hw.type": "battery"})
}

type MappingModel struct {
	EntityMappings []EntityMapping // List of entity mappings
	StateAttributes []string       // List of attributes defining the state of the entity
}

// MappingModelProvider is an interface that defines the contract for providing mapping models
type MappingModelProvider interface {
    GetMappingModel() MappingModel
}

// CompilePatterns is a helper to compile regex patterns.
func CompilePatterns(patterns []string) []*regexp.Regexp {
    var compiled []*regexp.Regexp
    for _, pattern := range patterns {
        compiled = append(compiled, regexp.MustCompile(pattern))
    }
    return compiled
}

// MappingModelRegistry is a registry for mapping models
type MappingModelRegistry struct {
    MappingModels map[string]MappingModel
}

// NewMappingModelRegistry creates a new mapping model registry
func NewMappingModelRegistry() *MappingModelRegistry {
    return &MappingModelRegistry{
        MappingModels: make(map[string]MappingModel),
    }
}

// RegisterMappingModel registers a mapping model by name
func (r *MappingModelRegistry) RegisterMappingModel(name string, mappingModel MappingModel) {
    r.MappingModels[name] = mappingModel
}

// GetMappingModels returns the mapping models in the registry
func (r *MappingModelRegistry) GetMappingModels() (map[string]MappingModel) {
    return r.MappingModels;
}