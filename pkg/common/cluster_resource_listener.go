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

// ClusterResourceListener is an interface for receiving notifications about changes in cluster metadata.
// It is useful for interceptors, metric reporters, serializers, and deserializers needing cluster metadata.
//
// Each metadata response triggers one OnUpdate call. Note that the cluster ID may be empty if the Atomstate broker version is below 0.10.1.0.
// If empty, the cluster ID will remain so unless there are brokers with different versions due to cluster upgrades.
//
// Example:
//
// type MyClusterListener struct{}
//
//	func (m *MyClusterListener) OnUpdate(clusterResource ClusterResource) {
//	    fmt.Printf("Cluster updated: %s\n", clusterResource.ClusterID)
//	}
//
//	func main() {
//	    listener := &MyClusterListener{}
//	    listener.OnUpdate(ClusterResource{ClusterID: "my-cluster-id"})
//	}
type ClusterResourceListener interface {
	// OnUpdate is a callback method that a user can implement to get updates for ClusterResource.
	//
	// Parameters:
	// - clusterResource: The cluster metadata that is updated.
	OnUpdate(clusterResource ClusterResource)
}
