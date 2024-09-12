package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCloseableIterator tests the CloseableIterator functionality.
func TestCloseableIterator(t *testing.T) {
	data := []int{1, 2, 3}
	inner := NewSimpleIterator(data)
	closerCalled := false

	// Create a CloseableIterator
	ci := NewCloseableIterator(inner, func() error {
		closerCalled = true
		return nil
	})

	// Test HasNext and Next
	assert.True(t, ci.HasNext())
	val, err := ci.Next()
	require.NoError(t, err)
	assert.Equal(t, 1, val)

	val, err = ci.Next()
	require.NoError(t, err)
	assert.Equal(t, 2, val)

	val, err = ci.Next()
	require.NoError(t, err)
	assert.Equal(t, 3, val)

	// Ensure no more elements
	assert.False(t, ci.HasNext())
	_, err = ci.Next()
	assert.Error(t, err)

	// Test Close
	err = ci.Close()
	assert.NoError(t, err)
	assert.True(t, closerCalled)
}

// TestFlattenedIterator tests the FlattenedIterator functionality.
func TestFlattenedIterator(t *testing.T) {
	outerData := []int{1, 2}
	innerData := map[int][]int{
		1: {10, 11},
		2: {20, 21},
	}

	outer := NewSimpleIterator(outerData)
	innerIteratorFunc := func(o int) Iterator[int] {
		return NewSimpleIterator(innerData[o])
	}

	fi := NewFlattenedIterator(outer, innerIteratorFunc)

	// Test HasNext and Next
	assert.True(t, fi.HasNext())
	val, err := fi.Next()
	require.NoError(t, err)
	assert.Equal(t, 10, val)

	val, err = fi.Next()
	require.NoError(t, err)
	assert.Equal(t, 11, val)

	assert.True(t, fi.HasNext())
	val, err = fi.Next()
	require.NoError(t, err)
	assert.Equal(t, 20, val)

	val, err = fi.Next()
	require.NoError(t, err)
	assert.Equal(t, 21, val)

	// Ensure no more elements
	assert.False(t, fi.HasNext())
	_, err = fi.Next()
	assert.Error(t, err)

	// Test Close
	err = fi.Close()
	assert.NoError(t, err)
}

// TestFlattenedIteratorWithEmptyInner tests the case where the inner iterator is empty.
/* func TestFlattenedIteratorWithEmptyInner(t *testing.T) {
	outerData := []int{1, 2}
	innerData := map[int][]int{
		1: {},       // Outer element 1 has no inner elements
		2: {20, 21}, // Outer element 2 has inner elements 20 and 21
	}

	outer := NewSimpleIterator(outerData)
	innerIteratorFunc := func(o int) Iterator[int] {
		return NewSimpleIterator(innerData[o])
	}

	fi := NewFlattenedIterator(outer, innerIteratorFunc)

	// Test HasNext and Next
	assert.True(t, fi.HasNext()) // Should be true for outer element 2
	_, err := fi.Next()          // Should not return an error for outer element 1
	require.NoError(t, err)

	assert.True(t, fi.HasNext()) // Should now be true for inner elements of outer element 2
	val, err := fi.Next()
	require.NoError(t, err)
	assert.Equal(t, 20, val) // Expecting 20

	val, err = fi.Next()
	require.NoError(t, err)
	assert.Equal(t, 21, val) // Expecting 21

	// Ensure no more elements
	assert.False(t, fi.HasNext()) // Should be false now
	_, err = fi.Next()
	assert.Error(t, err) // Expecting an error for no more elements

	// Test Close
	err = fi.Close()
	assert.NoError(t, err)
} */

// TestFlattenedIteratorWithError tests the case where the inner iterator returns an error.
func TestFlattenedIteratorWithError(t *testing.T) {
	outerData := []int{1, 2}
	innerData := map[int][]int{
		1: {10, 11},
		2: {},
	}

	outer := NewSimpleIterator(outerData)
	innerIteratorFunc := func(o int) Iterator[int] {
		if o == 2 {
			return &SimpleIterator{data: nil, index: 0} // Simulating an empty iterator
		}
		return NewSimpleIterator(innerData[o])
	}

	fi := NewFlattenedIterator(outer, innerIteratorFunc)

	// Test HasNext and Next
	assert.True(t, fi.HasNext())
	val, err := fi.Next()
	require.NoError(t, err)
	assert.Equal(t, 10, val)

	val, err = fi.Next()
	require.NoError(t, err)
	assert.Equal(t, 11, val)

	// Ensure no more elements
	assert.False(t, fi.HasNext())
	_, err = fi.Next()
	assert.Error(t, err)

	// Test Close
	err = fi.Close()
	assert.NoError(t, err)
}
