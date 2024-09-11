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
	"testing"
)

func TestNewCircularIterator(t *testing.T) {
	t.Run("Valid Collection", func(t *testing.T) {
		iterator, err := NewCircularIterator([]int{1, 2, 3})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if iterator == nil {
			t.Fatal("Expected iterator to be non-nil")
		}
	})

	t.Run("Empty Collection", func(t *testing.T) {
		_, err := NewCircularIterator([]int{})
		if err == nil {
			t.Fatal("Expected error for empty collection, got none")
		}
	})

	t.Run("Nil Collection", func(t *testing.T) {
		_, err := NewCircularIterator([]*int(nil)) // Specify a type for nil
		if err == nil {
			t.Fatal("Expected error for nil collection, got none")
		}
	})
}

func TestCIHasNext(t *testing.T) {
	iterator, _ := NewCircularIterator([]int{1, 2, 3})

	for i := 0; i < 10; i++ {
		if !iterator.HasNext() {
			t.Fatal("Expected HasNext to return true")
		}
		iterator.Next() // Advance the iterator
	}
}

func TestCINext(t *testing.T) {
	iterator, _ := NewCircularIterator([]int{1, 2, 3})

	expected := []int{1, 2, 3, 1, 2, 3}
	for i := 0; i < len(expected); i++ {
		if got := iterator.Next(); got != expected[i] {
			t.Fatalf("Expected %d, got %d", expected[i], got)
		}
	}
}

func TestCIPeek(t *testing.T) {
	iterator, _ := NewCircularIterator([]int{1, 2, 3})

	for i := 0; i < 3; i++ {
		if got := iterator.Peek(); got != 1+i {
			t.Fatalf("Expected %d, got %d", 1+i, got)
		}
		iterator.Next() // Advance the iterator
	}
}

func TestCircularBehavior(t *testing.T) {
	iterator, _ := NewCircularIterator([]int{1, 2, 3})

	// Advance through the iterator multiple times
	for i := 0; i < 10; i++ {
		got := iterator.Next()
		expected := (i % 3) + 1 // Expected values are 1, 2, 3, 1, 2, 3, ...
		if got != expected {
			t.Fatalf("Expected %d, got %d at iteration %d", expected, got, i)
		}
	}
}

func TestPeekWithoutAdvance(t *testing.T) {
	collection := []int{1, 2, 3}
	iterator, _ := NewCircularIterator(collection)

	// Get the expected value from the iterator
	expected := iterator.Peek()

	// Peek multiple times without advancing
	for i := 0; i < len(collection); i++ {
		if got := iterator.Peek(); got != expected {
			t.Fatalf("Expected %d, got %d", expected, got)
		}
		iterator.Next()            // Advance the iterator after the peek
		expected = iterator.Peek() // Update expected to the next value
	}
}
