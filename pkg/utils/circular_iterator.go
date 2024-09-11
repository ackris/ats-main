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
	"errors"
)

// CircularIterator is an iterator that cycles through a slice indefinitely.
// It is useful for scenarios where you need to repeatedly iterate over a collection.
type CircularIterator[T any] struct {
	collection []T // The underlying collection to iterate over
	index      int // The current index in the collection
}

// NewCircularIterator creates a new CircularIterator for the provided slice.
// It returns an error if the slice is empty or nil.
//
// Example usage:
//
//	iterator, err := NewCircularIterator([]int{1, 2, 3})
//	if err != nil {
//	    // handle error
//	}
//	for i := 0; i < 10; i++ {
//	    fmt.Println(iterator.Next()) // Will print numbers 1, 2, 3 repeatedly
//	}
func NewCircularIterator[T any](collection []T) (*CircularIterator[T], error) {
	if collection == nil {
		return nil, errors.New("collection cannot be nil")
	}
	if len(collection) == 0 {
		return nil, errors.New("CircularIterator can only be used on non-empty slices")
	}
	return &CircularIterator[T]{collection: collection}, nil
}

// HasNext always returns true since the iteration cycles indefinitely.
// This method does not need to be called to use the iterator, but it can be used
// to check if the iterator is still valid (it always is).
//
// Example usage:
//
//	if iterator.HasNext() {
//	    fmt.Println(iterator.Next())
//	}
func (ci *CircularIterator[T]) HasNext() bool {
	return true
}

// Next returns the next element in the collection, cycling back to the start if necessary.
// It is assumed that the caller checks HasNext() before calling Next().
//
// Example usage:
//
//	for i := 0; i < 10; i++ {
//	    fmt.Println(iterator.Next()) // Will print numbers 1, 2, 3 repeatedly
//	}
func (ci *CircularIterator[T]) Next() T {
	value := ci.collection[ci.index]
	ci.index = (ci.index + 1) % len(ci.collection) // Cycle index
	return value
}

// Peek returns the next element without advancing the iterator.
// It is assumed that the caller checks HasNext() before calling Peek().
//
// Example usage:
//
//	fmt.Println(iterator.Peek()) // Prints the next value without advancing
func (ci *CircularIterator[T]) Peek() T {
	return ci.collection[ci.index]
}
