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
	"fmt"
	"sync"
)

// AtomstateFuture represents the interface for a future with various async methods.
type AtomstateFuture[T any] interface {
	Complete(value T) error
	CompleteExceptionally(err error) error
	Get() (T, error)
	GetNow(defaultValue T) T
	Accept(action func(T))
	ToCompletionStage() *AtomstateFutureImpl[T]
	ToString() string
	ThenApply(function func(T) T) *AtomstateFutureImpl[T]
	ThenApplyAsync(function func(T) T) *AtomstateFutureImpl[T]
	WhenComplete(action func(T, error))
	WhenCompleteAsync(action func(T, error))
	Handle(biFunction func(T, error) T) *AtomstateFutureImpl[T]
	HandleAsync(biFunction func(T, error) T) *AtomstateFutureImpl[T]
	Exceptionally(function func(error) T) *AtomstateFutureImpl[T]
	AllOf(futures ...*AtomstateFutureImpl[T]) *AtomstateFutureImpl[struct{}]
	AnyOf(futures ...*AtomstateFutureImpl[T]) *AtomstateFutureImpl[T]
	Cancel() bool
	IsCancelled() bool
	IsDone() bool
}

// AtomstateFutureImpl implements AtomstateFuture interface.
type AtomstateFutureImpl[T any] struct {
	completableFuture *AtomstateCompletableFuture[T]
	cancelled         bool
	mu                sync.RWMutex
}

// NewAtomstateFuture creates a new AtomstateFutureImpl instance.
func NewAtomstateFuture[T any]() *AtomstateFutureImpl[T] {
	return &AtomstateFutureImpl[T]{
		completableFuture: NewAtomstateCompletableFuture[T](),
	}
}

// Complete completes the future with the provided value.
func (f *AtomstateFutureImpl[T]) Complete(value T) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.cancelled {
		return errors.New("future is cancelled")
	}
	return f.completableFuture.Complete(value)
}

// CompleteExceptionally completes the future with an error.
func (f *AtomstateFutureImpl[T]) CompleteExceptionally(err error) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.cancelled {
		return errors.New("future is cancelled")
	}
	return f.completableFuture.CompleteExceptionally(err)
}

// Get retrieves the value of the future, blocking until it is completed.
func (f *AtomstateFutureImpl[T]) Get() (T, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.cancelled {
		return *new(T), errors.New("future is cancelled")
	}
	return f.completableFuture.Get()
}

// GetNow retrieves the value of the future if it's completed, otherwise returns the default value.
func (f *AtomstateFutureImpl[T]) GetNow(defaultValue T) T {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.cancelled {
		return defaultValue
	}

	// Check if the future is done before calling Get
	if !f.IsDone() {
		return defaultValue
	}

	// Call Get to retrieve the value if it is done
	value, err := f.completableFuture.Get()
	if err != nil {
		return defaultValue
	}
	return value
}

// Accept runs the provided action with the value of this future.
func (f *AtomstateFutureImpl[T]) Accept(action func(T)) {
	go func() {
		value, err := f.completableFuture.Get()
		if err == nil {
			action(value)
		}
	}()
}

// ToCompletionStage returns the AtomstateFutureImpl itself.
func (f *AtomstateFutureImpl[T]) ToCompletionStage() *AtomstateFutureImpl[T] {
	return f
}

// ToString returns a string representation of the future's state.
func (f *AtomstateFutureImpl[T]) ToString() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.cancelled {
		return "Cancelled"
	}
	if f.IsDone() {
		value, err := f.completableFuture.Get()
		if err != nil {
			return fmt.Sprintf("CompletedExceptionally[%v]", err)
		}
		return fmt.Sprintf("Completed[%v]", value)
	}
	return "Incomplete"
}

// ThenApply returns a new AtomstateFuture that applies the provided function to the result of this future.
func (f *AtomstateFutureImpl[T]) ThenApply(function func(T) T) *AtomstateFutureImpl[T] {
	result := NewAtomstateFuture[T]()
	go func() {
		// Wait for the future to complete and then apply the function.
		value, err := f.Get()
		if err != nil {
			result.CompleteExceptionally(err)
			return
		}
		result.Complete(function(value))
	}()
	return result
}

