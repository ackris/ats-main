package common

import "testing"

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
