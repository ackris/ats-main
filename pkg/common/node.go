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
	"fmt"
	"strconv"
)

// Node represents information about a Atomstate node.
type Node struct {
	// ID is the unique identifier of the node.
	ID int
	// Host is the hostname or IP address of the node.
	Host string
	// Port is the port number of the node.
	Port int
	// RackName is the name of the rack where the node is located.
	RackName *string
}

// NoNode is a placeholder for a node that should be considered as non-existent.
var NoNode = Node{ID: -1, Host: "", Port: -1}

// IsEmpty checks whether the node is empty.
//
// Returns:
//
//	bool: true if the node is empty (Host is empty or Port is less than 0), false otherwise.
//
// Example:
//
//	node := Node{Host: "localhost", Port: 9092}
//	if node.IsEmpty() {
//	    fmt.Println("Node is empty")
//	} else {
//	    fmt.Println("Node is not empty")
//	}
func (n Node) IsEmpty() bool {
	return n.Host == "" || n.Port < 0
}

// IDString returns the string representation of the node ID.
//
// Returns:
//
//	string: the string representation of the node ID.
//
// Example:
//
//	node := Node{ID: 42}
//	idStr := node.IDString()
//	fmt.Println("Node ID:", idStr) // Output: Node ID: 42
func (n Node) IDString() string {
	return strconv.Itoa(n.ID)
}

// HasRack checks if the node has a defined rack.
//
// Returns:
//
//	bool: true if the node has a defined rack (RackName is not nil), false otherwise.
//
// Example:
//
//	rackName := "rack1"
//	node := Node{RackName: &rackName}
//	if node.HasRack() {
//	    fmt.Println("Node has a defined rack")
//	} else {
//	    fmt.Println("Node does not have a defined rack")
//	}
func (n Node) HasRack() bool {
	return n.RackName != nil
}

// Rack returns the rack for this node. It returns an empty string if the rack is not defined.
//
// Returns:
//
//	string: the rack name if defined, otherwise an empty string.
//
// Example:
//
//	rackName := "rack1"
//	node := Node{RackName: &rackName}
//	rack := node.Rack()
//	fmt.Println("Node rack:", rack) // Output: Node rack: rack1
func (n Node) Rack() string {
	if n.RackName != nil {
		return *n.RackName
	}
	return ""
}

// String provides a string representation of the node.
//
// Returns:
//
//	string: a string representation of the node in the format "Host:Port (id: IDString() rack: Rack())".
//
// Example:
//
//	rackName := "rack1"
//	node := Node{ID: 1, Host: "localhost", Port: 9092, RackName: &rackName}
//	nodeStr := node.String()
//	fmt.Println(nodeStr) // Output: localhost:9092 (id: 1 rack: rack1)
func (n Node) String() string {
	rackStr := "none"
	if n.HasRack() {
		rackStr = *n.RackName
	}
	return fmt.Sprintf("%s:%d (id: %s rack: %s)", n.Host, n.Port, n.IDString(), rackStr)
}

// Equal checks if two nodes are equal.
//
// Parameters:
//
//	other (Node): the other node to compare with.
//
// Returns:
//
//	bool: true if the nodes are equal (same ID, Port, Host, and RackName), false otherwise.
//
// Example:
//
//	node1 := Node{ID: 1, Host: "localhost", Port: 9092}
//	node2 := Node{ID: 1, Host: "localhost", Port: 9092}
//	if node1.Equal(node2) {
//	    fmt.Println("Nodes are equal")
//	} else {
//	    fmt.Println("Nodes are not equal")
//	}
func (n Node) Equal(other Node) bool {
	if n.ID != other.ID || n.Port != other.Port || n.Host != other.Host {
		return false
	}
	if n.RackName == nil && other.RackName == nil {
		return true
	}
	if n.RackName != nil && other.RackName != nil {
		return *n.RackName == *other.RackName
	}
	return false
}

// HashCode calculates the hash code for the node.
// This implementation uses a simple hash function that combines the individual components
// of the Node struct using a prime number (31) as the multiplier.
//
// Returns:
//
//	int: the hash code for the node.
//
// Example:
//
//	rackName := "rack1"
//	node := Node{ID: 1, Host: "localhost", Port: 9092, RackName: &rackName}
//	hash := node.HashCode()
//	fmt.Println("Node hash code:", hash)
func (n Node) HashCode() int {
	hash := 1
	hash = hash*31 + n.ID        // Add ID
	hash = hash*31 + len(n.Host) // Add length of Host
	hash = hash*31 + n.Port      // Add Port
	if n.RackName != nil {
		hash = hash*31 + len(*n.RackName) // Add length of RackName
	}
	return hash
}
