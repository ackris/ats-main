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

import "fmt"

// IsolationLevel represents the isolation levels for message streaming and event storage.
// It defines how transactions interact with messages and events in the system.
type IsolationLevel struct {
	id byte
}

// NewIsolationLevel creates a new IsolationLevel with the given id.
// It returns an error if the provided id is not a valid isolation level.
//
// Valid IDs:
//
//	0 - ReadUncommitted: Allows reading messages that are not yet committed.
//	1 - ReadCommitted: Only allows reading committed messages.
//
// Example:
//
//	il, err := NewIsolationLevel(0)
//	if err != nil {
//	  // Handle error
//	}
//	fmt.Println(il) // Output: &{0}
func NewIsolationLevel(id byte) (*IsolationLevel, error) {
	if id > 1 {
		return nil, fmt.Errorf("unknown isolation level %d", id)
	}
	return &IsolationLevel{id: id}, nil
}

// ID returns the id of the IsolationLevel.
//
// Example:
//
//	il, _ := NewIsolationLevel(1)
//	fmt.Println(il.ID()) // Output: 1
func (il *IsolationLevel) ID() byte {
	return il.id
}

// String returns the string representation of the isolation level.
//
// Example:
//
//	il, _ := NewIsolationLevel(0)
//	fmt.Println(il.String()) // Output: read_uncommitted
func (il *IsolationLevel) String() string {
	names := [...]string{"read_uncommitted", "read_committed"}
	if il.id > 1 {
		return "unknown"
	}
	return names[il.id]
}

// ForID returns the IsolationLevel corresponding to the given byte id.
// It returns an error if the provided id is not a valid isolation level.
//
// Example:
//
//	il, err := ForID(1)
//	if err != nil {
//	  // Handle error
//	}
//	fmt.Println(il) // Output: &{1}
func ForID(id byte) (*IsolationLevel, error) {
	return NewIsolationLevel(id)
}
