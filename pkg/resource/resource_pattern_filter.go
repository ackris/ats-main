package resource

import (
	"errors"
	"fmt"
	"strings"
)

// ResourcePatternFilter represents a filter that can match ResourcePattern.
type ResourcePatternFilter struct {
	resourceType ResourceType
	name         string
	patternType  PatternType
}

// NewResourcePatternFilter creates a new ResourcePatternFilter.
func NewResourcePatternFilter(resourceType ResourceType, name string, patternType PatternType) (*ResourcePatternFilter, error) {
	if resourceType == UNRECOGNIZED {
		return nil, errors.New("resourceType cannot be UNRECOGNIZED")
	}
	if patternType == UNKNOWN {
		return nil, errors.New("patternType cannot be UNKNOWN")
	}
	return &ResourcePatternFilter{
		resourceType: resourceType,
		name:         name,
		patternType:  patternType,
	}, nil
}

// AnyResourcePatternFilter creates a filter that matches any resource pattern.
func AnyResourcePatternFilter() *ResourcePatternFilter {
	filter, _ := NewResourcePatternFilter(ALL_RESOURCES, "", ANY)
	return filter
}

// IsUnknown checks if the filter has any UNKNOWN components.
func (f *ResourcePatternFilter) IsUnknown() bool {
	return f.resourceType == UNRECOGNIZED || f.patternType == UNKNOWN
}

// ResourceType returns the specific resource type this pattern matches.
func (f *ResourcePatternFilter) ResourceType() ResourceType {
	return f.resourceType
}

// Name returns the resource name.
func (f *ResourcePatternFilter) Name() string {
	return f.name
}

// PatternType returns the resource pattern type.
func (f *ResourcePatternFilter) PatternType() PatternType {
	return f.patternType
}

// Matches checks if the filter matches the given ResourcePattern.
func (f *ResourcePatternFilter) Matches(pattern *ResourcePattern) bool {
	if pattern == nil {
		return false
	}

	if f.resourceType != ALL_RESOURCES && f.resourceType != pattern.resourceType {
		return false
	}

	if f.patternType != ANY && f.patternType != MATCH && f.patternType != pattern.patternType {
		return false
	}

	if f.name == "" {
		return true
	}

	return f.nameMatches(pattern)
}

// nameMatches checks if the name matches based on the pattern type.
func (f *ResourcePatternFilter) nameMatches(pattern *ResourcePattern) bool {
	switch {
	case f.patternType == ANY:
		return true
	case f.patternType == pattern.patternType:
		return f.name == pattern.name || f.name == WILDCARD_RESOURCE
	case f.patternType == PREFIXED:
		// Match if the pattern name starts with the filter's name
		return strings.HasPrefix(pattern.name, f.name)
	case pattern.patternType == LITERAL:
		return f.name == pattern.name || pattern.name == WILDCARD_RESOURCE
	default:
		return false // Return false for unsupported pattern types
	}
}

// MatchesAtMostOne checks if the filter could only match one pattern.
func (f *ResourcePatternFilter) MatchesAtMostOne() bool {
	return f.findIndefiniteField() == ""
}

// findIndefiniteField returns a string describing any indefinite field, or an empty string if none.
func (f *ResourcePatternFilter) findIndefiniteField() string {
	switch {
	case f.resourceType == ALL_RESOURCES:
		return "Resource type is ALL_RESOURCES."
	case f.resourceType == UNRECOGNIZED:
		return "Resource type is UNRECOGNIZED."
	case f.name == "":
		return "Resource name is empty."
	case f.patternType == MATCH:
		return "Resource pattern type is MATCH."
	case f.patternType == UNKNOWN:
		return "Resource pattern type is UNKNOWN."
	default:
		return ""
	}
}

// String returns a string representation of the ResourcePatternFilter.
func (f *ResourcePatternFilter) String() string {
	return fmt.Sprintf("ResourcePatternFilter{resourceType=%s, name=%q, patternType=%s}", f.resourceType.String(), f.name, f.patternType.String())
}
