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
	"errors"
	"sync"
	"testing"
	"time"
)

// TestAtomstateFutureImpl_Get tests the Get method of AtomstateFutureImpl.
func TestAtomstateFutureImpl_Get(t *testing.T) {
	f := NewAtomstateFuture[int]()
	err := f.Complete(42)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	value, err := f.Get()
	if err != nil || value != 42 {
		t.Fatalf("expected 42, got %v with error %v", value, err)
	}
}

// TestAtomstateFutureImpl_CompleteExceptionally tests the CompleteExceptionally method.
func TestAtomstateFutureImpl_CompleteExceptionally(t *testing.T) {
	f := NewAtomstateFuture[int]()
	err := f.CompleteExceptionally(errors.New("an error"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	value, err := f.Get()
	if value != 0 || err == nil {
		t.Fatalf("expected error, got %v with error %v", value, err)
	}
}

// TestAtomstateFutureImpl_GetNow tests the GetNow method.
func TestAtomstateFutureImpl_GetNow(t *testing.T) {
	f := NewAtomstateFuture[int]()
	value := f.GetNow(99)
	if value != 99 {
		t.Fatalf("expected 99, got %v", value)
	}
	err := f.Complete(42)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	value = f.GetNow(99)
	if value != 42 {
		t.Fatalf("expected 42, got %v", value)
	}
}

// TestAtomstateFutureImpl_ThenApply tests the ThenApply method.
// TestAtomstateFutureImpl_ThenApply tests the ThenApply method.
func TestAtomstateFutureImpl_ThenApply(t *testing.T) {
	f := NewAtomstateFuture[int]()
	err := f.Complete(42)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	result := f.ThenApply(func(v int) int {
		return v + 1
	})
	// Wait for a short period to ensure the result has time to be computed.
	time.Sleep(100 * time.Millisecond)
	value, err := result.Get()
	if err != nil || value != 43 {
		t.Fatalf("expected 43, got %v with error %v", value, err)
	}
}

// TestAtomstateFutureImpl_ThenApplyAsync tests the ThenApplyAsync method.
func TestAtomstateFutureImpl_ThenApplyAsync(t *testing.T) {
	f := NewAtomstateFuture[int]()
	err := f.Complete(42)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Use a wait group to wait for the asynchronous result
	var wg sync.WaitGroup
	wg.Add(1)

	result := f.ThenApplyAsync(func(v int) int {
		defer wg.Done() // Signal that the async operation is done
		return v + 1
	})

	// Wait for the asynchronous operation to complete
	wg.Wait()

	// Verify the result
	value, err := result.Get()
	if err != nil || value != 43 {
		t.Fatalf("expected 43, got %v with error %v", value, err)
	}
}

// TestAtomstateFutureImpl_WhenComplete tests the WhenComplete method.
func TestAtomstateFutureImpl_WhenComplete(t *testing.T) {
	// Create a new future and complete it
	f := NewAtomstateFuture[int]()
	f.Complete(42)

	// Use a channel to receive the result from WhenComplete
	resultChan := make(chan struct {
		value int
		err   error
	})

	// Call WhenComplete and pass a function that sends the result to the channel
	f.WhenComplete(func(value int, err error) {
		resultChan <- struct {
			value int
			err   error
		}{value, err}
	})

	// Wait for the result and check it
	result := <-resultChan
	if result.value != 42 || result.err != nil {
		t.Fatalf("expected value 42 and no error, got %v with error %v", result.value, result.err)
	}
}

// TestAtomstateFutureImpl_AllOf tests the AllOf method.
func TestAtomstateFutureImpl_AllOf(t *testing.T) {
	// Create futures with type int
	f1 := NewAtomstateFuture[int]()
	f2 := NewAtomstateFuture[int]()

	// Complete futures with different values
	err := f1.Complete(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = f2.Complete(2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Use interface{} to create AllOf future
	allOfFuture := NewAtomstateFuture[interface{}]().AllOf(
		NewAtomstateFuture[interface{}](),
		NewAtomstateFuture[interface{}](),
	)

	// Complete the futures used in AllOf
	err = allOfFuture.Complete(struct{}{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check the result
	_, err = allOfFuture.Get()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestAtomstateFutureImpl_AnyOf tests the AnyOf method.
func TestAtomstateFutureImpl_AnyOf(t *testing.T) {
	f1 := NewAtomstateFuture[int]()
	f2 := NewAtomstateFuture[int]()

	anyOfFuture := NewAtomstateFuture[int]().AnyOf(f1, f2)

	// Complete the first future with a value
	err := f1.Complete(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Create a channel to receive results from the goroutine
	resultChan := make(chan struct {
		value int
		err   error
	})

	go func() {
		value, err := anyOfFuture.Get()
		resultChan <- struct {
			value int
			err   error
		}{value: value, err: err}
	}()

	select {
	case res := <-resultChan:
		if res.err != nil {
			t.Fatalf("expected no error, got %v", res.err)
		}
		if res.value != 1 {
			t.Fatalf("expected value 1, got %v", res.value)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Test timed out")
	}
}

// TestAtomstateFutureImpl_Cancel tests the Cancel method.
func TestAtomstateFutureImpl_Cancel(t *testing.T) {
	f := NewAtomstateFuture[int]()
	cancelled := f.Cancel()
	if !cancelled {
		t.Fatalf("expected future to be cancelled")
	}
	cancelled = f.Cancel()
	if cancelled {
		t.Fatalf("expected future to be already cancelled")
	}
}
