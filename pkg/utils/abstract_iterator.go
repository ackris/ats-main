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

// AbstractIterator provides a way to iterate over a sequence of integers using a custom function.
type AbstractIterator struct {
	state        State
	next         int
	iteratorFunc func() (int, bool, error)
}

// State enum to manage the state of the iterator.
type State int

const (
	// READY indicates that the iterator is ready to return the next element.
	READY State = iota
	// NOT_READY indicates that the next element has not been computed yet.
	NOT_READY
	// DONE indicates that there are no more elements to iterate over.
	DONE
	// FAILED indicates that an error occurred while trying to compute the next element.
	FAILED
)

// NewAbstractIterator creates a new AbstractIterator with the provided function.
// The function should return the next integer, a boolean indicating if there are more elements,
// and an error if something goes wrong.
//
// Example usage:
//
//	iterFunc := func() (int, bool, error) {
//	    // Logic to return the next integer, for example from a slice.
//	}
//	iterator := NewAbstractIterator(iterFunc)
func NewAbstractIterator(f func() (int, bool, error)) *AbstractIterator {
	return &AbstractIterator{
		state:        NOT_READY,
		iteratorFunc: f,
	}
}

// HasNext checks if there are more elements to iterate over.
// It returns true if the iterator is ready to provide the next element,
// false otherwise (including when it has reached the end or encountered an error).
//
// Example usage:
//
//	if iterator.HasNext() {
//	    // Proceed to get the next element.
//	}
func (iter *AbstractIterator) HasNext() bool {
	switch iter.state {
	case FAILED:
		return false
	case DONE:
		return false
	case READY:
		return true
	case NOT_READY:
		return iter.maybeComputeNext()
	default:
		panic("invalid state")
	}
}

// Next retrieves the next element from the iterator.
// It returns the next integer and an error if there are no more elements.
//
// Example usage:
//
//	nextValue, err := iterator.Next()
//	if err != nil {
//	    // Handle the error (e.g., no more elements).
//	}
func (iter *AbstractIterator) Next() (int, error) {
	if !iter.HasNext() {
		return 0, errors.New("no more elements")
	}
	nextValue := iter.next
	iter.state = NOT_READY // Prepare to compute the next value
	return nextValue, nil
}

// Peek returns the next element without advancing the iterator.
// It allows you to see the next element without changing the iterator's state.
//
// Example usage:
//
//	nextValue, err := iterator.Peek()
//	if err != nil {
//	    // Handle the error (e.g., no more elements).
//	}
func (iter *AbstractIterator) Peek() (int, error) {
	if !iter.HasNext() {
		return 0, errors.New("no more elements")
	}
	return iter.next, nil
}

// maybeComputeNext tries to compute the next element.
// It updates the iterator's state based on the result of the iterator function.
// Returns true if a next element was successfully computed, false otherwise.
func (iter *AbstractIterator) maybeComputeNext() bool {
	if iter.state == DONE {
		return false
	}
	iter.state = FAILED
	var ok bool
	var err error
	iter.next, ok, err = iter.iteratorFunc()
	if err != nil {
		iter.state = FAILED
		return false
	}
	if !ok {
		iter.state = DONE
		return false
	}
	iter.state = READY
	return true
}