// ThenApplyAsync returns a new AtomstateFuture that applies the provided function asynchronously.
func (f *AtomstateFutureImpl[T]) ThenApplyAsync(function func(T) T) *AtomstateFutureImpl[T] {
	result := NewAtomstateFuture[T]()
	go func() {
		value, err := f.Get() // Get the value from the future
		if err != nil {
			result.CompleteExceptionally(err)
			return
		}
		// Apply the function and complete the result future
		result.Complete(function(value))
	}()
	return result
}

// WhenComplete runs the provided action when this future completes.
func (f *AtomstateFutureImpl[T]) WhenComplete(action func(T, error)) {
	go func() {
		value, err := f.completableFuture.Get()
		action(value, err)
	}()
}

// WhenCompleteAsync runs the provided action asynchronously.
func (f *AtomstateFutureImpl[T]) WhenCompleteAsync(action func(T, error)) {
	f.WhenComplete(action)
}

// Handle returns a new AtomstateFuture that handles both the result and error.
func (f *AtomstateFutureImpl[T]) Handle(biFunction func(T, error) T) *AtomstateFutureImpl[T] {
	result := NewAtomstateFuture[T]()
	go func() {
		value, err := f.completableFuture.Get()
		result.Complete(biFunction(value, err))
	}()
	return result
}

// HandleAsync returns a new AtomstateFuture that handles both the result and error asynchronously.
func (f *AtomstateFutureImpl[T]) HandleAsync(biFunction func(T, error) T) *AtomstateFutureImpl[T] {
	return f.Handle(biFunction)
}

// Exceptionally returns a new AtomstateFuture that handles errors using the provided function.
func (f *AtomstateFutureImpl[T]) Exceptionally(function func(error) T) *AtomstateFutureImpl[T] {
	return f.Handle(func(value T, err error) T {
		if err != nil {
			return function(err)
		}
		return value
	})
}

// AllOf returns a new AtomstateFuture that is completed when all of the given futures are completed.
func (f *AtomstateFutureImpl[T]) AllOf(futures ...*AtomstateFutureImpl[interface{}]) *AtomstateFutureImpl[struct{}] {
	result := NewAtomstateFuture[struct{}]()
	var wg sync.WaitGroup
	var once sync.Once
	var firstError error

	// Increment the wait group counter for each future
	wg.Add(len(futures))

	for _, future := range futures {
		go func(f *AtomstateFutureImpl[interface{}]) {
			defer wg.Done()   // Decrement the wait group counter when done
			_, err := f.Get() // Get the result of the future
			if err != nil {
				once.Do(func() { firstError = err }) // Record the first error
			}
		}(future)
	}

	// Use a goroutine to handle the result after all futures are done
	go func() {
		wg.Wait() // Wait for all futures to complete
		if firstError != nil {
			result.CompleteExceptionally(firstError)
		} else {
			result.Complete(struct{}{})
		}
	}()

	return result
}

// AnyOf returns a new AtomstateFuture that is completed when any of the given futures are completed.
func (f *AtomstateFutureImpl[T]) AnyOf(futures ...*AtomstateFutureImpl[T]) *AtomstateFutureImpl[T] {
	result := NewAtomstateFuture[T]()
	var once sync.Once

	for _, future := range futures {
		go func(f *AtomstateFutureImpl[T]) {
			value, err := f.Get()
			once.Do(func() {
				if err != nil {
					result.CompleteExceptionally(err)
				} else {
					result.Complete(value)
				}
			})
		}(future)
	}

	return result
}

// Cancel cancels the future if possible.
func (f *AtomstateFutureImpl[T]) Cancel() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.cancelled {
		return false
	}
	f.cancelled = true
	return true
}

// IsCancelled checks if the future is cancelled.
func (f *AtomstateFutureImpl[T]) IsCancelled() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.cancelled
}

// IsDone checks if the future is completed.
func (f *AtomstateFutureImpl[T]) IsDone() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	select {
	case <-f.completableFuture.done:
		return true
	default:
		return false
	}
}
