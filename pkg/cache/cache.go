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
	"container/list"
	"sync"
)

// Cache defines the interface for cache operations.
type Cache[K comparable, V any] interface {
	// Get retrieves a value from the cache by its key.
	// Returns the value and a boolean indicating whether the key exists.
	Get(key K) (V, bool)

	// Put adds a key-value pair to the cache.
	// If the key already exists, it updates the value and moves the key to the front.
	Put(key K, value V)

	// Remove deletes a key from the cache.
	// Returns true if the key was found and removed, false otherwise.
	Remove(key K) bool

	// Size returns the number of entries in the cache.
	Size() int
}

// LRUCache implements a Least Recently Used (LRU) cache.
type LRUCache[K comparable, V any] struct {
	maxSize int                 // Maximum size of the cache
	cache   map[K]*list.Element // Map to hold the keys and their corresponding list elements
	order   *list.List          // List to maintain the order of keys based on usage
	values  map[K]V             // Map to hold the actual values
	mu      sync.RWMutex        // Mutex for thread-safe access
}

// NewLRUCache creates a new LRUCache with a specified maximum size.
// It panics if maxSize is less than or equal to zero.
//
// Example usage:
//
//	cache := NewLRUCache[string, int](2)
//	cache.Put("one", 1)
//	value, exists := cache.Get("one") // value = 1, exists = true
func NewLRUCache[K comparable, V any](maxSize int) *LRUCache[K, V] {
	if maxSize <= 0 {
		panic("maxSize must be greater than 0")
	}
	return &LRUCache[K, V]{
		maxSize: maxSize,
		cache:   make(map[K]*list.Element),
		order:   list.New(),
		values:  make(map[K]V),
	}
}

// Get retrieves a value from the cache by its key.
// If the key exists, it moves the key to the front of the usage order.
// Returns the value and a boolean indicating whether the key exists.
//
// Example usage:
//
//	value, exists := cache.Get("one")
//	if exists {
//	    fmt.Println("Value:", value)
//	} else {
//	    fmt.Println("Key not found")
//	}
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if elem, ok := c.cache[key]; ok {
		c.order.MoveToFront(elem)
		return c.values[key], true
	}
	var zero V
	return zero, false
}

// Put adds a key-value pair to the cache.
// If the key already exists, it updates the value and moves the key to the front.
// If the cache exceeds its maximum size, it evicts the least recently used item.
//
// Example usage:
//
//	cache.Put("two", 2)
//	cache.Put("three", 3) // This may evict the least recently used item.
func (c *LRUCache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.order.MoveToFront(elem)
		c.values[key] = value
		return
	}

	if c.order.Len() >= c.maxSize {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			oldKey := oldest.Value.(K)
			delete(c.cache, oldKey)
			delete(c.values, oldKey)
		}
	}

	elem := c.order.PushFront(key)
	c.cache[key] = elem
	c.values[key] = value
}

// Remove deletes a key from the cache.
// Returns true if the key was found and removed, false otherwise.
//
// Example usage:
//
//	removed := cache.Remove("two")
//	if removed {
//	    fmt.Println("Key 'two' removed")
//	} else {
//	    fmt.Println("Key 'two' not found")
//	}
func (c *LRUCache[K, V]) Remove(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.order.Remove(elem)
		delete(c.cache, key)
		delete(c.values, key)
		return true
	}
	return false
}

// Size returns the number of entries in the cache.
//
// Example usage:
//
//	size := cache.Size()
//	fmt.Println("Cache size:", size)
func (c *LRUCache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.values)
}

// SynchronizedCache is a thread-safe wrapper around another Cache.
type SynchronizedCache[K comparable, V any] struct {
	underlying Cache[K, V]
	mu         sync.RWMutex
}

// NewSynchronizedCache creates a new thread-safe cache wrapper.
// It panics if the underlying cache is nil.
//
// Example usage:
//
//	underlying := NewLRUCache[string, int](2)
//	syncCache := NewSynchronizedCache(underlying)
func NewSynchronizedCache[K comparable, V any](underlying Cache[K, V]) *SynchronizedCache[K, V] {
	if underlying == nil {
		panic("underlying cache cannot be nil")
	}
	return &SynchronizedCache[K, V]{underlying: underlying}
}

// Get retrieves a value from the cache by its key.
// It delegates the call to the underlying cache and ensures thread safety.
//
// Example usage:
//
//	value, exists := syncCache.Get("one")
//	if exists {
//	    fmt.Println("Value:", value)
//	} else {
//	    fmt.Println("Key not found")
//	}
func (s *SynchronizedCache[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.underlying.Get(key)
}

// Put adds a key-value pair to the cache.
// It delegates the call to the underlying cache and ensures thread safety.
//
// Example usage:
//
//	syncCache.Put("two", 2)
func (s *SynchronizedCache[K, V]) Put(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.underlying.Put(key, value)
}

// Remove deletes a key from the cache.
// It delegates the call to the underlying cache and ensures thread safety.
// Returns true if the key was found and removed, false otherwise.
//
// Example usage:
//
//	removed := syncCache.Remove("two")
//	if removed {
//	    fmt.Println("Key 'two' removed")
//	} else {
//	    fmt.Println("Key 'two' not found")
//	}
func (s *SynchronizedCache[K, V]) Remove(key K) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.underlying.Remove(key)
}

// Size returns the number of entries in the cache.
// It delegates the call to the underlying cache and ensures thread safety.
//
// Example usage:
//
//	size := syncCache.Size()
//	fmt.Println("Cache size:", size)
func (s *SynchronizedCache[K, V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.underlying.Size()
}
