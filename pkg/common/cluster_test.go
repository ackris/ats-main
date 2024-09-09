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

import (
	"testing"
)

// TestNewCluster verifies that the cluster is created correctly.
func TestNewCluster(t *testing.T) {
	cluster := NewCluster("test-cluster", false, nil, nil, nil, nil, nil, nil, nil)

	if cluster.ClusterID() != "test-cluster" {
		t.Errorf("Expected cluster ID to be 'test-cluster', got '%s'", cluster.ClusterID())
	}
	if len(cluster.Nodes()) != 0 {
		t.Errorf("Expected no nodes, got %d", len(cluster.Nodes()))
	}
}

// TestAddNode verifies that nodes are added correctly.
func TestAddNode(t *testing.T) {
	cluster := EmptyCluster()
	node := &Node{ID: 1, Host: "localhost:9092"}

	err := cluster.AddNode(node)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(cluster.Nodes()) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(cluster.Nodes()))
	}

	// Test adding the same node again
	err = cluster.AddNode(node)
	if err == nil || err.Error() != "node with the same ID already exists" {
		t.Fatalf("Expected error for duplicate node, got %v", err)
	}
}

// TestAddPartition verifies that partitions are added correctly.
func TestAddPartition(t *testing.T) {
	cluster := EmptyCluster()
	node := &Node{ID: 1, Host: "localhost:9092"}
	cluster.AddNode(node)

	partition := &PartitionInfo{
		Topic:     "test-topic",
		Partition: 0,
		Leader:    node,
		Replicas:  []*Node{node}, // Ensure the node is a replica
	}
	err := cluster.AddPartition(partition)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(cluster.PartitionsForTopic("test-topic")) != 1 {
		t.Fatalf("Expected 1 partition for topic 'test-topic', got %d", len(cluster.PartitionsForTopic("test-topic")))
	}

	// Test adding the same partition again
	err = cluster.AddPartition(partition)
	if err == nil || err.Error() != "partition already exists" {
		t.Fatalf("Expected error for duplicate partition, got %v", err)
	}
}

// TestRemoveNode verifies that nodes are removed correctly.
func TestRemoveNode(t *testing.T) {
	cluster := EmptyCluster()
	node := &Node{ID: 1, Host: "localhost:9092"}
	cluster.AddNode(node)

	partition := &PartitionInfo{
		Topic:     "test-topic",
		Partition: 0,
		Leader:    node,
		Replicas:  []*Node{node}, // Ensure the node is a replica
	}
	cluster.AddPartition(partition)

	cluster.RemoveNode(node.ID)

	if len(cluster.Nodes()) != 0 {
		t.Fatalf("Expected 0 nodes after removal, got %d", len(cluster.Nodes()))
	}

	if len(cluster.PartitionsForTopic("test-topic")) != 0 {
		t.Fatalf("Expected 0 partitions for topic 'test-topic' after node removal, got %d", len(cluster.PartitionsForTopic("test-topic")))
	}
}

// TestRemovePartition verifies that partitions are removed correctly.
func TestRemovePartition(t *testing.T) {
	cluster := EmptyCluster()
	node := &Node{ID: 1, Host: "localhost:9092"}
	cluster.AddNode(node)

	partition := &PartitionInfo{
		Topic:     "test-topic",
		Partition: 0,
		Leader:    node,
		Replicas:  []*Node{node}, // Ensure the node is a replica
	}
	cluster.AddPartition(partition)

	cluster.RemovePartition("test-topic", 0)

	if len(cluster.PartitionsForTopic("test-topic")) != 0 {
		t.Fatalf("Expected 0 partitions for topic 'test-topic' after removal, got %d", len(cluster.PartitionsForTopic("test-topic")))
	}
}

// TestLeaderForPartition verifies that the correct leader is returned for a partition.
func TestLeaderForPartition(t *testing.T) {
	cluster := EmptyCluster()
	node := &Node{ID: 1, Host: "localhost:9092"}
	cluster.AddNode(node)

	partition := &PartitionInfo{
		Topic:     "test-topic",
		Partition: 0,
		Leader:    node,
		Replicas:  []*Node{node}, // Ensure the node is a replica
	}
	cluster.AddPartition(partition)

	leader := cluster.LeaderFor(NewTopicPartition("test-topic", 0))
	if leader == nil || leader.ID != node.ID {
		t.Fatalf("Expected leader to be node with ID %d, got %v", node.ID, leader)
	}
}

