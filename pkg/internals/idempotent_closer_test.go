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

package internals

import (
	"sync"
	"testing"
)

// TestNewIdempotentCloser tests the creation of a new IdempotentCloser.
func TestNewIdempotentCloser(t *testing.T) {
	closer := NewIdempotentCloser()
	if closer.IsClosed() {
		t.Error("Expected IdempotentCloser to be open, but it is closed.")
	}
}

// TestAssertOpen tests the AssertOpen method.
func TestAssertOpen(t *testing.T) {
	closer := NewIdempotentCloser()

	// Should not return an error when the closer is open
	err := closer.AssertOpen("Resource is closed")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Close the resource
	closer.Close(nil, nil)

	// Should return an error when the closer is closed
	err = closer.AssertOpen("Resource is closed")
	if err == nil {
		t.Error("Expected an error, but got none.")
	} else if err.Error() != "Resource is closed" {
		t.Errorf("Expected error message to be 'Resource is closed', but got: %v", err)
	}
}

// TestIsClosed tests the IsClosed method.
func TestIsClosed(t *testing.T) {
	closer := NewIdempotentCloser()

	if closer.IsClosed() {
		t.Error("Expected IdempotentCloser to be open, but it is closed.")
	}

	closer.Close(nil, nil)

	if !closer.IsClosed() {
		t.Error("Expected IdempotentCloser to be closed, but it is open.")
	}
}

// TestClose tests the Close method with initial and subsequent close callbacks.
func TestClose(t *testing.T) {
	var initialClosed, subsequentClosed bool

	closer := NewIdempotentCloser()

	// Close for the first time
	closer.Close(func() {
		initialClosed = true
	}, nil)

	if !initialClosed {
		t.Error("Expected initial close callback to be executed.")
	}

	// Close again
	closer.Close(nil, func() {
		subsequentClosed = true
	})

	if !subsequentClosed {
		t.Error("Expected subsequent close callback to be executed.")
	}
}

// TestConcurrentAccess tests the IdempotentCloser for concurrent access.
func TestConcurrentAccess(t *testing.T) {
	closer := NewIdempotentCloser()
	var wg sync.WaitGroup

	// Start multiple goroutines to test concurrent access
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			closer.AssertOpen("Resource is closed")
		}()
	}

	wg.Wait()

	// Close the resource
	closer.Close(nil, nil)

	// Test concurrent access after closing
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := closer.AssertOpen("Resource is closed")
			if err == nil {
				t.Error("Expected an error after closing, but got none.")
			}
		}()
	}

	wg.Wait()
}

// TestStringMethod tests the String method for correct output.
func TestStringMethod(t *testing.T) {
	closer := NewIdempotentCloser()
	expectedOpenString := "IdempotentCloser{isClosed: false}"

	if closer.String() != expectedOpenString {
		t.Errorf("Expected string: %s, but got: %s", expectedOpenString, closer.String())
	}

	// Close the resource
	closer.Close(nil, nil)
	expectedClosedString := "IdempotentCloser{isClosed: true}"

	if closer.String() != expectedClosedString {
		t.Errorf("Expected string: %s, but got: %s", expectedClosedString, closer.String())
	}
}
