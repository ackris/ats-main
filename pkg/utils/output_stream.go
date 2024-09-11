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
	"bytes"
)

// ByteBufferOutputStream is a dynamically expanding buffer that implements an output stream.
// It allows writing data into an internal buffer and expands the buffer as needed.
//
// Usage example:
//
//	package main
//
//	import (
//	    "fmt"
//	    "utils"
//	)
//
//	func main() {
//	    // Create a new ByteBufferOutputStream with an initial capacity of 10 bytes.
//	    stream := utils.NewByteBufferOutputStream(10)
//
//	    // Write a byte slice to the stream.
//	    data := []byte("Hello, World!")
//	    _, err := stream.Write(data)
//	    if err != nil {
//	        fmt.Println("Error writing data:", err)
//	        return
//	    }
//
//	    // Retrieve and print the current buffer content.
//	    fmt.Println("Buffer content:", string(stream.Buffer().Bytes()))
//
//	    // Change the position in the buffer.
//	    stream.SetPosition(15)
//	    fmt.Println("Buffer length after setting position:", stream.Len())
//	}
type ByteBufferOutputStream struct {
	buffer          *bytes.Buffer
	initialCapacity int
}

// NewByteBufferOutputStream creates a new ByteBufferOutputStream with the specified initial capacity.
//
// Parameters:
//
//	initialCapacity int: The initial capacity of the buffer.
//
// Returns:
//
//	*ByteBufferOutputStream: A new instance of ByteBufferOutputStream with the specified initial capacity.
//
// Example:
//
//	stream := NewByteBufferOutputStream(1024)
func NewByteBufferOutputStream(initialCapacity int) *ByteBufferOutputStream {
	return &ByteBufferOutputStream{
		buffer:          bytes.NewBuffer(make([]byte, 0, initialCapacity)),
		initialCapacity: initialCapacity,
	}
}

// Write writes the byte slice to the buffer, expanding the buffer if necessary.
//
// Parameters:
//
//	p []byte: The byte slice to write to the buffer.
//
// Returns:
//
//	int: The number of bytes written to the buffer.
//	error: An error if the write operation fails.
//
// Example:
//
//	data := []byte("Hello")
//	_, err := stream.Write(data)
//	if err != nil {
//	    fmt.Println("Error writing data:", err)
//	}
func (b *ByteBufferOutputStream) Write(p []byte) (n int, err error) {
	b.ensureCapacity(len(p))
	return b.buffer.Write(p)
}

// WriteByte writes a single byte to the buffer, expanding the buffer if necessary.
//
// Parameters:
//
//	p byte: The byte to write to the buffer.
//
// Returns:
//
//	error: An error if the write operation fails.
//
// Example:
//
//	err := stream.WriteByte('A')
//	if err != nil {
//	    fmt.Println("Error writing byte:", err)
//	}
func (b *ByteBufferOutputStream) WriteByte(p byte) error {
	b.ensureCapacity(1)
	return b.buffer.WriteByte(p)
}

// Buffer returns the underlying buffer, which contains the written data.
//
// Returns:
//
//	*bytes.Buffer: The underlying buffer.
//
// Example:
//
//	buf := stream.Buffer()
//	fmt.Println("Buffer content:", string(buf.Bytes()))
func (b *ByteBufferOutputStream) Buffer() *bytes.Buffer {
	return b.buffer
}

// Len returns the number of bytes written to the buffer.
//
// Returns:
//
//	int: The length of the buffer.
//
// Example:
//
//	length := stream.Len()
//	fmt.Println("Buffer length:", length)
func (b *ByteBufferOutputStream) Len() int {
	return b.buffer.Len()
}

// Position returns the current position in the buffer.
// In this implementation, it is effectively the length of the buffer.
//
// Returns:
//
//	int: The current position in the buffer.
//
// Example:
//
//	pos := stream.Position()
//	fmt.Println("Buffer position:", pos)
func (b *ByteBufferOutputStream) Position() int {
	return b.buffer.Len() // The position is effectively the length of the buffer in this context.
}

// Remaining returns the remaining capacity in the buffer.
// It calculates the difference between the buffer's capacity and length.
//
// Returns:
//
//	int: The remaining capacity in the buffer.
//
// Example:
//
//	remaining := stream.Remaining()
//	fmt.Println("Remaining capacity:", remaining)
func (b *ByteBufferOutputStream) Remaining() int {
	return b.buffer.Cap() - b.buffer.Len()
}

// Limit returns the total capacity of the buffer.
//
// Returns:
//
//	int: The total capacity of the buffer.
//
// Example:
//
//	limit := stream.Limit()
//	fmt.Println("Buffer limit:", limit)
func (b *ByteBufferOutputStream) Limit() int {
	return b.buffer.Cap()
}

// SetPosition sets the position to a specific value and expands the buffer if
// necessary to accommodate the new position.
//
// Parameters:
//
//	position int: The position to set in the buffer.
//
// Example:
//
//	stream.SetPosition(50)
//	fmt.Println("Buffer length after setting position:", stream.Len())
func (b *ByteBufferOutputStream) SetPosition(position int) {
	if position > b.buffer.Len() {
		b.ensureCapacity(position - b.buffer.Len())
	}
	// If position is beyond the current length, adjust the buffer length to match the position
	if position > b.buffer.Len() {
		b.buffer.Write(make([]byte, position-b.buffer.Len()))
	}
}

// InitialCapacity returns the initial capacity of the buffer that was set when the
// ByteBufferOutputStream was created.
//
// Returns:
//
//	int: The initial capacity of the buffer.
//
// Example:
//
//	initialCap := stream.InitialCapacity()
//	fmt.Println("Initial capacity:", initialCap)
func (b *ByteBufferOutputStream) InitialCapacity() int {
	return b.initialCapacity
}

// ensureCapacity ensures that the buffer has enough capacity to write the specified number of bytes.
// It expands the buffer if necessary based on a growth factor.
//
// Parameters:
//
//	required int: The number of bytes required.
//
// Example:
//
//	stream.ensureCapacity(100)
//	fmt.Println("Buffer capacity ensured.")
func (b *ByteBufferOutputStream) ensureCapacity(required int) {
	if b.buffer.Len()+required > b.buffer.Cap() {
		// Calculate the new capacity with a growth factor
		newCapacity := int(float64(b.buffer.Cap()) * 1.1)
		// Ensure that the new capacity is large enough to accommodate the required space
		if b.buffer.Len()+required > newCapacity {
			newCapacity = b.buffer.Len() + required
		}
		// Create a new buffer with the calculated capacity and copy the old buffer's contents
		newBuffer := bytes.NewBuffer(make([]byte, 0, newCapacity))
		newBuffer.Write(b.buffer.Bytes())
		b.buffer = newBuffer
	}
}
