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

func TestFixedOrderMap(t *testing.T) {
	// Create a new FixedOrderMap
	m := NewFixedOrderMap()

	// Test Set and Get
	if err := m.Set("key1", "value1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	value, exists := m.Get("key1")
	if !exists {
		t.Fatal("expected key1 to exist")
	}
	if value != "value1" {
		t.Fatalf("expected value1, got %v", value)
	}

	// Test duplicate key
	if err := m.Set("key1", "newValue"); err != ErrKeyExists {
		t.Fatalf("expected ErrKeyExists, got %v", err)
	}

	// Test Order
	if order := m.Order(); len(order) != 1 || order[0] != "key1" {
		t.Fatalf("expected order [key1], got %v", order)
	}

	// Test Clone
	clone := m.Clone()
	if clone == nil {
		t.Fatal("expected a valid clone")
	}
	if value, exists := clone.Get("key1"); !exists || value != "value1" {
		t.Fatalf("expected clone to have key1 with value1, got %v", value)
	}

	// Ensure clone is independent
	if err := clone.Set("key2", "value2"); err != nil {
		t.Fatalf("expected no error when setting key2 in clone, got %v", err)
	}
	if _, exists := m.Get("key2"); exists {
		t.Fatal("original map should not have key2")
	}
	if order := clone.Order(); len(order) != 2 || order[1] != "key2" {
		t.Fatalf("expected clone order [key1 key2], got %v", order)
	}

	// Test Remove operations
	if err := m.Remove("key1"); err != ErrRemovalNotAllowed {
		t.Fatalf("expected ErrRemovalNotAllowed, got %v", err)
	}
	if err := m.RemoveEntry("key1", "value1"); err != ErrRemovalNotAllowed {
		t.Fatalf("expected ErrRemovalNotAllowed, got %v", err)
	}
}

func BenchmarkFixedOrderMap(b *testing.B) {
	m := NewFixedOrderMap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := m.Set(i, i); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}

	for i := 0; i < b.N; i++ {
		_, exists := m.Get(i)
		if !exists {
			b.Fatalf("key %d should exist", i)
		}
	}
}
