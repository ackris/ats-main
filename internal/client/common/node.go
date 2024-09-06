package common

import (
	"fmt"
	"strconv"
)

// Node represents information about a Kafka node.
type Node struct {
	ID       int
	Host     string
	Port     int
	RackName *string
}

// NoNode is a placeholder for a node that should be considered as non-existent.
var NoNode = Node{ID: -1, Host: "", Port: -1}

// IsEmpty checks whether the node is empty.
func (n Node) IsEmpty() bool {
	return n.Host == "" || n.Port < 0
}

// IDString returns the string representation of the node ID.
func (n Node) IDString() string {
	return strconv.Itoa(n.ID)
}

// HasRack checks if the node has a defined rack.
func (n Node) HasRack() bool {
	return n.RackName != nil
}

// Rack returns the rack for this node. It returns an empty string if the rack is not defined.
func (n Node) Rack() string {
	if n.RackName != nil {
		return *n.RackName
	}
	return ""
}

// String provides a string representation of the node.
func (n Node) String() string {
	rackStr := "none"
	if n.HasRack() {
		rackStr = *n.RackName
	}
	return fmt.Sprintf("%s:%d (id: %s rack: %s)", n.Host, n.Port, n.IDString(), rackStr)
}

// Equal checks if two nodes are equal.
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
func (n Node) HashCode() int {
	hash := 31
	hash = hash*31 + len(n.Host)
	hash = hash*31 + n.ID
	hash = hash*31 + n.Port
	if n.RackName != nil {
		hash = hash*31 + len(*n.RackName)
	}
	return hash
}
