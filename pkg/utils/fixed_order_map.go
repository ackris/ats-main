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

package utils

import (
	"errors"
	"sync"
)

// Custom errors for FixedOrderMap operations.
var (
	// ErrKeyExists is returned when attempting to add a key that already exists.
	ErrKeyExists = errors.New("key already exists")

	// ErrRemovalNotAllowed is returned when attempting to remove an entry.
	ErrRemovalNotAllowed = errors.New("removal of entries is not allowed")
)

// FixedOrderMap is a map that maintains the insertion order of keys.
// Once keys are added, they cannot be removed or modified.
// It provides thread-safe operations for setting and getting values,
// while ensuring that the order of insertion is preserved.
type FixedOrderMap struct {
	mu    sync.RWMutex
	data  map[interface{}]interface{}
	order []interface{}
}

// NewFixedOrderMap creates and returns a new instance of FixedOrderMap.
// Example usage:
//
//	m := NewFixedOrderMap()
//	err := m.Set("key1", "value1")
func NewFixedOrderMap() *FixedOrderMap {
	return &FixedOrderMap{
		data:  make(map[interface{}]interface{}),
		order: []interface{}{},
	}
}

// Set adds a new key-value pair to the map.
// It returns an error if the key already exists in the map.
// Example usage:
//
//	err := m.Set("key1", "value1")
//	if err != nil {
//	    log.Fatalf("error setting value: %v", err)
//	}
func (m *FixedOrderMap) Set(key, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[key]; exists {
		return ErrKeyExists
	}

	m.data[key] = value
	m.order = append(m.order, key)
	return nil
}

// Get retrieves the value associated with the given key.
// It returns the value and a boolean indicating whether the key exists in the map.
// Example usage:
//
//	value, exists := m.Get("key1")
//	if exists {
//	    fmt.Println("Value:", value)
//	}
func (m *FixedOrderMap) Get(key interface{}) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.data[key]
	return value, exists
}

// Order returns a slice of keys in the order they were added.
// The slice is a copy, ensuring immutability of the original order slice.
// Example usage:
//
//	order := m.Order()
//	fmt.Println("Keys in order:", order)
func (m *FixedOrderMap) Order() []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy of the order slice for immutability
	orderCopy := make([]interface{}, len(m.order))
	copy(orderCopy, m.order)
	return orderCopy
}

// Remove does nothing and always returns an error indicating removal is not allowed.
// Example usage:
//
//	err := m.Remove("key1")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	}
func (m *FixedOrderMap) Remove(key interface{}) error {
	return ErrRemovalNotAllowed
}

// RemoveEntry does nothing and always returns an error indicating removal is not allowed.
// Example usage:
//
//	err := m.RemoveEntry("key1", "value1")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	}
func (m *FixedOrderMap) RemoveEntry(key, value interface{}) error {
	return ErrRemovalNotAllowed
}

// Clone creates and returns a new FixedOrderMap with the same key-value pairs and insertion order.
// The cloned map is independent of the original map.
// Example usage:
//
//	clone := m.Clone()
//	fmt.Println("Cloned map order:", clone.Order())
func (m *FixedOrderMap) Clone() *FixedOrderMap {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clone := NewFixedOrderMap()
	for _, key := range m.order {
		clone.data[key] = m.data[key]
	}
	// Copy the order slice
	clone.order = make([]interface{}, len(m.order))
	copy(clone.order, m.order)
	return clone
}
