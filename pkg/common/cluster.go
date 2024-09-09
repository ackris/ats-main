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
	"errors"
	"fmt"
	"hash/fnv"
	"reflect"
	"sort"
	"sync"
)

// Cluster represents an immutable representation of a Atomstate cluster.
type Cluster struct {
	isBootstrapConfigured      bool
	nodes                      []*Node
	unauthorizedTopics         map[string]struct{}
	invalidTopics              map[string]struct{}
	internalTopics             map[string]struct{}
	controller                 *Node
	partitionsByTopicPartition map[*TopicPartition]*PartitionInfo
	partitionsByTopic          map[string][]*PartitionInfo
	availablePartitionsByTopic map[string][]*PartitionInfo
	partitionsByNode           map[int][]*PartitionInfo
	nodesById                  map[int]*Node
	clusterResource            *ClusterResource
	topicIds                   map[string]Uuid
	topicNames                 map[Uuid]string
	mu                         sync.RWMutex // for concurrent safety
}

// NewCluster creates a new Cluster instance with the given parameters.
//
// Parameters:
//   - clusterID: The unique ID of the cluster.
//   - isBootstrapConfigured: Flag indicating if the bootstrap configuration is set.
//   - nodes: List of nodes in the cluster.
//   - partitionsByTopicPartition: Mapping of TopicPartition to PartitionInfo.
//   - unauthorizedTopics: Set of topics that are unauthorized.
//   - invalidTopics: Set of invalid topics.
//   - internalTopics: Set of internal topics.
//   - controller: The controller node of the cluster.
//   - topicIds: Mapping of topic names to their UUIDs.
//
// Returns:
//   - A new instance of Cluster.
//
// Example:
//
// cluster := NewCluster(
//
//	"cluster-id",
//	true,
//	[]*Node{
//		{ID: 1, Host: "localhost:9092"},
//	},
//	map[*TopicPartition]*PartitionInfo{
//		&TopicPartition{topic: "my-topic", partition: 0}: &PartitionInfo{Topic: "my-topic", Partition: 0},
//	},
//	map[string]struct{}{"unauthorized-topic": {}},
//	map[string]struct{}{"invalid-topic": {}},
//	map[string]struct{}{"internal-topic": {}},
//	&Node{ID: 1, Host: "localhost:9092"},
//	map[string]Uuid{"my-topic": "uuid-1"},
//
// )
func NewCluster(
	clusterID string,
	isBootstrapConfigured bool,
	nodes []*Node,
	partitionsByTopicPartition map[*TopicPartition]*PartitionInfo,
	unauthorizedTopics map[string]struct{},
	invalidTopics map[string]struct{},
	internalTopics map[string]struct{},
	controller *Node,
	topicIds map[string]Uuid,
) *Cluster {
	cl := &Cluster{
		isBootstrapConfigured:      isBootstrapConfigured,
		nodes:                      make([]*Node, len(nodes)),
		unauthorizedTopics:         copySet(unauthorizedTopics),
		invalidTopics:              copySet(invalidTopics),
		internalTopics:             copySet(internalTopics),
		controller:                 controller,
		partitionsByTopicPartition: make(map[*TopicPartition]*PartitionInfo),
		partitionsByTopic:          make(map[string][]*PartitionInfo),
		availablePartitionsByTopic: make(map[string][]*PartitionInfo),
		partitionsByNode:           make(map[int][]*PartitionInfo),
		nodesById:                  make(map[int]*Node),
		clusterResource:            &ClusterResource{clusterID: clusterID},
		topicIds:                   copyMap(topicIds),
		topicNames:                 make(map[Uuid]string),
	}

	// Copy nodes and sort by ID
	copy(cl.nodes, nodes)
	sort.Slice(cl.nodes, func(i, j int) bool { return cl.nodes[i].ID < cl.nodes[j].ID })

	// Index nodes
	for _, node := range cl.nodes {
		cl.nodesById[node.ID] = node
		cl.partitionsByNode[node.ID] = nil // Initialize as nil to avoid unnecessary allocations
	}

	// Index partitions
	for tp, p := range partitionsByTopicPartition {
		cl.partitionsByTopicPartition[tp] = p
		cl.partitionsByTopic[p.Topic] = append(cl.partitionsByTopic[p.Topic], p)

		if p.Leader != nil {
			cl.partitionsByNode[p.Leader.ID] = append(cl.partitionsByNode[p.Leader.ID], p)
		}
	}

	// Compute available partitions
	cl.updateAvailablePartitions()

	// Create topic names map
	for topic, id := range topicIds {
		cl.topicNames[id] = topic
	}

	return cl
}

