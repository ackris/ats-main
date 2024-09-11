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
	"testing"
)

// TestNewAbstractIterator tests the creation of a new AbstractIterator
func TestNewAbstractIterator(t *testing.T) {
	f := func() (int, bool, error) {
		return 0, false, nil
	}
	iter := NewAbstractIterator(f)

	if iter.state != NOT_READY {
		t.Errorf("Expected state to be NOT_READY, got %v", iter.state)
	}
}

// TestHasNext tests the HasNext method
func TestHasNext(t *testing.T) {
	count := 0
	f := func() (int, bool, error) {
		if count == 0 {
			count++
			return 1, true, nil
		}
		return 0, false, nil
	}
	iter := NewAbstractIterator(f)

	if !iter.HasNext() {
		t.Error("Expected HasNext to return true")
	}

	// Consume the next element
	iter.Next()

	if iter.HasNext() {
		t.Error("Expected HasNext to return false after consuming the only element")
	}
}

// TestNext tests the Next method
func TestNext(t *testing.T) {
	count := 0
	f := func() (int, bool, error) {
		if count < 3 {
			count++
			return count, true, nil
		}
		return 0, false, nil
	}
	iter := NewAbstractIterator(f)

	for i := 1; i <= 3; i++ {
		val, err := iter.Next()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}

	_, err := iter.Next()
	if err == nil {
		t.Error("Expected error when calling Next after all elements consumed")
	}
}

// TestPeek tests the Peek method
func TestPeek(t *testing.T) {
	count := 0
	f := func() (int, bool, error) {
		if count < 2 {
			count++
			return count, true, nil
		}
		return 0, false, nil
	}
	iter := NewAbstractIterator(f)

	val, err := iter.Peek()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}

	iter.Next() // Consume the first element

	val, err = iter.Peek()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}

	iter.Next() // Consume the second element

	_, err = iter.Peek()
	if err == nil {
		t.Error("Expected error when peeking after all elements consumed")
	}
}

// TestFailedIterator tests the iterator when an error occurs
func TestFailedIterator(t *testing.T) {
	f := func() (int, bool, error) {
		return 0, false, errors.New("some error")
	}
	iter := NewAbstractIterator(f)

	if iter.HasNext() {
		t.Error("Expected HasNext to return false due to failure")
	}

	_, err := iter.Next()
	if err == nil {
		t.Error("Expected error when calling Next on a failed iterator")
	}
}

// TestPeekOnFailedIterator tests Peek on a failed iterator
func TestPeekOnFailedIterator(t *testing.T) {
	f := func() (int, bool, error) {
		return 0, false, errors.New("some error")
	}
	iter := NewAbstractIterator(f)

	_, err := iter.Peek()
	if err == nil {
		t.Error("Expected error when peeking on a failed iterator")
	}
}

// TestMultipleCalls tests multiple calls to HasNext, Next, and Peek
func TestMultipleCalls(t *testing.T) {
	count := 0
	f := func() (int, bool, error) {
		if count < 3 {
			count++
			return count, true, nil
		}
		return 0, false, nil
	}
	iter := NewAbstractIterator(f)

	// First call to HasNext
	if !iter.HasNext() {
		t.Error("Expected HasNext to return true")
	}

	// First call to Next
	val, err := iter.Next()
	if err != nil || val != 1 {
		t.Errorf("Expected 1, got %d, error: %v", val, err)
	}

	// Check Peek after first Next
	val, err = iter.Peek()
	if err != nil || val != 2 {
		t.Errorf("Expected 2, got %d, error: %v", val, err)
	}

	// Second call to Next
	val, err = iter.Next()
	if err != nil || val != 2 {
		t.Errorf("Expected 2, got %d, error: %v", val, err)
	}

	// Check HasNext again
	if !iter.HasNext() {
		t.Error("Expected HasNext to return true")
	}

	// Third call to Next
	val, err = iter.Next()
	if err != nil || val != 3 {
		t.Errorf("Expected 3, got %d, error: %v", val, err)
	}

	// Check HasNext after consuming all elements
	if iter.HasNext() {
		t.Error("Expected HasNext to return false after consuming all elements")
	}
}
