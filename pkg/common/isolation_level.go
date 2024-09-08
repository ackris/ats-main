package common

import "fmt"

// IsolationLevel represents the isolation levels for database transactions.
type IsolationLevel struct {
	id byte
}

// NewIsolationLevel creates a new IsolationLevel with the given id.
func NewIsolationLevel(id byte) (*IsolationLevel, error) {
	if id > 1 {
		return nil, fmt.Errorf("unknown isolation level %d", id)
	}
	return &IsolationLevel{id: id}, nil
}

// ID returns the id of the IsolationLevel.
func (il *IsolationLevel) ID() byte {
	return il.id
}

// String returns the string representation of the isolation level.
func (il *IsolationLevel) String() string {
	names := [...]string{"read_uncommitted", "read_committed"}
	if il.id > 1 {
		return "unknown"
	}
	return names[il.id]
}

// ForID returns the IsolationLevel corresponding to the given byte id.
func ForID(id byte) (*IsolationLevel, error) {
	return NewIsolationLevel(id)
}
