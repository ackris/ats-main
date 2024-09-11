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

// TestNewByteBufferOutputStream verifies the creation of a ByteBufferOutputStream with initial capacity.
func TestNewByteBufferOutputStream(t *testing.T) {
	initialCapacity := 1024
	stream := NewByteBufferOutputStream(initialCapacity)
	if stream.InitialCapacity() != initialCapacity {
		t.Errorf("Expected initial capacity %d, got %d", initialCapacity, stream.InitialCapacity())
	}
	if cap(stream.Buffer().Bytes()) != initialCapacity {
		t.Errorf("Expected buffer capacity %d, got %d", initialCapacity, cap(stream.Buffer().Bytes()))
	}
}

// TestWrite verifies writing to the buffer and buffer expansion.
func TestWrite(t *testing.T) {
	stream := NewByteBufferOutputStream(10)
	data := []byte("Hello, World!")

	n, err := stream.Write(data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d bytes", len(data), n)
	}
	if stream.Len() != len(data) {
		t.Errorf("Expected buffer length %d, got %d", len(data), stream.Len())
	}
	if stream.Remaining() < 0 {
		t.Errorf("Expected remaining capacity to be non-negative, got %d", stream.Remaining())
	}
}

// TestWriteByte verifies writing a single byte to the buffer.
func TestWriteByte(t *testing.T) {
	stream := NewByteBufferOutputStream(10)
	err := stream.WriteByte('A')
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if stream.Len() != 1 {
		t.Errorf("Expected buffer length 1, got %d", stream.Len())
	}
	if stream.Buffer().Bytes()[0] != 'A' {
		t.Errorf("Expected buffer to start with byte 'A', got %c", stream.Buffer().Bytes()[0])
	}
}

// TestSetPosition verifies that the position can be set and the buffer is correctly expanded if necessary.
func TestSetPosition(t *testing.T) {
	stream := NewByteBufferOutputStream(10)
	data := []byte("Hello")
	_, err := stream.Write(data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	stream.SetPosition(20) // Move position beyond current length
	if stream.Len() != 20 {
		t.Errorf("Expected buffer length to be 20 after setting position, got %d", stream.Len())
	}
}

// TestBufferExpansion verifies the behavior of buffer expansion.
func TestBufferExpansion(t *testing.T) {
	initialCapacity := 10
	stream := NewByteBufferOutputStream(initialCapacity)
	data := make([]byte, initialCapacity*2) // Data larger than initial capacity

	_, err := stream.Write(data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if stream.Len() != len(data) {
		t.Errorf("Expected buffer length %d, got %d", len(data), stream.Len())
	}
	if stream.Remaining() < 0 {
		t.Errorf("Expected remaining capacity to be non-negative, got %d", stream.Remaining())
	}
}

// TestPosition verifies the correct position tracking in the buffer.
func TestPosition(t *testing.T) {
	stream := NewByteBufferOutputStream(20)
	data := []byte("Data")
	_, err := stream.Write(data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if pos := stream.Position(); pos != len(data) {
		t.Errorf("Expected position %d, got %d", len(data), pos)
	}
}

// TestInitialCapacity verifies that the initial capacity is correctly stored and returned.
func TestInitialCapacity(t *testing.T) {
	initialCapacity := 512
	stream := NewByteBufferOutputStream(initialCapacity)
	if initialCapacityReturned := stream.InitialCapacity(); initialCapacityReturned != initialCapacity {
		t.Errorf("Expected initial capacity %d, got %d", initialCapacity, initialCapacityReturned)
	}
}
