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

package common

import (
	"errors"
	"sync"
	"time"
)

// AtomstateFuture represents a future that can be completed with a value or an error.
type AtomstateFuture[T any] struct {
	result    T
	err       error
	done      chan struct{}
	cancelled bool
	mu        sync.Mutex
	once      sync.Once
}

// NewAtomstateFuture creates a new AtomstateFuture.
func NewAtomstateFuture[T any]() *AtomstateFuture[T] {
	return &AtomstateFuture[T]{done: make(chan struct{})}
}

// CompletedFuture returns a AtomstateFuture that is already completed with the given value.
func CompletedFuture[T any](value T) *AtomstateFuture[T] {
	future := NewAtomstateFuture[T]()
	future.Complete(value)
	return future
}

// Complete sets the result of the future.
func (af *AtomstateFuture[T]) Complete(value T) {
	af.once.Do(func() {
		af.mu.Lock()
		defer af.mu.Unlock()
		if af.err == nil && !af.cancelled {
			af.result = value
			close(af.done)
		}
	})
}

// CompleteExceptionally sets the error of the future.
func (af *AtomstateFuture[T]) CompleteExceptionally(err error) {
	af.once.Do(func() {
		af.mu.Lock()
		defer af.mu.Unlock()
		if af.err == nil && !af.cancelled {
			af.err = err
			close(af.done)
		}
	})
}

// IsDone checks if the future is completed.
func (af *AtomstateFuture[T]) IsDone() bool {
	select {
	case <-af.done:
		return true
	default:
		return false
	}
}

// IsCancelled checks if the future was cancelled.
func (af *AtomstateFuture[T]) IsCancelled() bool {
	af.mu.Lock()
	defer af.mu.Unlock()
	return af.cancelled
}

// GetNow returns the result if completed, else returns the given valueIfAbsent.
func (af *AtomstateFuture[T]) GetNow(valueIfAbsent T) (T, error) {
	af.mu.Lock()
	defer af.mu.Unlock()
	if af.err != nil {
		return valueIfAbsent, af.err
	}
	return af.result, nil
}

// Get waits for the future to complete and returns the result or an error.
func (af *AtomstateFuture[T]) Get() (T, error) {
	<-af.done
	af.mu.Lock()
	defer af.mu.Unlock()
	return af.result, af.err
}

// GetWithTimeout waits for the future to complete for the specified duration.
func (af *AtomstateFuture[T]) GetWithTimeout(timeout time.Duration) (T, error) {
	select {
	case <-time.After(timeout):
		af.mu.Lock()
		defer af.mu.Unlock()
		return af.result, errors.New("timeout waiting for future")
	case <-af.done:
		return af.result, af.err
	}
}

// AllOf returns a new AtomstateFuture that is completed when all the given futures have completed.
func AllOf[T any](futures ...*AtomstateFuture[T]) *AtomstateFuture[struct{}] {
	result := NewAtomstateFuture[struct{}]()
	var wg sync.WaitGroup
	wg.Add(len(futures))

	for _, future := range futures {
		go func(f *AtomstateFuture[T]) {
			defer wg.Done()
			_, err := f.Get()
			if err != nil {
				result.CompleteExceptionally(err)
			}
		}(future)
	}

	go func() {
		wg.Wait()
		if result.err == nil {
			result.Complete(struct{}{})
		}
	}()

	return result
}

// ThenApply executes the provided function when the future completes.
func (af *AtomstateFuture[T]) ThenApply(fn func(T) interface{}) *AtomstateFuture[interface{}] {
	result := NewAtomstateFuture[interface{}]()
	go func() {
		value, err := af.Get()
		if err != nil {
			result.CompleteExceptionally(err)
			return
		}
		result.Complete(fn(value))
	}()
	return result
}

// WhenComplete executes the provided action when the future completes.
func (af *AtomstateFuture[T]) WhenComplete(action func(T, error)) *AtomstateFuture[T] {
	go func() {
		value, err := af.Get()
		action(value, err)
	}()
	return af
}

// Cancel cancels the future if not already completed.
func (af *AtomstateFuture[T]) Cancel(mayInterruptIfRunning bool) bool {
	af.once.Do(func() {
		af.mu.Lock()
		defer af.mu.Unlock()
		if af.cancelled || af.IsDone() {
			return
		}
		af.cancelled = true
		close(af.done)
	})
	return true
}

// BaseFunction is a function type that takes an input of type A and returns a value of type B.
type BaseFunction[A any, B any] func(A) B

// BiConsumer is a function type that takes two inputs of types A and B.
type BiConsumer[A any, B any] func(A, B)
