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
	"testing"
)

// Helper function to check if two buffers are the same object
func buffersAreSame(b1, b2 []byte) bool {
	return &b1[0] == &b2[0]
}

// Test NoCachingBufferSupplier
func TestNoCachingBufferSupplier(t *testing.T) {
	supplier := NewNoCachingBufferSupplier()

	// Test buffer allocation
	buf := supplier.Get(100)
	if len(buf) != 100 {
		t.Fatalf("expected buffer of length 100, got %d", len(buf))
	}

	// Test releasing buffer
	supplier.Release(buf)

	// No caching, so no need to test further
}

// Test DefaultBufferSupplier
func TestDefaultBufferSupplier(t *testing.T) {
	supplier := NewDefaultBufferSupplier()

	// Test buffer allocation and caching
	buf1 := supplier.Get(128)
	buf2 := supplier.Get(128)
	if len(buf1) != 128 {
		t.Fatalf("expected buffer of length 128, got %d", len(buf1))
	}
	if buffersAreSame(buf1, buf2) {
		t.Fatalf("expected different buffers, got same buffer")
	}

	supplier.Release(buf1)
	buf3 := supplier.Get(128)
	if !buffersAreSame(buf1, buf3) {
		t.Fatalf("expected to reuse buffer, got different buffer")
	}

	// Test releasing and closing
	supplier.Close()
}

// Test GrowableBufferSupplier
func TestGrowableBufferSupplier(t *testing.T) {
	supplier := NewGrowableBufferSupplier()

	// Test buffer allocation
	buf1 := supplier.Get(256)
	if len(buf1) != 256 {
		t.Fatalf("expected buffer of length 256, got %d", len(buf1))
	}

	// Test releasing buffer
	supplier.Release(buf1)

	// Test buffer reuse
	buf2 := supplier.Get(128)
	if len(buf2) != 128 {
		t.Fatalf("expected buffer of length 128, got %d", len(buf2))
	}
	if cap(buf1) != cap(buf2) {
		t.Fatalf("expected buffer capacity %d, got %d", cap(buf1), cap(buf2))
	}

	// Test closing
	supplier.Close()
}

// Test concurrent access to DefaultBufferSupplier
func TestDefaultBufferSupplierConcurrency(t *testing.T) {
	supplier := NewDefaultBufferSupplier()
	var wg sync.WaitGroup
	const numGoroutines = 100

	// Concurrently get and release buffers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := supplier.Get(64)
			if len(buf) != 64 {
				t.Errorf("expected buffer of length 64, got %d", len(buf))
			}
			supplier.Release(buf)
		}()
	}

	wg.Wait()
	supplier.Close()
}

// Test concurrent access to GrowableBufferSupplier
func TestGrowableBufferSupplierConcurrency(t *testing.T) {
	supplier := NewGrowableBufferSupplier()
	var wg sync.WaitGroup
	const numGoroutines = 100

	// Concurrently get and release buffers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := supplier.Get(64)
			if len(buf) != 64 {
				t.Errorf("expected buffer of length 64, got %d", len(buf))
			}
			supplier.Release(buf)
		}()
	}

	wg.Wait()
	supplier.Close()
}
