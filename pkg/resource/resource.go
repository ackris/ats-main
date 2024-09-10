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

import "fmt"

// Resource represents a cluster resource with a tuple of (type, name).
type Resource struct {
	resourceType ResourceType
	name         string
}

// CLUSTER_NAME is the name of the CLUSTER resource.
const CLUSTER_NAME = "atomstate-cluster"

// ClUSTER_RESOURCE is a resource representing the whole cluster.
var ClUSTER_RESOURCE = NewResource(CLUSTER, CLUSTER_NAME)

// NewResource creates an instance of Resource with the provided parameters.
//
// Parameters:
//   - resourceType: The type of the resource. Must not be UNRECOGNIZED.
//   - name: The name of the resource. Must not be empty.
//
// Returns:
//   - A pointer to a new Resource instance.
//
// Example usage:
//
//	resource := NewResource(TOPIC, "my-topic")
//	fmt.Println(resource.String()) // Output: (resourceType=TOPIC, name=my-topic)
func NewResource(resourceType ResourceType, name string) *Resource {
	if resourceType == UNRECOGNIZED {
		panic("resourceType cannot be UNRECOGNIZED")
	}
	if name == "" {
		panic("name cannot be empty")
	}
	return &Resource{resourceType: resourceType, name: name}
}

// ResourceType returns the resource type of the Resource instance.
//
// Returns:
//   - The ResourceType of the resource.
//
// Example usage:
//
//	resource := NewResource(GROUP, "my-group")
//	fmt.Println(resource.ResourceType()) // Output: GROUP
func (r *Resource) ResourceType() ResourceType {
	return r.resourceType
}

// Name returns the name of the Resource instance.
//
// Returns:
//   - The name of the resource.
//
// Example usage:
//
//	resource := NewResource(TOPIC, "my-topic")
//	fmt.Println(resource.Name()) // Output: my-topic
func (r *Resource) Name() string {
	return r.name
}

// IsUnknown returns true if this Resource has any UNKNOWN components.
//
// Returns:
//   - true if the resource type is unrecognized; otherwise, false.
//
// Example usage:
//
//	resource := NewResource(UNRECOGNIZED, "unknown-resource")
//	fmt.Println(resource.IsUnknown()) // Output: true
func (r *Resource) IsUnknown() bool {
	return r.resourceType.IsUnrecognized()
}

// String returns a string representation of the Resource instance.
//
// Returns:
//   - A formatted string representing the resource.
//
// Example usage:
//
//	resource := NewResource(TOPIC, "my-topic")
//	fmt.Println(resource.String()) // Output: (resourceType=TOPIC, name=my-topic)
func (r *Resource) String() string {
	return fmt.Sprintf("(resourceType=%s, name=%s)", r.resourceType.String(), r.name)
}

// Equals checks if two Resource instances are equal.
//
// Parameters:
//   - other: A pointer to another Resource instance to compare with.
//
// Returns:
//   - true if both resources are equal; otherwise, false.
//
// Example usage:
//
//	resource1 := NewResource(TOPIC, "my-topic")
//	resource2 := NewResource(TOPIC, "my-topic")
//	fmt.Println(resource1.Equals(resource2)) // Output: true
func (r *Resource) Equals(other *Resource) bool {
	if other == nil {
		return false
	}
	return r.resourceType == other.resourceType && r.name == other.name
}
