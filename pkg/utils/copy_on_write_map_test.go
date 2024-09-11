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
	"testing"
)

// TestNew tests the creation of a new CopyOnWriteMap instance.
func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("Expected a new instance of CopyOnWriteMap, got nil")
	}
	if m.Size() != 0 {
		t.Fatalf("Expected size 0, got %d", m.Size())
	}
}

// TestPutAndGet tests putting and getting values.
func TestPutAndGet(t *testing.T) {
	m := New()
	m.Put("key1", "value1")

	value, exists := m.Get("key1")
	if !exists {
		t.Fatal("Expected key1 to exist")
	}
	if value != "value1" {
		t.Fatalf("Expected value1, got %v", value)
	}
}

// TestContainsKey tests checking if a key exists.
func TestContainsKey(t *testing.T) {
	m := New()
	m.Put("key1", "value1")

	if !m.ContainsKey("key1") {
		t.Fatal("Expected key1 to be present")
	}
	if m.ContainsKey("key2") {
		t.Fatal("Expected key2 to be absent")
	}
}

// TestContainsValue tests checking if a value exists.
func TestContainsValue(t *testing.T) {
	m := New()
	m.Put("key1", "value1")

	if !m.ContainsValue("value1") {
		t.Fatal("Expected value1 to be present")
	}
	if m.ContainsValue("value2") {
		t.Fatal("Expected value2 to be absent")
	}
}

// TestPutIfAbsent tests putting a key-value pair only if the key is absent.
func TestPutIfAbsent(t *testing.T) {
	m := New()
	m.PutIfAbsent("key1", "value1")

	if value := m.PutIfAbsent("key1", "newValue"); value != "value1" {
		t.Fatalf("Expected old value to be value1, got %v", value)
	}
	if value := m.PutIfAbsent("key2", "value2"); value != nil {
		t.Fatalf("Expected nil as key2 was absent, got %v", value)
	}
}

// TestRemove tests removing a key-value pair.
func TestRemove(t *testing.T) {
	m := New()
	m.Put("key1", "value1")

	oldValue := m.Remove("key1")
	if oldValue != "value1" {
		t.Fatalf("Expected removed value to be value1, got %v", oldValue)
	}
	if m.ContainsKey("key1") {
		t.Fatal("Expected key1 to be absent")
	}
}

// TestRemoveIfValueMatches tests removing a key if the value matches.
func TestRemoveIfValueMatches(t *testing.T) {
	m := New()
	m.Put("key1", "value1")

	if !m.RemoveIfValueMatches("key1", "value1") {
		t.Fatal("Expected RemoveIfValueMatches to return true")
	}
	if m.ContainsKey("key1") {
		t.Fatal("Expected key1 to be absent")
	}
	if m.RemoveIfValueMatches("key1", "value1") {
		t.Fatal("Expected RemoveIfValueMatches to return false when key is absent")
	}
}

// TestReplace tests replacing a value if the key and old value match.
func TestReplace(t *testing.T) {
	m := New()
	m.Put("key1", "value1")

	if !m.Replace("key1", "value1", "newValue") {
		t.Fatal("Expected Replace to return true when the value is replaced")
	}

	if value, exists := m.Get("key1"); !exists || value != "newValue" {
		t.Fatalf("Expected newValue, got %v", value)
	}
}

// TestReplaceValue tests replacing a value for an existing key.
func TestReplaceValue(t *testing.T) {
	m := New()
	m.Put("key1", "value1")

	oldValue := m.ReplaceValue("key1", "newValue")
	if oldValue != "value1" {
		t.Fatalf("Expected old value to be value1, got %v", oldValue)
	}
	if value, exists := m.Get("key1"); !exists || value != "newValue" {
		t.Fatalf("Expected newValue, got %v", value)
	}
}

// TestClear tests clearing the map.
func TestClear(t *testing.T) {
	m := New()
	m.Put("key1", "value1")
	m.Clear()

	if !m.IsEmpty() {
		t.Fatal("Expected map to be empty after Clear")
	}
}

// TestSize tests the Size method.
func TestSize(t *testing.T) {
	m := New()
	if m.Size() != 0 {
		t.Fatalf("Expected size 0, got %d", m.Size())
	}

	m.Put("key1", "value1")
	if m.Size() != 1 {
		t.Fatalf("Expected size 1, got %d", m.Size())
	}

	m.Put("key2", "value2")
	if m.Size() != 2 {
		t.Fatalf("Expected size 2, got %d", m.Size())
	}

	m.Remove("key1")
	if m.Size() != 1 {
		t.Fatalf("Expected size 1, got %d", m.Size())
	}
}

// TestKeySet tests the KeySet method.
func TestKeySet(t *testing.T) {
	m := New()
	m.Put("key1", "value1")
	m.Put("key2", "value2")

	keys := m.KeySet()
	if len(keys) != 2 {
		t.Fatalf("Expected 2 keys, got %d", len(keys))
	}
	if !contains(keys, "key1") || !contains(keys, "key2") {
		t.Fatalf("Expected keySet to contain key1 and key2")
	}
}

// TestValues tests the Values method.
func TestValues(t *testing.T) {
	m := New()
	m.Put("key1", "value1")
	m.Put("key2", "value2")

	values := m.Values()
	if len(values) != 2 {
		t.Fatalf("Expected 2 values, got %d", len(values))
	}
	if !containsValue(values, "value1") || !containsValue(values, "value2") {
		t.Fatalf("Expected values to contain value1 and value2")
	}
}

// Helper function to check if a slice contains a specific key.
func contains(slice []interface{}, item interface{}) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Helper function to check if a slice contains a specific value.
func containsValue(slice []interface{}, item interface{}) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
