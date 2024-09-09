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
// It allows you to handle asynchronous operations and their results or errors.
type AtomstateFuture[T any] interface {
	// Complete completes the future with the provided value.
	Complete(value T) error

	// CompleteExceptionally completes the future with an error.
	CompleteExceptionally(err error) error

	// Get retrieves the value of the future, blocking until it is completed.
	Get() (T, error)

	// GetNow retrieves the value of the future if it's completed, otherwise returns the default value.
	GetNow(defaultValue T) T

	// Accept runs the provided action with the value of this future.
	Accept(action func(T))

	// ToCompletionStage returns the AtomstateFutureImplementation itself.
	ToCompletionStage() *AtomstateFutureImplementation[T]

	// ToString returns a string representation of the future's state.
	ToString() string

	// ThenApply returns a new AtomstateFuture that applies the provided function to the result of this future.
	ThenApply(function func(T) T) *AtomstateFutureImplementation[T]

	// ThenApplyAsync returns a new AtomstateFuture that applies the provided function asynchronously.
	ThenApplyAsync(function func(T) T) *AtomstateFutureImplementation[T]

	// WhenComplete runs the provided action when this future completes.
	WhenComplete(action func(T, error))

	// WhenCompleteAsync runs the provided action asynchronously.
	WhenCompleteAsync(action func(T, error))

	// Handle returns a new AtomstateFuture that handles both the result and error.
	Handle(biFunction func(T, error) T) *AtomstateFutureImplementation[T]

	// HandleAsync returns a new AtomstateFuture that handles both the result and error asynchronously.
	HandleAsync(biFunction func(T, error) T) *AtomstateFutureImplementation[T]

	// Exceptionally returns a new AtomstateFuture that handles errors using the provided function.
	Exceptionally(function func(error) T) *AtomstateFutureImplementation[T]

	// AllOf returns a new AtomstateFuture that is completed when all of the given futures are completed.
	AllOf(futures ...*AtomstateFutureImplementation[interface{}]) *AtomstateFutureImplementation[struct{}]

	// AnyOf returns a new AtomstateFuture that is completed when any of the given futures are completed.
	AnyOf(futures ...*AtomstateFutureImplementation[T]) *AtomstateFutureImplementation[T]

	// Cancel cancels the future if possible.
	Cancel() bool

	// IsCancelled checks if the future is cancelled.
	IsCancelled() bool

	// IsDone checks if the future is completed.
	IsDone() bool
}

// AtomstateFutureImplementation implements AtomstateFuture interface.
// It provides methods to manage and retrieve results from an asynchronous operation.
type AtomstateFutureImplementation[T any] struct {
	completableFuture *AtomstateCompletableFuture[T]
	cancelled         bool
	mu                sync.RWMutex
}

// NewAtomstateFuture creates a new AtomstateFutureImplementation instance.
//
// Example:
//
// future := NewAtomstateFuture[int]()
func NewAtomstateFuture[T any]() *AtomstateFutureImplementation[T] {
	return &AtomstateFutureImplementation[T]{
		completableFuture: NewAtomstateCompletableFuture[T](),
	}
}

// Complete completes the future with the provided value.
// Returns an error if the future was already cancelled.
// Example usage:
// err := future.Complete(42)
//
//	if err != nil {
//	    fmt.Println("Error completing future:", err)
//	}
func (f *AtomstateFutureImplementation[T]) Complete(value T) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.cancelled {
		return errors.New("future is cancelled")
	}
	return f.completableFuture.Complete(value)
}

// CompleteExceptionally completes the future with an error.
// Returns an error if the future was already cancelled.
// Example usage:
// err := future.CompleteExceptionally(errors.New("something went wrong"))
//
//	if err != nil {
//	    fmt.Println("Error completing future exceptionally:", err)
//	}
func (f *AtomstateFutureImplementation[T]) CompleteExceptionally(err error) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.cancelled {
		return errors.New("future is cancelled")
	}
	return f.completableFuture.CompleteExceptionally(err)
}