// copySet creates a copy of the set.
//
// Parameters:
//   - original: The original set to copy.
//
// Returns:
//   - A new copy of the set.
func copySet(original map[string]struct{}) map[string]struct{} {
	copy := make(map[string]struct{}, len(original))
	for k := range original {
		copy[k] = struct{}{}
	}
	return copy
}

// copyMap creates a copy of the map.
//
// Parameters:
//   - original: The original map to copy.
//
// Returns:
//   - A new copy of the map.
func copyMap(original map[string]Uuid) map[string]Uuid {
	copy := make(map[string]Uuid, len(original))
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

// EmptyCluster creates an empty Cluster instance.
//
// Returns:
//   - An empty Cluster instance.
func EmptyCluster() *Cluster {
	return &Cluster{
		nodes:                      []*Node{}, // Initialize as an empty slice instead of nil
		unauthorizedTopics:         make(map[string]struct{}),
		invalidTopics:              make(map[string]struct{}),
		internalTopics:             make(map[string]struct{}),
		partitionsByTopicPartition: make(map[*TopicPartition]*PartitionInfo),
		partitionsByTopic:          make(map[string][]*PartitionInfo),
		availablePartitionsByTopic: make(map[string][]*PartitionInfo),
		partitionsByNode:           make(map[int][]*PartitionInfo), // No need for initial capacity here
		nodesById:                  make(map[int]*Node),
		clusterResource:            &ClusterResource{clusterID: ""},
		topicIds:                   make(map[string]Uuid),
		topicNames:                 make(map[Uuid]string),
	}
}

// BootstrapCluster creates a bootstrap Cluster with the given list of addresses.
//
// Parameters:
//   - addresses: List of node addresses to be used for bootstrap.
//
// Returns:
//   - A new bootstrap Cluster instance.
func BootstrapCluster(addresses []string) *Cluster {
	const negativeIDOffset = -1 // Define a constant for negative ID offset
	nodes := make([]*Node, len(addresses))
	for i, addr := range addresses {
		nodes[i] = &Node{ID: negativeIDOffset * (i + 1), Host: addr}
	}
	return NewCluster("", true, nodes, nil, nil, nil, nil, nil, nil)
}

// WithPartitions returns a new Cluster with updated partitions.
//
// Parameters:
//   - partitions: Mapping of TopicPartition to PartitionInfo to be updated.
//
// Returns:
//   - A new Cluster instance with the updated partitions.
func (cl *Cluster) WithPartitions(partitions map[*TopicPartition]*PartitionInfo) *Cluster {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	// Create a new copy of the partitionsByTopicPartition map
	newPartitions := copyPartitionsMap(cl.partitionsByTopicPartition)
	for k, v := range partitions {
		newPartitions[k] = v
	}

	// Create a new Cluster with the updated partitions
	return NewCluster(cl.clusterResource.clusterID, cl.isBootstrapConfigured, cl.nodes, newPartitions, cl.unauthorizedTopics, cl.invalidTopics, cl.internalTopics, cl.controller, cl.topicIds)
}

// copyPartitionsMap creates a copy of the map[*TopicPartition]*PartitionInfo.
//
// Parameters:
//   - original: The original map to copy.
//
// Returns:
//   - A new copy of the map.
func copyPartitionsMap(original map[*TopicPartition]*PartitionInfo) map[*TopicPartition]*PartitionInfo {
	copy := make(map[*TopicPartition]*PartitionInfo, len(original))
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

// Nodes returns the list of known nodes in the cluster.
//
// Returns:
//   - A slice of nodes in the cluster.
func (cl *Cluster) Nodes() []*Node {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	nodes := make([]*Node, len(cl.nodes))
	copy(nodes, cl.nodes)
	return nodes
}

// NodeByID returns the node with the specified ID.
//
// Parameters:
//   - id: The ID of the node to retrieve.
//
// Returns:
//   - The node with the given ID, or nil if not found.
func (cl *Cluster) NodeByID(id int) *Node {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return cl.nodesById[id]
}

// NodeIfOnline returns the node if it's online for the given partition.
//
// Parameters:
//   - partition: The partition for which to check the node's status.
//   - id: The ID of the node to check.
//
// Returns:
//   - The node if it is online for the partition, otherwise nil.
func (cl *Cluster) NodeIfOnline(partition *TopicPartition, id int) *Node {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	node := cl.NodeByID(id)
	if node == nil {
		return nil
	}

	p, exists := cl.partitionsByTopicPartition[partition]
	if !exists {
		return nil
	}

	// Check if the node is offline or not a replica
	for _, offline := range p.OfflineReplicas {
		if offline.ID == id {
			return nil
		}
	}
	for _, replica := range p.Replicas {
		if replica.ID == id {
			return node
		}
	}
	return nil
}

// LeaderFor returns the leader node for the given topic-partition.
//
// Parameters:
//   - tp: The TopicPartition for which to find the leader.
//
// Returns:
//   - The leader node for the topic-partition, or nil if not found.
func (cl *Cluster) LeaderFor(tp *TopicPartition) *Node {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	p, exists := cl.partitionsByTopicPartition[tp]
	if !exists {
		return nil
	}
	return p.Leader
}

// Partition returns the metadata for the specified partition.
//
// Parameters:
//   - tp: The TopicPartition for which to retrieve the partition info.
//
// Returns:
//   - The PartitionInfo for the specified TopicPartition, or nil if not found.
func (cl *Cluster) Partition(tp *TopicPartition) *PartitionInfo {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return cl.partitionsByTopicPartition[tp]
}

// PartitionsForTopic returns the list of partitions for a given topic.
//
// Parameters:
//   - topic: The name of the topic for which to retrieve partitions.
//
// Returns:
//   - A slice of PartitionInfo for the specified topic.
func (cl *Cluster) PartitionsForTopic(topic string) []*PartitionInfo {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	partitions := cl.partitionsByTopic[topic]
	return append([]*PartitionInfo(nil), partitions...) // Return a copy of the slice
}

// PartitionCountForTopic returns the number of partitions for a given topic.
//
// Parameters:
//   - topic: The name of the topic for which to count partitions.
//
// Returns:
//   - The number of partitions for the specified topic.
func (cl *Cluster) PartitionCountForTopic(topic string) int {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return len(cl.partitionsByTopic[topic])
}

// AvailablePartitionsForTopic returns the list of available partitions for a given topic.
//
// Parameters:
//   - topic: The name of the topic for which to retrieve available partitions.
//
// Returns:
//   - A slice of available PartitionInfo for the specified topic.
func (cl *Cluster) AvailablePartitionsForTopic(topic string) []*PartitionInfo {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	partitions := cl.availablePartitionsByTopic[topic]
	return append([]*PartitionInfo(nil), partitions...) // Return a copy of the slice
}

// PartitionsForNode returns the list of partitions for the given node ID.
//
// Parameters:
//   - nodeID: The ID of the node for which to retrieve partitions.
//
// Returns:
//   - A slice of PartitionInfo for the specified node.
func (cl *Cluster) PartitionsForNode(nodeID int) []*PartitionInfo {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	partitions := cl.partitionsByNode[nodeID]
	return append([]*PartitionInfo(nil), partitions...) // Return a copy of the slice
}

// Topics returns the set of all topics in the cluster.
//
// Returns:
//   - A map of topic names to empty structs representing the set of topics.
func (cl *Cluster) Topics() map[string]struct{} {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	topics := make(map[string]struct{}, len(cl.partitionsByTopic))
	for topic := range cl.partitionsByTopic {
		topics[topic] = struct{}{}
	}
	return topics
}

// Controller returns the controller node of the cluster.
//
// Returns:
//   - The controller node, or nil if not set.
func (cl *Cluster) Controller() *Node {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return cl.controller
}

// ClusterID returns the ID of the cluster.
//
// Returns:
//   - The ID of the cluster.
func (cl *Cluster) ClusterID() string {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return cl.clusterResource.clusterID
}

// TopicID returns the UUID for a given topic.
//
// Parameters:
//   - topic: The name of the topic for which to retrieve the UUID.
//
// Returns:
//   - The UUID of the topic and a boolean indicating if it exists.
func (cl *Cluster) TopicID(topic string) (Uuid, bool) {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	id, exists := cl.topicIds[topic]
	return id, exists
}

// TopicName returns the topic name for a given UUID.
//
// Parameters:
//   - id: The UUID of the topic for which to retrieve the name.
//
// Returns:
//   - The name of the topic and a boolean indicating if it exists.
func (cl *Cluster) TopicName(id Uuid) (string, bool) {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	name, exists := cl.topicNames[id]
	return name, exists
}

// AddNode adds a new node to the cluster.
//
// Parameters:
//   - node: The Node to be added.
//
// Returns:
//   - An error if the node already exists, otherwise nil.
//
// Example usage:
//
//	cluster := EmptyCluster()
//	err := cluster.AddNode(&Node{ID: 1, Host: "localhost:9092"})
//	if err != nil {
//	    log.Fatal(err)
//	}
func (cl *Cluster) AddNode(node *Node) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	if _, exists := cl.nodesById[node.ID]; exists {
		return errors.New("node with the same ID already exists") // Return an error if the node exists
	}

	cl.nodes = append(cl.nodes, node)
	cl.nodesById[node.ID] = node
	cl.partitionsByNode[node.ID] = nil // Initialize as nil to avoid unnecessary allocations
	return nil
}

// AddPartition adds a new partition to the cluster.
//
// Parameters:
//   - partition: The PartitionInfo to be added.
//
// Returns:
//   - An error if the partition already exists, otherwise nil.
//
// Example usage:
//
//	cluster := EmptyCluster()
//	err := cluster.AddPartition(&PartitionInfo{Topic: "my-topic", Partition: 0})
//	if err != nil {
//	    log.Fatal(err)
//	}
func (cl *Cluster) AddPartition(partition *PartitionInfo) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	tp := NewTopicPartition(partition.Topic, partition.Partition)
	if existingPartition, exists := cl.partitionsByTopicPartition[tp]; exists && existingPartition.Equals(partition) {
		return errors.New("partition already exists") // Correctly handle duplicate partitions
	}

	cl.partitionsByTopicPartition[tp] = partition
	cl.partitionsByTopic[partition.Topic] = append(cl.partitionsByTopic[partition.Topic], partition)

	if partition.Leader != nil {
		cl.partitionsByNode[partition.Leader.ID] = append(cl.partitionsByNode[partition.Leader.ID], partition)
	}

	cl.updateAvailablePartitions()
	return nil
}