// TestNodeIfOnline verifies that the correct node is returned if it is online.
func TestNodeIfOnline(t *testing.T) {
	cluster := EmptyCluster()
	node := &Node{ID: 1, Host: "localhost:9092"}
	cluster.AddNode(node)

	partition := &PartitionInfo{
		Topic:     "test-topic",
		Partition: 0,
		Leader:    node,
		Replicas:  []*Node{node}, // Ensure the node is a replica
	}
	cluster.AddPartition(partition)

	onlineNode := cluster.NodeIfOnline(NewTopicPartition("test-topic", 0), node.ID)
	if onlineNode == nil || onlineNode.ID != node.ID {
		t.Fatalf("Expected online node to be node with ID %d, got %v", node.ID, onlineNode)
	}

	// Test with an offline node
	offlineNode := cluster.NodeIfOnline(NewTopicPartition("test-topic", 0), 2)
	if offlineNode != nil {
		t.Fatalf("Expected no node for ID 2, got %v", offlineNode)
	}
}

// TestAvailablePartitionsForTopic verifies the available partitions for a topic.
func TestAvailablePartitionsForTopic(t *testing.T) {
	cluster := EmptyCluster()
	node := &Node{ID: 1, Host: "localhost:9092"}
	cluster.AddNode(node)

	partition1 := &PartitionInfo{
		Topic:     "test-topic",
		Partition: 0,
		Leader:    node,
		Replicas:  []*Node{node}, // Ensure the node is a replica
	}
	partition2 := &PartitionInfo{
		Topic:     "test-topic",
		Partition: 1,         // No leader
		Replicas:  []*Node{}, // No replicas
	}

	cluster.AddPartition(partition1)
	cluster.AddPartition(partition2)

	availablePartitions := cluster.AvailablePartitionsForTopic("test-topic")
	if len(availablePartitions) != 1 {
		t.Fatalf("Expected 1 available partition for topic 'test-topic', got %d", len(availablePartitions))
	}
}

// Helper function to create a sample Cluster
func createSampleCluster() *Cluster {
	return &Cluster{
		clusterResource:       &ClusterResource{clusterID: "test-cluster"},
		unauthorizedTopics:    map[string]struct{}{"topic1": {}},
		invalidTopics:         map[string]struct{}{"topic2": {}},
		isBootstrapConfigured: true,
		nodes:                 []*Node{{ID: 1, Host: "localhost:9092"}},
	}
}

// TestString tests the String() method of Cluster
func TestClusterString(t *testing.T) {
	cluster := createSampleCluster()
	expected := "Cluster{id='test-cluster', unauthorizedTopics=map[topic1:{}], invalidTopics=map[topic2:{}], isBootstrapConfigured=true, nodes=[{ID:1 Host:localhost:9092}]}"
	if got := cluster.String(); got != expected {
		t.Errorf("String() = %v, want %v", got, expected)
	}
}

// TestEquals tests the Equals() method of Cluster
func TestClusterEquals(t *testing.T) {
	cluster1 := createSampleCluster()
	cluster2 := createSampleCluster()

	if !cluster1.Equals(cluster2) {
		t.Errorf("Equals() = false, want true")
	}

	cluster2.isBootstrapConfigured = false
	if cluster1.Equals(cluster2) {
		t.Errorf("Equals() = true, want false")
	}
}

// TestHashCode tests the HashCode() method of Cluster
func TestClusterHashCode(t *testing.T) {
	cluster1 := createSampleCluster()
	cluster2 := createSampleCluster()

	if cluster1.HashCode() != cluster2.HashCode() {
		t.Errorf("HashCode() = %v, want %v", cluster1.HashCode(), cluster2.HashCode())
	}

	cluster2.isBootstrapConfigured = false
	if cluster1.HashCode() == cluster2.HashCode() {
		t.Errorf("HashCode() = %v, want different", cluster2.HashCode())
	}
}

// TestIsBootstrapConfigured tests the IsBootstrapConfigured() method of Cluster
func TestIsBootstrapConfigured(t *testing.T) {
	cluster := createSampleCluster()
	if got := cluster.IsBootstrapConfigured(); !got {
		t.Errorf("IsBootstrapConfigured() = %v, want true", got)
	}

	// Testing with a cluster where bootstrap is not configured
	cluster.isBootstrapConfigured = false
	if got := cluster.IsBootstrapConfigured(); got {
		t.Errorf("IsBootstrapConfigured() = %v, want false", got)
	}
}
