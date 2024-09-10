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

package cache

import (
	"testing"
)

func TestLRUCache(t *testing.T) {
	cache := NewLRUCache[string, int](2)

	// Test Put and Get
	cache.Put("one", 1)
	if val, ok := cache.Get("one"); !ok || val != 1 {
		t.Errorf("Expected 1, got %v", val)
	}

	// Test Size
	if size := cache.Size(); size != 1 {
		t.Errorf("Expected size 1, got %d", size)
	}

	// Test eviction of least recently used item
	cache.Put("two", 2)
	cache.Put("three", 3) // This should evict "one"
	if _, ok := cache.Get("one"); ok {
		t.Error("Expected 'one' to be evicted")
	}

	// Test that "two" is still present
	if val, ok := cache.Get("two"); !ok || val != 2 {
		t.Errorf("Expected 2, got %v", val)
	}

	// Test updating an existing key
	cache.Put("two", 22)
	if val, ok := cache.Get("two"); !ok || val != 22 {
		t.Errorf("Expected 22, got %v", val)
	}

	// Test Remove
	if removed := cache.Remove("two"); !removed {
		t.Error("Expected 'two' to be removed")
	}
	if _, ok := cache.Get("two"); ok {
		t.Error("Expected 'two' to be absent after removal")
	}

	// Test Size after removal
	if size := cache.Size(); size != 1 {
		t.Errorf("Expected size 1 after removal, got %d", size)
	}

	// Test eviction with multiple items
	cache.Put("four", 4)
	cache.Put("five", 5) // This should evict "three"
	if _, ok := cache.Get("three"); ok {
		t.Error("Expected 'three' to be evicted")
	}
}

func TestSynchronizedCache(t *testing.T) {
	underlyingCache := NewLRUCache[string, int](2)
	syncCache := NewSynchronizedCache(underlyingCache)

	// Test Put and Get
	syncCache.Put("one", 1)
	if val, ok := syncCache.Get("one"); !ok || val != 1 {
		t.Errorf("Expected 1, got %v", val)
	}

	// Test Size
	if size := syncCache.Size(); size != 1 {
		t.Errorf("Expected size 1, got %d", size)
	}

	// Test eviction of least recently used item
	syncCache.Put("two", 2)
	syncCache.Put("three", 3) // This should evict "one"
	if _, ok := syncCache.Get("one"); ok {
		t.Error("Expected 'one' to be evicted")
	}

	// Test that "two" is still present
	if val, ok := syncCache.Get("two"); !ok || val != 2 {
		t.Errorf("Expected 2, got %v", val)
	}

	// Test updating an existing key
	syncCache.Put("two", 22)
	if val, ok := syncCache.Get("two"); !ok || val != 22 {
		t.Errorf("Expected 22, got %v", val)
	}

	// Test Remove
	if removed := syncCache.Remove("two"); !removed {
		t.Error("Expected 'two' to be removed")
	}
	if _, ok := syncCache.Get("two"); ok {
		t.Error("Expected 'two' to be absent after removal")
	}

	// Test Size after removal
	if size := syncCache.Size(); size != 1 {
		t.Errorf("Expected size 1 after removal, got %d", size)
	}

	// Test eviction with multiple items
	syncCache.Put("four", 4)
	syncCache.Put("five", 5) // This should evict "three"
	if _, ok := syncCache.Get("three"); ok {
		t.Error("Expected 'three' to be evicted")
	}
}
