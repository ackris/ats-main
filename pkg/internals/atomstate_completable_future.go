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
	"time"
)

// AtomstateCompletableFuture represents a future value that can be completed only by internal Atomstate clients.
// It provides methods to complete the future with a value or an error, and to retrieve the value asynchronously.
//
// Example Usage:
//
//	future := NewAtomstateCompletableFuture[string]()
//	err := future.Complete("Hello, world!")
//	if err != nil {
//	    // Handle error
//	}
//	value, err := future.Get()
//	if err != nil {
//	    // Handle error
//	} else {
//	    fmt.Println(value) // Output: Hello, world!
//	}
type AtomstateCompletableFuture[T any] struct {
	value    T
	err      error
	complete bool
	mu       sync.Mutex
	done     chan struct{}
}

// NewAtomstateCompletableFuture creates a new instance of AtomstateCompletableFuture.
// It initializes the future with a done channel to notify waiting goroutines when the future is completed.
//
// Example Usage:
//
//	future := NewAtomstateCompletableFuture[int]()
func NewAtomstateCompletableFuture[T any]() *AtomstateCompletableFuture[T] {
	return &AtomstateCompletableFuture[T]{done: make(chan struct{})}
}

// Complete completes the future with the provided value. Only for internal use.
// It sets the value and marks the future as complete, then closes the done channel to notify waiting goroutines.
// If the future is already complete, it returns an error.
//
// Example Usage:
//
//	err := future.Complete("Completed value")
//	if err != nil {
//	    // Handle error
//	}
func (f *AtomstateCompletableFuture[T]) Complete(value T) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.complete {
		return errors.New("future is already complete")
	}

	f.value = value
	f.complete = true
	close(f.done) // Notify all waiting goroutines
	return nil
}

// CompleteExceptionally completes the future with an error. Only for internal use.
// It sets the error and marks the future as complete, then closes the done channel to notify waiting goroutines.
// If the future is already complete, it returns an error.
//
// Example Usage:
//
//	err := future.CompleteExceptionally(errors.New("an error occurred"))
//	if err != nil {
//	    // Handle error
//	}
func (f *AtomstateCompletableFuture[T]) CompleteExceptionally(err error) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.complete {
		return errors.New("future is already complete")
	}

	f.err = err
	f.complete = true
	close(f.done) // Notify all waiting goroutines
	return nil
}

// Get retrieves the value of the future, blocking until it is completed.
// It waits for the done channel to be closed, then returns the value and any associated error.
//
// Example Usage:
//
//	value, err := future.Get()
//	if err != nil {
//	    // Handle error
//	} else {
//	    fmt.Println(value) // Use the retrieved value
//	}
func (f *AtomstateCompletableFuture[T]) Get() (T, error) {
	<-f.done // Wait until the future is completed

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return f.value, f.err
	}
	return f.value, nil
}

// CompleteAsync allows for asynchronous completion of the future with a value.
// It starts a new goroutine that sleeps for the specified timeout, then completes the future with the value returned by the supplier function.
//
// Example Usage:
//
//	future.CompleteAsync(func() string {
//	    return "Async value"
//	}, 2*time.Second)
func (f *AtomstateCompletableFuture[T]) CompleteAsync(supplier func() T, timeout time.Duration) error {
	go func() {
		time.Sleep(timeout)
		f.Complete(supplier())
	}()
	return nil
}

// CompleteOnTimeout completes the future with a value after a timeout.
// It uses time.AfterFunc to schedule a function that completes the future with the specified value after the timeout.
//
// Example Usage:
//
//	future.CompleteOnTimeout("Timeout value", 5*time.Second)
func (f *AtomstateCompletableFuture[T]) CompleteOnTimeout(value T, timeout time.Duration) {
	time.AfterFunc(timeout, func() {
		f.Complete(value)
	})
}

// ObtrudeValue forcibly sets the value of the future. Not allowed for user code.
// This method is intended for internal use only and will return an error if called from user code.
//
// Example Usage (not recommended):
//
//	err := future.ObtrudeValue("Forced value") // This will return an error
func (f *AtomstateCompletableFuture[T]) ObtrudeValue(value T) error {
	return f.ErroneousCompletionException()
}

// ObtrudeException forcibly sets the error of the future. Not allowed for user code.
// This method is intended for internal use only and will return an error if called from user code.
//
// Example Usage (not recommended):
//
//	err := future.ObtrudeException(errors.New("Forced error")) // This will return an error
func (f *AtomstateCompletableFuture[T]) ObtrudeException(err error) error {
	return f.ErroneousCompletionException()
}

// ErroneousCompletionException returns an error indicating that user code should not complete futures.
// This method is intended for internal use only.
func (f *AtomstateCompletableFuture[T]) ErroneousCompletionException() error {
	return errors.New("user code should not complete futures returned from Atomstate clients")
}
