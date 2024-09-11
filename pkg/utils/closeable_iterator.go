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

// Iterator defines the basic iteration interface.
// It provides methods to check if there are more elements to iterate over
// and to retrieve the next element.
type Iterator[T any] interface {
	// HasNext returns true if there are more elements to iterate over.
	HasNext() bool

	// Next retrieves the next element in the iteration.
	// It returns an error if there are no more elements or if an error occurs.
	Next() (T, error)
}

// Closeable defines the interface for closing or releasing resources.
// It ensures that any resources held by the iterator can be properly released.
type Closeable interface {
	// Close releases any resources held by the Closeable instance.
	// It returns an error if the closing operation fails.
	Close() error
}

// CloseableIterator combines the Iterator and Closeable interfaces.
// It provides an iterator that can also be closed after use.
type CloseableIterator[T any] interface {
	Iterator[T]
	Closeable
}

// Concrete implementation of CloseableIterator.
type closeableIterator[T any] struct {
	inner  Iterator[T]  // The underlying iterator
	closer func() error // Function to close resources
}

// NewCloseableIterator creates a new CloseableIterator instance.
// It takes an inner Iterator and a closer function as parameters.
// The closer function is called when the Close method is invoked.
func NewCloseableIterator[T any](inner Iterator[T], closer func() error) CloseableIterator[T] {
	return &closeableIterator[T]{inner: inner, closer: closer}
}

// Wrap takes an existing Iterator and wraps it into a CloseableIterator.
// This function is a convenience method to create a CloseableIterator
// from an existing Iterator.
func Wrap[T any](inner Iterator[T], closer func() error) CloseableIterator[T] {
	return NewCloseableIterator(inner, closer)
}

// HasNext checks if there are more elements to iterate over.
// It delegates the call to the inner iterator.
func (ci *closeableIterator[T]) HasNext() bool {
	return ci.inner.HasNext()
}

// Next retrieves the next element in the iteration.
// It delegates the call to the inner iterator.
func (ci *closeableIterator[T]) Next() (T, error) {
	return ci.inner.Next()
}

// Close releases any resources held by the closeable iterator.
// It calls the closer function and returns any error encountered.
func (ci *closeableIterator[T]) Close() error {
	return ci.closer()
}