// updateAvailablePartitions updates the available partitions for each topic.
func (cl *Cluster) updateAvailablePartitions() {
	for topic, partitions := range cl.partitionsByTopic {
		var availablePartitions []*PartitionInfo
		for _, p := range partitions {
			if p.Leader != nil {
				availablePartitions = append(availablePartitions, p)
			}
		}
		cl.availablePartitionsByTopic[topic] = availablePartitions
	}
}

// RemoveNode removes a node from the cluster.
//
// Parameters:
//   - nodeID: The ID of the node to remove.
//
// Example usage:
//
//	cluster := EmptyCluster()
//	cluster.RemoveNode(1)
func (cl *Cluster) RemoveNode(nodeID int) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	if _, exists := cl.nodesById[nodeID]; !exists {
		return
	}

	// Create a new slice to hold remaining nodes
	newNodes := make([]*Node, 0, len(cl.nodes)-1)
	for _, n := range cl.nodes {
		if n.ID != nodeID {
			newNodes = append(newNodes, n)
		}
	}
	cl.nodes = newNodes

	delete(cl.nodesById, nodeID)

	// Remove partitions associated with the node
	for topic, partitions := range cl.partitionsByTopic {
		var remainingPartitions []*PartitionInfo
		for _, p := range partitions {
			if p.Leader != nil && p.Leader.ID == nodeID {
				continue
			}
			remainingPartitions = append(remainingPartitions, p)
		}
		cl.partitionsByTopic[topic] = remainingPartitions
	}

	delete(cl.partitionsByNode, nodeID)

	// Update available partitions after node removal
	cl.updateAvailablePartitions()
}