// Get retrieves the value of the future, blocking until it is completed.
// Returns the value and any error that occurred during completion.
// Example usage:
// value, err := future.Get()
//
//	if err != nil {
//	    fmt.Println("Error getting future value:", err)
//	} else {
//
//	    fmt.Println("Future value:", value)
//	}
func (f *AtomstateFutureImplementation[T]) Get() (T, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.cancelled {
		return *new(T), errors.New("future is cancelled")
	}
	return f.completableFuture.Get()
}

// GetNow retrieves the value of the future if it's completed,
// otherwise returns the default value provided.
//
// Example usage:
//
// value := future.GetNow(0)
//
// fmt.Println("Future value or default:", value)
func (f *AtomstateFutureImplementation[T]) GetNow(defaultValue T) T {
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

// Accept runs the provided action with the value of this future once it is completed.
// The action is executed in a new goroutine.
// Example usage:
//
//	future.Accept(func(value int) {
//	    fmt.Println("Accepted value:", value)
//	})
func (f *AtomstateFutureImplementation[T]) Accept(action func(T)) {
	go func() {
		value, err := f.completableFuture.Get()
		if err == nil {
			action(value)
		}
	}()
}

// ToCompletionStage returns the AtomstateFutureImplementation itself.
// This allows method chaining with the future instance.
// Example usage:
// stage := future.ToCompletionStage()
// fmt.Println(stage)
func (f *AtomstateFutureImplementation[T]) ToCompletionStage() *AtomstateFutureImplementation[T] {
	return f
}

// ToString returns a string representation of the future's state.
// Example usage:
//
// fmt.Println(future.ToString())
//
// Outputs "Incomplete", "Completed[42]", or "Cancelled"
func (f *AtomstateFutureImplementation[T]) ToString() string {
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
// The function is applied in a new goroutine.
// Example usage:
//
//	newFuture := future.ThenApply(func(v int) int {
//	    return v * 2
//	})
//
// value, err := newFuture.Get()
// fmt.Println("Transformed Value:", value)
func (f *AtomstateFutureImplementation[T]) ThenApply(function func(T) T) *AtomstateFutureImplementation[T] {
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
// The function is applied in a new goroutine.
// Example usage:
//
//	newFuture := future.ThenApplyAsync(func(v int) int {
//	    return v + 10
//	})
//
// value, err := newFuture.Get()
// fmt.Println("Asynchronously Transformed Value:", value)
func (f *AtomstateFutureImplementation[T]) ThenApplyAsync(function func(T) T) *AtomstateFutureImplementation[T] {
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
// The action is executed in a new goroutine and receives the result and error.
// Example usage:
//
//	future.WhenComplete(func(value int, err error) {
//	    if err != nil {
//	        fmt.Println("Completed with error:", err)
//	    } else {
//	        fmt.Println("Completed with value:", value)
//	    }
//	})
func (f *AtomstateFutureImplementation[T]) WhenComplete(action func(T, error)) {
	go func() {
		value, err := f.completableFuture.Get()
		action(value, err)
	}()
}

// WhenCompleteAsync runs the provided action asynchronously.
// It’s similar to WhenComplete but explicitly uses a new goroutine.
// Example usage:
//
//	future.WhenCompleteAsync(func(value int, err error) {
//	    if err != nil {
//	        fmt.Println("Asynchronously completed with error:", err)
//	    } else {
//	        fmt.Println("Asynchronously completed with value:", value)
//	    }
//	})
func (f *AtomstateFutureImplementation[T]) WhenCompleteAsync(action func(T, error)) {
	f.WhenComplete(action)
}

// Handle returns a new AtomstateFuture that handles both the result and error.
// The provided function is applied to the result and error and is completed with the result.
// Example usage:
//
//	handledFuture := future.Handle(func(value int, err error) int {
//	    if err != nil {
//	        return 0
//	    }
//	    return value * 2
//	})
//
// value, err := handledFuture.Get()
// fmt.Println("Handled Value:", value)
func (f *AtomstateFutureImplementation[T]) Handle(biFunction func(T, error) T) *AtomstateFutureImplementation[T] {
	result := NewAtomstateFuture[T]()
	go func() {
		value, err := f.completableFuture.Get()
		result.Complete(biFunction(value, err))
	}()
	return result
}

// HandleAsync returns a new AtomstateFuture that handles both the result and error asynchronously.
// It’s similar to Handle but explicitly uses a new goroutine.
// Example usage:
//
//	handledFuture := future.HandleAsync(func(value int, err error) int {
//	    if err != nil {
//	        return -1
//	    }
//	    return value + 1
//	})
//
// value, err := handledFuture.Get()
// fmt.Println("Asynchronously Handled Value:", value)
func (f *AtomstateFutureImplementation[T]) HandleAsync(biFunction func(T, error) T) *AtomstateFutureImplementation[T] {
	return f.Handle(biFunction)
}

// Exceptionally returns a new AtomstateFuture that handles errors using the provided function.
// The function is applied only if an error occurs.
// Example usage:
//
//	handledFuture := future.Exceptionally(func(err error) int {
//	    return -1
//	})
//
// value, err := handledFuture.Get()
// fmt.Println("Exceptionally Handled Value:", value)
func (f *AtomstateFutureImplementation[T]) Exceptionally(function func(error) T) *AtomstateFutureImplementation[T] {
	return f.Handle(func(value T, err error) T {
		if err != nil {
			return function(err)
		}
		return value
	})
}

// AllOf returns a new AtomstateFuture that is completed when all of the given futures are completed.
// It completes with an error if any of the futures complete exceptionally.
// Example usage:
// future1 := NewAtomstateFuture[int]()
// future2 := NewAtomstateFuture[int]()
// allOfFuture := future.AllOf(future1, future2)
// _, err := allOfFuture.Get()
// fmt.Println("All futures completed:", err == nil)
func (f *AtomstateFutureImplementation[T]) AllOf(futures ...*AtomstateFutureImplementation[interface{}]) *AtomstateFutureImplementation[struct{}] {
	result := NewAtomstateFuture[struct{}]()
	var wg sync.WaitGroup
	var once sync.Once
	var firstError error

	// Increment the wait group counter for each future
	wg.Add(len(futures))

	for _, future := range futures {
		go func(f *AtomstateFutureImplementation[interface{}]) {
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
// It completes with the value of the first future that completes successfully.
// Example usage:
// future1 := NewAtomstateFuture[int]()
// future2 := NewAtomstateFuture[int]()
// anyOfFuture := future.AnyOf(future1, future2)
// value, err := anyOfFuture.Get()
// fmt.Println("Any future completed with value:", value, "Error:", err)
func (f *AtomstateFutureImplementation[T]) AnyOf(futures ...*AtomstateFutureImplementation[T]) *AtomstateFutureImplementation[T] {
	result := NewAtomstateFuture[T]()
	var once sync.Once

	for _, future := range futures {
		go func(f *AtomstateFutureImplementation[T]) {
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
// Returns true if the future was successfully cancelled, false otherwise.
// Example usage:
// cancelled := future.Cancel()
// fmt.Println("Future cancelled:", cancelled)
func (f *AtomstateFutureImplementation[T]) Cancel() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.cancelled {
		return false
	}
	f.cancelled = true
	return true
}

// IsCancelled checks if the future is cancelled.
// Example usage:
//
//	if future.IsCancelled() {
//	    fmt.Println("Future is cancelled")
//	} else {
//
//	    fmt.Println("Future is not cancelled")
//	}
func (f *AtomstateFutureImplementation[T]) IsCancelled() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.cancelled
}

// IsDone checks if the future is completed.
// Example usage:
//
//	if future.IsDone() {
//	    fmt.Println("Future is completed")
//	} else {
//
//	    fmt.Println("Future is not completed")
//	}
func (f *AtomstateFutureImplementation[T]) IsDone() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	select {
	case <-f.completableFuture.done:
		return true
	default:
		return false
	}
}
