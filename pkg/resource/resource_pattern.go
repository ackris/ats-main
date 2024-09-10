package resource

import (
	"errors"
	"fmt"
)

// ResourcePattern represents a pattern used by ACLs to match zero or more resources.
type ResourcePattern struct {
	resourceType ResourceType
	name         string
	patternType  PatternType
}

const (
	// WILDCARD_RESOURCE is a special literal resource name that corresponds to 'all resources of a certain type'.
	WILDCARD_RESOURCE = "*"
)

// NewResourcePattern creates a new ResourcePattern using the supplied parameters.
func NewResourcePattern(resourceType ResourceType, name string, patternType PatternType) (*ResourcePattern, error) {
	if resourceType == UNRECOGNIZED {
		return nil, errors.New("resourceType cannot be UNRECOGNIZED")
	}
	if patternType == UNKNOWN {
		return nil, errors.New("patternType cannot be UNKNOWN")
	}
	return &ResourcePattern{
		resourceType: resourceType,
		name:         name,
		patternType:  patternType,
	}, nil
}

// ResourceType returns the specific resource type this pattern matches.
func (rp *ResourcePattern) ResourceType() ResourceType {
	return rp.resourceType
}

// Name returns the resource name.
func (rp *ResourcePattern) Name() string {
	return rp.name
}

// PatternType returns the resource pattern type.
func (rp *ResourcePattern) PatternType() PatternType {
	return rp.patternType
}

// ToFilter returns a filter which matches only this pattern.
func (rp *ResourcePattern) ToFilter() (*ResourcePatternFilter, error) {
	return NewResourcePatternFilter(rp.resourceType, rp.name, rp.patternType)
}

// IsUnknown returns true if this ResourcePattern has any UNKNOWN components.
func (rp *ResourcePattern) IsUnknown() bool {
	return rp.resourceType == UNRECOGNIZED || rp.patternType == UNKNOWN
}

// String returns a string representation of the ResourcePattern.
func (p *ResourcePattern) String() string {
	return fmt.Sprintf("ResourcePattern{resourceType=%s, name=%q, patternType=%s}", p.resourceType.String(), p.name, p.patternType.String())
}

// Equals checks if two ResourcePatterns are equal.
func (rp *ResourcePattern) Equals(other *ResourcePattern) bool {
	return other != nil && rp.resourceType == other.resourceType && rp.name == other.name && rp.patternType == other.patternType
}
