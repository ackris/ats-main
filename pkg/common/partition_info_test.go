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
	"reflect"
	"testing"
)

// TestFormatNodeIDs tests the FormatNodeIDs function.
func TestFormatNodeIDs(t *testing.T) {
	tests := []struct {
		nodes []*Node
		want  string
	}{
		{nil, "[]"},
		{[]*Node{}, "[]"},
		{[]*Node{{ID: 1}, {ID: 2}}, "[1,2]"},
		{[]*Node{{ID: 1}, nil, {ID: 2}}, "[1,2]"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := FormatNodeIDs(tt.nodes)
			if got != tt.want {
				t.Errorf("FormatNodeIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestPartitionInfoString tests the String method of PartitionInfo.
func TestPartitionInfoString(t *testing.T) {
	node1 := &Node{ID: 1, Host: "localhost", Port: 9092}
	node2 := &Node{ID: 2, Host: "localhost", Port: 9093}
	partitionInfo := &PartitionInfo{
		Topic:           "topic1",
		Partition:       0,
		Leader:          node1,
		Replicas:        []*Node{node1, node2},
		InSyncReplicas:  []*Node{node1},
		OfflineReplicas: []*Node{node2},
	}

	want := "Partition(topic = topic1, partition = 0, leader = 1, replicas = [1,2], isr = [1], offlineReplicas = [2])"
	got := partitionInfo.String()
	if got != want {
		t.Errorf("String() = %v, want %v", got, want)
	}
}

// TestPartitionInfoEquals tests the Equals method of PartitionInfo.
func TestPartitionInfoEquals(t *testing.T) {
	node1 := &Node{ID: 1}
	node2 := &Node{ID: 2}

	p1 := &PartitionInfo{
		Topic:           "topic1",
		Partition:       0,
		Leader:          node1,
		Replicas:        []*Node{node1, node2},
		InSyncReplicas:  []*Node{node1},
		OfflineReplicas: []*Node{node2},
	}

	p2 := &PartitionInfo{
		Topic:           "topic1",
		Partition:       0,
		Leader:          node1,
		Replicas:        []*Node{node1, node2},
		InSyncReplicas:  []*Node{node1},
		OfflineReplicas: []*Node{node2},
	}

	p3 := &PartitionInfo{
		Topic:           "topic1",
		Partition:       1,
		Leader:          node1,
		Replicas:        []*Node{node1, node2},
		InSyncReplicas:  []*Node{node1},
		OfflineReplicas: []*Node{node2},
	}

	tests := []struct {
		a, b *PartitionInfo
		want bool
	}{
		{p1, p2, true},
		{p1, p3, false},
		{p1, nil, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := tt.a.Equals(tt.b)
			if got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSlicesEqual tests the SlicesEqual function.
func TestSlicesEqual(t *testing.T) {
	node1 := &Node{ID: 1}
	node2 := &Node{ID: 2}

	tests := []struct {
		a, b []*Node
		want bool
	}{
		{nil, nil, true},
		{[]*Node{}, []*Node{}, true},
		{[]*Node{node1}, []*Node{node1}, true},
		{[]*Node{node1, node2}, []*Node{node1, node2}, true},
		{[]*Node{node1}, []*Node{node2}, false},
		{[]*Node{node1, nil}, []*Node{node1, nil}, true},
		{[]*Node{node1, nil}, []*Node{nil, node1}, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := SlicesEqual(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("SlicesEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPartitionInfo(t *testing.T) {
	leaderNode := &Node{ID: 1}
	replicaNodes := []*Node{{ID: 2}, {ID: 3}}
	inSyncNodes := []*Node{{ID: 2}}

	partitionInfo := NewPartitionInfo("example-topic", 0, leaderNode, replicaNodes, inSyncNodes)

	if partitionInfo.GetTopic() != "example-topic" {
		t.Errorf("Expected topic to be 'example-topic', got %s", partitionInfo.GetTopic())
	}

	if partitionInfo.GetPartition() != 0 {
		t.Errorf("Expected partition to be 0, got %d", partitionInfo.GetPartition())
	}

	if partitionInfo.GetLeader() != leaderNode {
		t.Errorf("Expected leader to be %v, got %v", leaderNode, partitionInfo.GetLeader())
	}

	if !reflect.DeepEqual(partitionInfo.GetReplicas(), replicaNodes) {
		t.Errorf("Expected replicas to be %v, got %v", replicaNodes, partitionInfo.GetReplicas())
	}

	if !reflect.DeepEqual(partitionInfo.GetInSyncReplicas(), inSyncNodes) {
		t.Errorf("Expected in-sync replicas to be %v, got %v", inSyncNodes, partitionInfo.GetInSyncReplicas())
	}

	if len(partitionInfo.GetOfflineReplicas()) != 0 {
		t.Errorf("Expected offline replicas to be empty, got %v", partitionInfo.GetOfflineReplicas())
	}
}

func TestNewPartitionInfoWithOffline(t *testing.T) {
	leaderNode := &Node{ID: 1}
	replicaNodes := []*Node{{ID: 2}, {ID: 3}}
	inSyncNodes := []*Node{{ID: 2}}
	offlineNodes := []*Node{{ID: 3}}

	partitionInfo := NewPartitionInfoWithOffline("example-topic", 0, leaderNode, replicaNodes, inSyncNodes, offlineNodes)

	if partitionInfo.GetTopic() != "example-topic" {
		t.Errorf("Expected topic to be 'example-topic', got %s", partitionInfo.GetTopic())
	}

	if partitionInfo.GetPartition() != 0 {
		t.Errorf("Expected partition to be 0, got %d", partitionInfo.GetPartition())
	}

	if partitionInfo.GetLeader() != leaderNode {
		t.Errorf("Expected leader to be %v, got %v", leaderNode, partitionInfo.GetLeader())
	}

	if !reflect.DeepEqual(partitionInfo.GetReplicas(), replicaNodes) {
		t.Errorf("Expected replicas to be %v, got %v", replicaNodes, partitionInfo.GetReplicas())
	}

	if !reflect.DeepEqual(partitionInfo.GetInSyncReplicas(), inSyncNodes) {
		t.Errorf("Expected in-sync replicas to be %v, got %v", inSyncNodes, partitionInfo.GetInSyncReplicas())
	}

	if !reflect.DeepEqual(partitionInfo.GetOfflineReplicas(), offlineNodes) {
		t.Errorf("Expected offline replicas to be %v, got %v", offlineNodes, partitionInfo.GetOfflineReplicas())
	}
}

func TestHash(t *testing.T) {
	leaderNode := &Node{ID: 1}
	replicaNodes := []*Node{{ID: 2}, {ID: 3}}
	inSyncNodes := []*Node{{ID: 2}}

	partitionInfo := NewPartitionInfo("example-topic", 0, leaderNode, replicaNodes, inSyncNodes)
	hash1 := partitionInfo.Hash()

	partitionInfo = NewPartitionInfo("example-topic", 0, leaderNode, replicaNodes, inSyncNodes)
	hash2 := partitionInfo.Hash()

	if hash1 != hash2 {
		t.Errorf("Expected hashes to be equal, got %d and %d", hash1, hash2)
	}
}