// RemovePartition removes a partition from the cluster.
//
// Parameters:
//   - topic: The name of the topic of the partition to remove.
//   - partition: The partition number to remove.
//
// Example usage:
//
//	cluster := EmptyCluster()
//	cluster.RemovePartition("my-topic", 0)
func (cl *Cluster) RemovePartition(topic string, partition int) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	tp := &TopicPartition{topic: topic, partition: partition}
	delete(cl.partitionsByTopicPartition, tp)

	partitions := cl.partitionsByTopic[topic]
	var remainingPartitions []*PartitionInfo
	for _, p := range partitions {
		if p.Partition != partition {
			remainingPartitions = append(remainingPartitions, p)
		}
	}
	cl.partitionsByTopic[topic] = remainingPartitions

	if len(remainingPartitions) == 0 {
		delete(cl.availablePartitionsByTopic, topic)
	}

	// Remove partitions from nodes
	for nodeID, nodePartitions := range cl.partitionsByNode {
		var updatedPartitions []*PartitionInfo
		for _, p := range nodePartitions {
			if p.Partition != partition {
				updatedPartitions = append(updatedPartitions, p)
			}
		}
		cl.partitionsByNode[nodeID] = updatedPartitions
	}

	// Update available partitions after partition removal
	cl.updateAvailablePartitions()
}

