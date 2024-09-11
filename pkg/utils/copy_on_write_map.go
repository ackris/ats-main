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
	"sync"
)

// CopyOnWriteMap is a thread-safe map optimized for read-heavy workloads.
// It uses a read-write lock to synchronize access and ensure that
// read operations are efficient while write operations are safe.
type CopyOnWriteMap struct {
	mu   sync.RWMutex
	data map[interface{}]interface{}
}

// New creates and returns a new instance of CopyOnWriteMap.
// The map is initially empty.
func New() *CopyOnWriteMap {
	return &CopyOnWriteMap{
		data: make(map[interface{}]interface{}),
	}
}

// ContainsKey checks if the map contains the specified key.
//
// Usage Example:
//
//	m := New()
//	m.Put("key1", "value1")
//	exists := m.ContainsKey("key1")
//	fmt.Println(exists) // Output: true
func (m *CopyOnWriteMap) ContainsKey(key interface{}) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.data[key]
	return exists
}

// ContainsValue checks if the map contains the specified value.
//
// Usage Example:
//
//	m := New()
//	m.Put("key1", "value1")
//	exists := m.ContainsValue("value1")
//	fmt.Println(exists) // Output: true
func (m *CopyOnWriteMap) ContainsValue(value interface{}) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range m.data {
		if v == value {
			return true
		}
	}
	return false
}

// Get retrieves the value associated with the specified key.
//
// Usage Example:
//
//	m := New()
//	m.Put("key1", "value1")
//	value, exists := m.Get("key1")
//	if exists {
//	    fmt.Println(value) // Output: value1
//	}
func (m *CopyOnWriteMap) Get(key interface{}) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.data[key]
	return value, exists
}

// IsEmpty checks if the map is empty.
//
// Usage Example:
//
//	m := New()
//	fmt.Println(m.IsEmpty()) // Output: true
//	m.Put("key1", "value1")
//	fmt.Println(m.IsEmpty()) // Output: false
func (m *CopyOnWriteMap) IsEmpty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data) == 0
}

// KeySet returns a slice of keys in the map.
//
// Usage Example:
//
//	m := New()
//	m.Put("key1", "value1")
//	m.Put("key2", "value2")
//	keys := m.KeySet()
//	fmt.Println(keys) // Output: [key1 key2]
func (m *CopyOnWriteMap) KeySet() []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]interface{}, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys
}

// Size returns the number of key-value pairs in the map.
//
// Usage Example:
//
//	m := New()
//	fmt.Println(m.Size()) // Output: 0
//	m.Put("key1", "value1")
//	fmt.Println(m.Size()) // Output: 1
func (m *CopyOnWriteMap) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Values returns a slice of values in the map.
//
// Usage Example:
//
//	m := New()
//	m.Put("key1", "value1")
//	m.Put("key2", "value2")
//	values := m.Values()
//	fmt.Println(values) // Output: [value1 value2]
func (m *CopyOnWriteMap) Values() []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	values := make([]interface{}, 0, len(m.data))
	for _, value := range m.data {
		values = append(values, value)
	}
	return values
}

// Clear removes all entries from the map.
//
// Usage Example:
//
//	m := New()
//	m.Put("key1", "value1")
//	m.Clear()
//	fmt.Println(m.IsEmpty()) // Output: true
func (m *CopyOnWriteMap) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[interface{}]interface{})
}

// Put adds or updates the key-value pair in the map and returns the previous value
// associated with the key, or nil if there was no mapping for the key.
//
// Usage Example:
//
//	m := New()
//	oldValue := m.Put("key1", "value1")
//	fmt.Println(oldValue) // Output: nil
//	oldValue = m.Put("key1", "newValue")
//	fmt.Println(oldValue) // Output: value1
func (m *CopyOnWriteMap) Put(key, value interface{}) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	oldValue := m.data[key]
	m.data[key] = value
	return oldValue
}

// PutAll adds or updates multiple key-value pairs in the map.
//
// Usage Example:
//
//	m := New()
//	entries := map[interface{}]interface{}{
//	    "key1": "value1",
//	    "key2": "value2",
//	}
//	m.PutAll(entries)
//	fmt.Println(m.Size()) // Output: 2
func (m *CopyOnWriteMap) PutAll(entries map[interface{}]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, value := range entries {
		m.data[key] = value
	}
}

// Remove deletes the key-value pair associated with the specified key and
// returns the previous value associated with the key, or nil if there was no mapping for the key.
//
// Usage Example:
//
//	m := New()
//	m.Put("key1", "value1")
//	oldValue := m.Remove("key1")
//	fmt.Println(oldValue) // Output: value1
//	oldValue = m.Remove("key2")
//	fmt.Println(oldValue) // Output: nil
func (m *CopyOnWriteMap) Remove(key interface{}) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	oldValue := m.data[key]
	delete(m.data, key)
	return oldValue
}

// PutIfAbsent adds a key-value pair if the key is not already present and
// returns the current value associated with the key, or nil if the key was absent.
//
// Usage Example:
//
//	m := New()
//	oldValue := m.PutIfAbsent("key1", "value1")
//	fmt.Println(oldValue) // Output: nil
//	oldValue = m.PutIfAbsent("key1", "newValue")
//	fmt.Println(oldValue) // Output: value1
func (m *CopyOnWriteMap) PutIfAbsent(key, value interface{}) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.data[key]; !exists {
		m.data[key] = value
		return nil
	}
	return m.data[key]
}

// RemoveIfValueMatches removes the key if the value matches the specified value
// and returns true if the key was removed, false otherwise.
//
// Usage Example:
//
//	m := New()
//	m.Put("key1", "value1")
//	removed := m.RemoveIfValueMatches("key1", "value1")
//	fmt.Println(removed) // Output: true
//	removed = m.RemoveIfValueMatches("key1", "value2")
//	fmt.Println(removed) // Output: false
func (m *CopyOnWriteMap) RemoveIfValueMatches(key, value interface{}) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, exists := m.data[key]; exists && v == value {
		delete(m.data, key)
		return true
	}
	return false
}

// Replace updates the value associated with the key if the key exists and the original value matches.
// Returns true if the value was replaced, false otherwise.
//
// Usage Example:
//
//	m := New()
//	m.Put("key1", "value1")
//	replaced := m.Replace("key1", "value1", "newValue")
//	fmt.Println(replaced) // Output: true
//	replaced = m.Replace("key1", "value1", "anotherValue")
//	fmt.Println(replaced) // Output: false
func (m *CopyOnWriteMap) Replace(key, original, replacement interface{}) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, exists := m.data[key]; exists && v == original {
		m.data[key] = replacement
		return true
	}
	return false
}

// ReplaceValue updates the value associated with the key if the key exists.
// Returns the previous value associated with the key, or nil if the key was not present.
//
// Usage Example:
//
//	m := New()
//	m.Put("key1", "value1")
//	oldValue := m.ReplaceValue("key1", "newValue")
//	fmt.Println(oldValue) // Output: value1
func (m *CopyOnWriteMap) ReplaceValue(key, value interface{}) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, exists := m.data[key]; exists {
		m.data[key] = value
		return v
	}
	return nil
}
