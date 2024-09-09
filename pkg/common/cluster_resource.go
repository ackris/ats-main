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

package common

import "fmt"

// ClusterResource encapsulates metadata for an Atomstate cluster.
type ClusterResource struct {
	clusterID string
}

// NewClusterResource creates a new ClusterResource with the given cluster ID.
// Note that cluster ID may be an empty string if the metadata request was sent
// to a broker without support for cluster IDs.
//
// Example usage:
//
//	cr := NewClusterResource("my-cluster-id")
//	fmt.Println(cr.ClusterID()) // Output: my-cluster-id
func NewClusterResource(clusterID string) *ClusterResource {
	return &ClusterResource{
		clusterID: clusterID,
	}
}

// ClusterID returns the cluster ID.
// Note that it may be an empty string if the metadata request was sent to a
// broker without support for cluster IDs.
//
// Example usage:
//
//	cr := NewClusterResource("my-cluster-id")
//	id := cr.ClusterID()
//	fmt.Println(id) // Output: my-cluster-id
func (cr *ClusterResource) ClusterID() string {
	return cr.clusterID
}

// String returns a string representation of the ClusterResource.
//
// Example usage:
//
//	cr := NewClusterResource("my-cluster-id")
//	fmt.Println(cr.String()) // Output: ClusterResource(ClusterID=my-cluster-id)
func (cr *ClusterResource) String() string {
	return fmt.Sprintf("ClusterResource(ClusterID=%s)", cr.ClusterID())
}

// Equals checks if two ClusterResource instances are equal.
// It returns true if the other instance is not nil and has the same cluster ID.
//
// Example usage:
//
//	cr1 := NewClusterResource("my-cluster-id")
//	cr2 := NewClusterResource("my-cluster-id")
//	cr3 := NewClusterResource("another-cluster-id")
//	fmt.Println(cr1.Equals(cr2)) // Output: true
//	fmt.Println(cr1.Equals(cr3)) // Output: false
func (cr *ClusterResource) Equals(other *ClusterResource) bool {
	if other == nil {
		return false
	}
	return cr.clusterID == other.clusterID
}

// Hash returns a simple hash value for the ClusterResource.
// This can be useful for using ClusterResource as a key in maps.
//
// Example usage:
//
//	cr := NewClusterResource("my-cluster-id")
//	hashValue := cr.Hash()
//	fmt.Println(hashValue) // Output: (some integer value)
func (cr *ClusterResource) Hash() int {
	hash := 0
	for _, char := range cr.clusterID {
		hash += int(char)
	}
	return hash
}