// String returns a formatted string representation of the Cluster instance.
// It includes the cluster ID, unauthorized topics, invalid topics,
// bootstrap configuration status, and the list of nodes.
//
// Example:
//
//	cluster := &Cluster{
//		clusterResource: &ClusterResource{clusterID: "test-cluster"},
//		unauthorizedTopics: map[string]struct{}{"topic1": {}},
//		invalidTopics: map[string]struct{}{"topic2": {}},
//		isBootstrapConfigured: true,
//		nodes: []*Node{{ID: 1, Host: "localhost:9092"}},
//	}
//	fmt.Println(cluster.String())
//
// Output: Cluster{id='test-cluster', unauthorizedTopics=map[topic1:{}], invalidTopics=map[topic2:{}], isBootstrapConfigured=true, nodes=[{ID:1 Host:localhost:9092}]}
func (cl *Cluster) String() string {
	return fmt.Sprintf("Cluster{id='%s', unauthorizedTopics=%v, invalidTopics=%v, isBootstrapConfigured=%v, nodes=%v}",
		cl.clusterResource.clusterID,
		cl.unauthorizedTopics,
		cl.invalidTopics,
		cl.isBootstrapConfigured,
		cl.nodes)
}

// Equals compares the current Cluster instance with another Cluster instance.
// Two Cluster instances are considered equal if they have the same cluster ID,
// bootstrap configuration status, unauthorized topics, invalid topics, and nodes.
//
// Parameters:
//
//	other (*Cluster): The other Cluster instance to compare with.
//
// Returns:
//
//	bool: true if both Cluster instances are equal, false otherwise.
//
// Example:
//
//	cluster1 := &Cluster{
//		clusterResource: &ClusterResource{clusterID: "test-cluster"},
//		unauthorizedTopics: map[string]struct{}{"topic1": {}},
//		invalidTopics: map[string]struct{}{"topic2": {}},
//		isBootstrapConfigured: true,
//		nodes: []*Node{{ID: 1, Host: "localhost:9092"}},
//	}
//
//	cluster2 := &Cluster{
//		clusterResource: &ClusterResource{clusterID: "test-cluster"},
//		unauthorizedTopics: map[string]struct{}{"topic1": {}},
//		invalidTopics: map[string]struct{}{"topic2": {}},
//		isBootstrapConfigured: true,
//		nodes: []*Node{{ID: 1, Host: "localhost:9092"}},
//	}
//
//	fmt.Println(cluster1.Equals(cluster2)) // Output: true
func (cl *Cluster) Equals(other *Cluster) bool {
	if other == nil {
		return false
	}

	// Compare fields for equality
	return cl.clusterResource.clusterID == other.clusterResource.clusterID &&
		cl.isBootstrapConfigured == other.isBootstrapConfigured &&
		reflect.DeepEqual(cl.unauthorizedTopics, other.unauthorizedTopics) &&
		reflect.DeepEqual(cl.invalidTopics, other.invalidTopics) &&
		reflect.DeepEqual(cl.nodes, other.nodes)
}

// HashCode generates a hash code for the Cluster instance based on its
// cluster ID and bootstrap configuration status. This hash code can be used
// for hashing and lookup in hash-based data structures.
//
// Returns:
//
//	uint64: The hash code for the Cluster instance.
//
// Example:
//
//	cluster := &Cluster{
//		clusterResource: &ClusterResource{clusterID: "test-cluster"},
//		isBootstrapConfigured: true,
//	}
//
//	hashCode := cluster.HashCode()
//	fmt.Println(hashCode) // Output: A hash code as uint64
func (cl *Cluster) HashCode() uint64 {
	h := fnv.New64a()
	h.Write([]byte(cl.clusterResource.clusterID))

	if cl.isBootstrapConfigured {
		h.Write([]byte("true"))
	} else {
		h.Write([]byte("false"))
	}

	// Consider including other fields in the hash if necessary
	return h.Sum64()
}

// IsBootstrapConfigured returns whether the bootstrap configuration is enabled for the cluster.
//
// Returns:
//
//	bool: true if bootstrap is configured, false otherwise.
//
// Example:
//
//	cluster := &Cluster{
//		isBootstrapConfigured: true,
//	}
//
//	fmt.Println(cluster.IsBootstrapConfigured()) // Output: true
func (cl *Cluster) IsBootstrapConfigured() bool {
	return cl.isBootstrapConfigured
}
