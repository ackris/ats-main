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
	"hash/fnv"
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
//
// Parameters:
//   - resourceType: The type of resource to match (e.g., TOPIC, GROUP).
//   - name: The name of the resource. Can be a specific name or the wildcard "*".
//
// Returns:
//   - A pointer to a ResourcePattern and an error if the input is invalid.
//
// Example:
//
//	rp, err := NewResourcePattern(TOPIC, "my-topic", LITERAL)
//	if err != nil {
//	    // Handle error
//	}
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
//
// Returns:
//   - The ResourceType of this ResourcePattern.
//
// Example:
//
//	resourceType := rp.ResourceType()
func (rp *ResourcePattern) ResourceType() ResourceType {
	return rp.resourceType
}

// Name returns the resource name.
//
// Returns:
//   - The name of the resource as a string.
//
// Example:
//
//	resourceName := rp.Name()
func (rp *ResourcePattern) Name() string {
	return rp.name
}

// PatternType returns the resource pattern type.
//
// Returns:
//   - The PatternType of this ResourcePattern.
//
// Example:
//
//	patternType := rp.PatternType()
func (rp *ResourcePattern) PatternType() PatternType {
	return rp.patternType
}

// ToFilter returns a filter which matches only this pattern.
//
// Returns:
//   - A pointer to a ResourcePatternFilter and an error if the creation fails.
//
// Example:
//
//	filter, err := rp.ToFilter()
//	if err != nil {
//	    // Handle error
//	}
func (rp *ResourcePattern) ToFilter() (*ResourcePatternFilter, error) {
	return NewResourcePatternFilter(rp.resourceType, rp.name, rp.patternType)
}

// IsUnknown returns true if this ResourcePattern has any UNKNOWN components.
//
// Returns:
//   - true if the resource type or pattern type is unknown; otherwise false.
//
// Example:
//
//	if rp.IsUnknown() {
//	    // Handle unknown pattern
//	}
func (rp *ResourcePattern) IsUnknown() bool {
	return rp.resourceType == UNRECOGNIZED || rp.patternType == UNKNOWN
}

// String returns a string representation of the ResourcePattern.
//
// Returns:
//   - A formatted string describing the ResourcePattern.
//
// Example:
//
//	fmt.Println(rp.String())
func (rp *ResourcePattern) String() string {
	return fmt.Sprintf("ResourcePattern{resourceType=%s, name=%q, patternType=%s}",
		rp.resourceType.String(),
		rp.name,
		rp.patternType.String())
}

// Equals checks if two ResourcePatterns are equal.
//
// Parameters:
//   - other: Another ResourcePattern to compare against.
//
// Returns:
//   - true if the two ResourcePatterns are equal; otherwise false.
//
// Example:
//
//	if rp1.Equals(rp2) {
//	    // Patterns are equal
//	}
func (rp *ResourcePattern) Equals(other *ResourcePattern) bool {
	return other != nil && rp.resourceType == other.resourceType &&
		rp.name == other.name && rp.patternType == other.patternType
}

// HashCode computes the hash code for the Resource struct.
func (r *ResourcePattern) HashCode() uint64 {
	if r == nil {
		return 0 // Return a default hash value if the ResourcePattern is nil
	}

	h := fnv.New64a() // Create a new FNV-1a hash

	// Combine fields into a single byte slice to minimize writes
	data := []byte(r.ResourceType().String() + r.Name() + r.PatternType().String())
	h.Write(data) // Write the combined data to the hash

	return h.Sum64() // Return the computed hash code
}
