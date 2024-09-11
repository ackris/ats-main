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

// BufferSupplier is an interface for supplying and managing byte buffers.
// It allows for retrieving and releasing buffers with specified capacities,
// and should handle buffer reuse and resource cleanup.
type BufferSupplier interface {
	// Get returns a buffer with at least the specified capacity.
	// The buffer may be newly allocated or retrieved from a cache.
	// Example:
	//     supplier := NewDefaultBufferSupplier()
	//     buf := supplier.Get(1024) // Get a buffer of at least 1024 bytes
	//     // Use the buffer
	//     supplier.Release(buf)    // Release the buffer for future reuse
	Get(capacity int) []byte

	// Release returns a buffer to the supplier for future reuse.
	// The buffer must be reset (e.g., with clear or reset methods) before release.
	// Example:
	//     supplier.Release(buf) // Return the buffer to the supplier
	Release(buffer []byte)

	// Close releases any resources associated with the supplier.
	// This should be called when the supplier is no longer needed.
	// Example:
	//     supplier.Close() // Clean up resources associated with the supplier
	Close()
}

// NoCachingBufferSupplier is a buffer supplier that does not cache buffers.
// It always allocates a new buffer when Get is called.
type NoCachingBufferSupplier struct{}

// NewNoCachingBufferSupplier creates a new instance of NoCachingBufferSupplier.
func NewNoCachingBufferSupplier() *NoCachingBufferSupplier {
	return &NoCachingBufferSupplier{}
}

// Get returns a new buffer of the specified capacity.
// Example:
//
//	supplier := NewNoCachingBufferSupplier()
//	buf := supplier.Get(1024) // Allocate a new buffer of 1024 bytes
func (n *NoCachingBufferSupplier) Get(capacity int) []byte {
	return make([]byte, capacity)
}

// Release does nothing as this supplier does not cache buffers.
// Example:
//
//	supplier.Release(buf) // No caching, so this method does nothing
func (n *NoCachingBufferSupplier) Release(buffer []byte) {}

// Close releases resources associated with the supplier (no-op here).
// Example:
//
//	supplier.Close() // No resources to clean up
func (n *NoCachingBufferSupplier) Close() {}

// DefaultBufferSupplier is a buffer supplier that caches buffers by size.
// It maintains a cache of buffers of each size and reuses them if available.
type DefaultBufferSupplier struct {
	cache sync.Map
}

// NewDefaultBufferSupplier creates a new instance of DefaultBufferSupplier.
func NewDefaultBufferSupplier() *DefaultBufferSupplier {
	return &DefaultBufferSupplier{}
}

// Get retrieves a buffer with at least the specified capacity.
// It reuses a buffer from the cache if available; otherwise, it allocates a new buffer.
// Example:
//
//	supplier := NewDefaultBufferSupplier()
//	buf := supplier.Get(2048) // Get a buffer of at least 2048 bytes
//	// Use the buffer
//	supplier.Release(buf)    // Release the buffer to cache it for future use
func (d *DefaultBufferSupplier) Get(capacity int) []byte {
	if value, ok := d.cache.Load(capacity); ok {
		if buf, ok := value.(*[]byte); ok && buf != nil {
			d.cache.Delete(capacity)
			return *buf
		}
	}
	return make([]byte, capacity)
}

// Release returns a buffer to the cache for future reuse.
// The buffer must be reset before it is released.
// Example:
//
//	supplier.Release(buf) // Return the buffer to the cache
func (d *DefaultBufferSupplier) Release(buffer []byte) {
	capacity := len(buffer)
	buf := buffer
	d.cache.Store(capacity, &buf)
}

// Close clears the buffer cache.
// Example:
//
//	supplier.Close() // Clear the cache and clean up resources
func (d *DefaultBufferSupplier) Close() {
	d.cache = sync.Map{}
}

// GrowableBufferSupplier is a buffer supplier that caches a single buffer.
// It reuses this buffer, growing it as needed, and releases it when done.
type GrowableBufferSupplier struct {
	mu           sync.Mutex
	cachedBuffer []byte
}

// NewGrowableBufferSupplier creates a new instance of GrowableBufferSupplier.
func NewGrowableBufferSupplier() *GrowableBufferSupplier {
	return &GrowableBufferSupplier{}
}

// Get retrieves a buffer with at least the specified capacity.
// It reuses a previously cached buffer if it is large enough, otherwise allocates a new one.
// Example:
//
//	supplier := NewGrowableBufferSupplier()
//	buf := supplier.Get(4096) // Get a buffer of at least 4096 bytes
//	// Use the buffer
//	supplier.Release(buf)    // Release the buffer to cache it for future use
func (g *GrowableBufferSupplier) Get(minCapacity int) []byte {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.cachedBuffer != nil && cap(g.cachedBuffer) >= minCapacity {
		buf := g.cachedBuffer[:minCapacity]
		g.cachedBuffer = nil
		return buf
	}
	return make([]byte, minCapacity)
}

// Release returns a buffer to be cached for future reuse.
// The buffer must be reset before release.
// Example:
//
//	supplier.Release(buf) // Return the buffer to be cached
func (g *GrowableBufferSupplier) Release(buffer []byte) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cachedBuffer = buffer
}

// Close clears the cached buffer.
// Example:
//
//	supplier.Close() // Clear the cached buffer and clean up resources
func (g *GrowableBufferSupplier) Close() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cachedBuffer = nil
}
