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
