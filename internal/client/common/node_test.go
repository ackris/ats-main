package common

import (
	"testing"
)

func TestNode_IsEmpty(t *testing.T) {
	tests := []struct {
		node     Node
		expected bool
	}{
		{Node{ID: 1, Host: "localhost", Port: 9092, RackName: nil}, false},
		{Node{ID: 1, Host: "", Port: 9092, RackName: nil}, true},
		{Node{ID: 1, Host: "localhost", Port: -1, RackName: nil}, true},
		{NoNode, true},
	}

	for _, test := range tests {
		result := test.node.IsEmpty()
		if result != test.expected {
			t.Errorf("IsEmpty() for node %v = %v; expected %v", test.node, result, test.expected)
		}
	}
}

func TestNode_IDString(t *testing.T) {
	node := Node{ID: 42, Host: "localhost", Port: 9092, RackName: nil}
	expected := "42"
	result := node.IDString()
	if result != expected {
		t.Errorf("IDString() = %v; expected %v", result, expected)
	}
}

func TestNode_HasRack(t *testing.T) {
	tests := []struct {
		node     Node
		expected bool
	}{
		{Node{ID: 1, Host: "localhost", Port: 9092, RackName: nil}, false},
		{Node{ID: 2, Host: "localhost", Port: 9092, RackName: new(string)}, true},
	}

	for _, test := range tests {
		result := test.node.HasRack()
		if result != test.expected {
			t.Errorf("HasRack() for node %v = %v; expected %v", test.node, result, test.expected)
		}
	}
}

func TestNode_Rack(t *testing.T) {
	rackName := "rack1"
	tests := []struct {
		node     Node
		expected string
	}{
		{Node{ID: 1, Host: "localhost", Port: 9092, RackName: nil}, ""},
		{Node{ID: 2, Host: "localhost", Port: 9092, RackName: &rackName}, "rack1"},
	}

	for _, test := range tests {
		result := test.node.Rack()
		if result != test.expected {
			t.Errorf("Rack() for node %v = %v; expected %v", test.node, result, test.expected)
		}
	}
}

func TestNode_String(t *testing.T) {
	rackName := "rack1"
	tests := []struct {
		node     Node
		expected string
	}{
		{Node{ID: 1, Host: "localhost", Port: 9092, RackName: nil}, "localhost:9092 (id: 1 rack: none)"},
		{Node{ID: 2, Host: "localhost", Port: 9093, RackName: &rackName}, "localhost:9093 (id: 2 rack: rack1)"},
	}

	for _, test := range tests {
		result := test.node.String()
		if result != test.expected {
			t.Errorf("String() for node %v = %v; expected %v", test.node, result, test.expected)
		}
	}
}

func TestNode_Equal(t *testing.T) {
	rackName1 := "rack1"
	rackName2 := "rack2"

	tests := []struct {
		node1    Node
		node2    Node
		expected bool
	}{
		{Node{ID: 1, Host: "localhost", Port: 9092, RackName: nil}, Node{ID: 1, Host: "localhost", Port: 9092, RackName: nil}, true},
		{Node{ID: 1, Host: "localhost", Port: 9092, RackName: nil}, Node{ID: 1, Host: "localhost", Port: 9093, RackName: nil}, false},
		{Node{ID: 1, Host: "localhost", Port: 9092, RackName: &rackName1}, Node{ID: 1, Host: "localhost", Port: 9092, RackName: &rackName1}, true},
		{Node{ID: 1, Host: "localhost", Port: 9092, RackName: &rackName1}, Node{ID: 1, Host: "localhost", Port: 9092, RackName: &rackName2}, false},
		{Node{ID: 1, Host: "localhost", Port: 9092, RackName: nil}, Node{ID: 1, Host: "localhost", Port: 9092, RackName: &rackName1}, false},
	}

	for _, test := range tests {
		result := test.node1.Equal(test.node2)
		if result != test.expected {
			t.Errorf("Equal() for nodes %v and %v = %v; expected %v", test.node1, test.node2, result, test.expected)
		}
	}
}

func TestNode_HashCode(t *testing.T) {
	rackName := "rack1"
	node := Node{ID: 1, Host: "localhost", Port: 9092, RackName: &rackName}
	expectedHash := 1243818 // Update this value based on the correct hash calculation

	result := node.HashCode()
	if result != expectedHash {
		t.Errorf("HashCode() for node %v = %v; expected %v", node, result, expectedHash)
	}
}
