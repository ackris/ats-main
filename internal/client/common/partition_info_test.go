package common

import "testing"

// TestNode_IDString tests the IDString method of the Node struct.
func TestNode_IDString(t *testing.T) {
	node := &Node{ID: "node1"}
	expected := "node1"
	if got := node.IDString(); got != expected {
		t.Errorf("Node.IDString() = %v, want %v", got, expected)
	}
}

// TestFormatNodeIDs tests the FormatNodeIDs function.
func TestFormatNodeIDs(t *testing.T) {
	tests := []struct {
		nodes    []*Node
		expected string
	}{
		{
			nodes:    []*Node{{ID: "node1"}, {ID: "node2"}},
			expected: "[node1,node2]",
		},
		{
			nodes:    nil,
			expected: "[]",
		},
	}

	for _, tt := range tests {
		if got := FormatNodeIDs(tt.nodes); got != tt.expected {
			t.Errorf("FormatNodeIDs(%v) = %v, want %v", tt.nodes, got, tt.expected)
		}
	}
}

// TestPartitionInfo_String tests the String method of the PartitionInfo struct.
func TestPartitionInfo_String(t *testing.T) {
	node1 := &Node{ID: "node1"}
	node2 := &Node{ID: "node2"}

	partitionInfo := &PartitionInfo{
		Topic:           "topic1",
		Partition:       0,
		Leader:          node1,
		Replicas:        []*Node{node1, node2},
		InSyncReplicas:  []*Node{node1},
		OfflineReplicas: []*Node{node2},
	}

	expected := "Partition(topic = topic1, partition = 0, leader = node1, replicas = [node1,node2], isr = [node1], offlineReplicas = [node2])"
	if got := partitionInfo.String(); got != expected {
		t.Errorf("PartitionInfo.String() = %v, want %v", got, expected)
	}
}

// TestPartitionInfo_Equals tests the Equals method of the PartitionInfo struct.
func TestPartitionInfo_Equals(t *testing.T) {
	node1 := &Node{ID: "node1"}
	node2 := &Node{ID: "node2"}

	tests := []struct {
		name     string
		p1, p2   *PartitionInfo
		expected bool
	}{
		{
			name: "equal partitions",
			p1: &PartitionInfo{
				Topic:           "topic1",
				Partition:       0,
				Leader:          node1,
				Replicas:        []*Node{node1, node2},
				InSyncReplicas:  []*Node{node1},
				OfflineReplicas: []*Node{node2},
			},
			p2: &PartitionInfo{
				Topic:           "topic1",
				Partition:       0,
				Leader:          node1,
				Replicas:        []*Node{node1, node2},
				InSyncReplicas:  []*Node{node1},
				OfflineReplicas: []*Node{node2},
			},
			expected: true,
		},
		{
			name: "different topic",
			p1: &PartitionInfo{
				Topic:           "topic1",
				Partition:       0,
				Leader:          node1,
				Replicas:        []*Node{node1, node2},
				InSyncReplicas:  []*Node{node1},
				OfflineReplicas: []*Node{node2},
			},
			p2: &PartitionInfo{
				Topic:           "topic2",
				Partition:       0,
				Leader:          node1,
				Replicas:        []*Node{node1, node2},
				InSyncReplicas:  []*Node{node1},
				OfflineReplicas: []*Node{node2},
			},
			expected: false,
		},
		{
			name: "different partition",
			p1: &PartitionInfo{
				Topic:           "topic1",
				Partition:       0,
				Leader:          node1,
				Replicas:        []*Node{node1, node2},
				InSyncReplicas:  []*Node{node1},
				OfflineReplicas: []*Node{node2},
			},
			p2: &PartitionInfo{
				Topic:           "topic1",
				Partition:       1,
				Leader:          node1,
				Replicas:        []*Node{node1, node2},
				InSyncReplicas:  []*Node{node1},
				OfflineReplicas: []*Node{node2},
			},
			expected: false,
		},
		{
			name: "different leader",
			p1: &PartitionInfo{
				Topic:           "topic1",
				Partition:       0,
				Leader:          node1,
				Replicas:        []*Node{node1, node2},
				InSyncReplicas:  []*Node{node1},
				OfflineReplicas: []*Node{node2},
			},
			p2: &PartitionInfo{
				Topic:           "topic1",
				Partition:       0,
				Leader:          node2,
				Replicas:        []*Node{node1, node2},
				InSyncReplicas:  []*Node{node1},
				OfflineReplicas: []*Node{node2},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p1.Equals(tt.p2); got != tt.expected {
				t.Errorf("PartitionInfo.Equals() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestSlicesEqual tests the SlicesEqual function.
func TestSlicesEqual(t *testing.T) {
	node1 := &Node{ID: "node1"}
	node2 := &Node{ID: "node2"}

	tests := []struct {
		a, b     []*Node
		expected bool
	}{
		{
			a:        []*Node{node1, node2},
			b:        []*Node{node1, node2},
			expected: true,
		},
		{
			a:        []*Node{node1, node2},
			b:        []*Node{node2, node1},
			expected: false,
		},
		{
			a:        []*Node{node1},
			b:        []*Node{node1, node2},
			expected: false,
		},
	}

	for _, tt := range tests {
		if got := SlicesEqual(tt.a, tt.b); got != tt.expected {
			t.Errorf("SlicesEqual(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.expected)
		}
	}
}

// TestContains tests the contains function.
func TestContains(t *testing.T) {
	m := map[string]struct{}{
		"node1": {},
		"node2": {},
	}

	tests := []struct {
		key      string
		expected bool
	}{
		{"node1", true},
		{"node2", true},
		{"node3", false},
	}

	for _, tt := range tests {
		if got := contains(m, tt.key); got != tt.expected {
			t.Errorf("contains(%v) = %v, want %v", tt.key, got, tt.expected)
		}
	}
}
