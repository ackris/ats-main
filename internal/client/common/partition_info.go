package common

import (
	"fmt"
	"strings"
)

// Node represents a node in the cluster.
// It holds the ID of the node.
type Node struct {
	ID string
}

// IDString returns the string representation of the node's ID.
// This method is used to get the unique identifier of the node as a string.
//
// Returns:
//
//	The ID of the node as a string.
//
// Example:
//
//	node := &Node{ID: "node1"}
//	fmt.Println(node.IDString()) // Output: node1
func (n *Node) IDString() string {
	return n.ID
}

// PartitionInfo describes the state of a partition within a topic.
// It includes information about the topic, partition ID, leader node,
// replicas, in-sync replicas, and offline replicas.
type PartitionInfo struct {
	Topic           string  // The name of the topic
	Partition       int     // The partition ID
	Leader          *Node   // The node currently acting as the leader for this partition (nil if no leader)
	Replicas        []*Node // The complete set of replicas for this partition
	InSyncReplicas  []*Node // The subset of replicas that are in sync with the leader
	OfflineReplicas []*Node // The subset of replicas that are offline
}

// FormatNodeIDs formats the node IDs from a slice of Node pointers for display.
// This method creates a string representation of node IDs, enclosed in square brackets,
// and separates each ID by a comma.
//
// Parameters:
//
//	nodes: A slice of Node pointers whose IDs are to be formatted.
//
// Returns:
//
//	A string representing the formatted node IDs.
//
// Example:
//
//	nodes := []*Node{
//	    &Node{ID: "node1"},
//	    &Node{ID: "node2"},
//	}
//	fmt.Println(FormatNodeIDs(nodes)) // Output: [node1,node2]
func FormatNodeIDs(nodes []*Node) string {
	if nodes == nil {
		return "[]"
	}
	var ids []string
	for _, node := range nodes {
		if node != nil {
			ids = append(ids, node.IDString())
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(ids, ","))
}

// String returns a string representation of the PartitionInfo instance.
// This method formats the partition information into a readable string format,
// including topic name, partition ID, leader, replicas, in-sync replicas, and offline replicas.
//
// Returns:
//
//	A string representation of the PartitionInfo.
//
// Example:
//
//	node1 := &Node{ID: "node1"}
//	node2 := &Node{ID: "node2"}
//	partitionInfo := &PartitionInfo{
//	    Topic:           "topic1",
//	    Partition:       0,
//	    Leader:          node1,
//	    Replicas:        []*Node{node1, node2},
//	    InSyncReplicas:  []*Node{node1},
//	    OfflineReplicas: []*Node{node2},
//	}
//	fmt.Println(partitionInfo.String())
//	// Output: Partition(topic = topic1, partition = 0, leader = node1, replicas = [node1,node2], isr = [node1], offlineReplicas = [node2])
func (p *PartitionInfo) String() string {
	leaderID := "none"
	if p.Leader != nil {
		leaderID = p.Leader.IDString()
	}
	return fmt.Sprintf(
		"Partition(topic = %s, partition = %d, leader = %s, replicas = %s, isr = %s, offlineReplicas = %s)",
		p.Topic,
		p.Partition,
		leaderID,
		FormatNodeIDs(p.Replicas),
		FormatNodeIDs(p.InSyncReplicas),
		FormatNodeIDs(p.OfflineReplicas),
	)
}

// Equals compares two PartitionInfo instances for equality.
// This method checks if the other PartitionInfo has the same topic, partition ID,
// leader, replicas, in-sync replicas, and offline replicas as the current instance.
//
// Parameters:
//
//	other: The PartitionInfo instance to compare with the current instance.
//
// Returns:
//
//	True if the two PartitionInfo instances are equal, false otherwise.
//
// Example:
//
//	node1 := &Node{ID: "node1"}
//	node2 := &Node{ID: "node2"}
//	p1 := &PartitionInfo{
//	    Topic:           "topic1",
//	    Partition:       0,
//	    Leader:          node1,
//	    Replicas:        []*Node{node1, node2},
//	    InSyncReplicas:  []*Node{node1},
//	    OfflineReplicas: []*Node{node2},
//	}
//	p2 := &PartitionInfo{
//	    Topic:           "topic1",
//	    Partition:       0,
//	    Leader:          node1,
//	    Replicas:        []*Node{node1, node2},
//	    InSyncReplicas:  []*Node{node1},
//	    OfflineReplicas: []*Node{node2},
//	}
//	fmt.Println(p1.Equals(p2)) // Output: true
func (p *PartitionInfo) Equals(other *PartitionInfo) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if p.Topic != other.Topic || p.Partition != other.Partition {
		return false
	}
	if (p.Leader == nil) != (other.Leader == nil) {
		return false
	}
	if p.Leader != nil && p.Leader.IDString() != other.Leader.IDString() {
		return false
	}
	if !SlicesEqual(p.Replicas, other.Replicas) {
		return false
	}
	if !SlicesEqual(p.InSyncReplicas, other.InSyncReplicas) {
		return false
	}
	if !SlicesEqual(p.OfflineReplicas, other.OfflineReplicas) {
		return false
	}
	return true
}

// SlicesEqual checks if two slices of Node pointers are equal.
// This method compares the slices by their contents, taking into account
// the possible presence of nil values in the slices.
//
// Parameters:
//
//	a: The first slice of Node pointers.
//	b: The second slice of Node pointers.
//
// Returns:
//
//	True if both slices contain the same elements, false otherwise.
//
// Example:
//
//	a := []*Node{
//	    &Node{ID: "node1"},
//	    &Node{ID: "node2"},
//	}
//	b := []*Node{
//	    &Node{ID: "node1"},
//	    &Node{ID: "node2"},
//	}
//	fmt.Println(SlicesEqual(a, b)) // Output: true
func SlicesEqual(a, b []*Node) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// contains checks if a map contains the specified key.
// This method is used to check if a key is present in a map.
//
// Parameters:
//
//	m: The map to check for the presence of the key.
//	key: The key to look for in the map.
//
// Returns:
//
//	True if the key is in the map, false otherwise.
//
// Example:
//
//	m := map[string]struct{}{
//	    "node1": {},
//	    "node2": {},
//	}
//	fmt.Println(contains(m, "node1")) // Output: true
//	fmt.Println(contains(m, "node3")) // Output: false
func contains(m map[string]struct{}, key string) bool {
	_, ok := m[key]
	return ok
}
