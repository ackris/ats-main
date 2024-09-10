// Copyright 2024 Atomstate Technologies Private Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"errors"
	"fmt"
	"strings"
)

// ResourcePatternFilter represents a filter that can match ResourcePattern.
// This filter can be used to determine if a given ResourcePattern matches specific criteria
// based on resource type, name, and pattern type.
type ResourcePatternFilter struct {
	resourceType ResourceType // The type of resource this filter applies to.
	name         string       // The name of the resource, can be empty or a specific name.
	patternType  PatternType  // The type of pattern this filter uses.
}

// NewResourcePatternFilter creates a new ResourcePatternFilter.
//
// Parameters:
// - resourceType: The type of resource to filter. Must not be UNRECOGNIZED.
// - name: The name of the resource to filter. Can be empty or a specific name.
// - patternType: The type of pattern to filter. Must not be UNKNOWN.
//
// Returns:
// - A pointer to a new ResourcePatternFilter instance.
// - An error if the resourceType is UNRECOGNIZED or patternType is UNKNOWN.
//
// Example usage:
// filter, err := NewResourcePatternFilter(TOPIC, "payments.", PREFIXED)
//
//	if err != nil {
//	    // Handle error
//	}
func NewResourcePatternFilter(
	resourceType ResourceType,
	name string,
	patternType PatternType) (*ResourcePatternFilter, error) {
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
// This is a convenience function that returns a filter with ALL_RESOURCES type,
// an empty name, and ANY pattern type.
//
// Returns:
// - A pointer to a new ResourcePatternFilter instance that matches any resource pattern.
//
// Example usage:
// filter := AnyResourcePatternFilter()
func AnyResourcePatternFilter() *ResourcePatternFilter {
	filter, _ := NewResourcePatternFilter(ALL_RESOURCES, "", ANY)
	return filter
}

// IsUnknown checks if the filter has any UNKNOWN components.
//
// Returns:
// - true if the resourceType is UNRECOGNIZED or the patternType is UNKNOWN, false otherwise.
//
// Example usage:
//
//	if filter.IsUnknown() {
//	    // Handle unknown filter
//	}
func (f *ResourcePatternFilter) IsUnknown() bool {
	return f.resourceType == UNRECOGNIZED || f.patternType == UNKNOWN
}

// ResourceType returns the specific resource type this pattern matches.
//
// Returns:
// - The resource type of the filter.
//
// Example usage:
// resType := filter.ResourceType()
func (f *ResourcePatternFilter) ResourceType() ResourceType {
	return f.resourceType
}

// Name returns the resource name.
//
// Returns:
// - The name of the resource in the filter.
//
// Example usage:
// name := filter.Name()
func (f *ResourcePatternFilter) Name() string {
	return f.name
}

// PatternType returns the resource pattern type.
//
// Returns:
// - The pattern type of the filter.
//
// Example usage:
// patternType := filter.PatternType()
func (f *ResourcePatternFilter) PatternType() PatternType {
	return f.patternType
}

// Matches checks if the filter matches the given ResourcePattern.
//
// Parameters:
// - pattern: A pointer to a ResourcePattern to match against this filter.
//
// Returns:
// - true if the filter matches the pattern, false otherwise.
//
// Example usage:
// pattern, _ := NewResourcePattern(TOPIC, "payments.received", LITERAL)
//
//	if filter.Matches(pattern) {
//	    // The filter matches the pattern
//	}
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
// This method is used internally by the Matches method.
//
// Parameters:
// - pattern: A pointer to a ResourcePattern to match against this filter.
//
// Returns:
// - true if the name matches according to the filter's pattern type, false otherwise.
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
//
// Returns:
// - true if there are no ANY or UNKNOWN fields in the filter, false otherwise.
//
// Example usage:
//
//	if filter.MatchesAtMostOne() {
//	    // The filter can only match one pattern
//	}
func (f *ResourcePatternFilter) MatchesAtMostOne() bool {
	return f.findIndefiniteField() == ""
}

// findIndefiniteField returns a string describing any indefinite field, or an empty string if none.
// This method is used internally to determine if the filter has any indefinite components.
//
// Returns:
// - A string describing any indefinite field or an empty string if none.
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
