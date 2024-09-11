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
	"io"
)

// ByteBufferInputStream represents a reader backed by a byte slice.
// It allows reading data from a buffer in a manner similar to an InputStream in Java.
//
// Example usage:
//
//	buffer := []byte("Hello, World!")
//	reader := NewByteBufferInputStream(buffer)
//
//	// Reading data
//	p := make([]byte, 5)
//	n, err := reader.Read(p)
//	if err != nil && err != io.EOF {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(p[:n])) // Output: Hello
//
//	// Checking available bytes
//	available := reader.Available()
//	fmt.Println(available) // Output: 8
//
//	// Reading remaining data
//	p = make([]byte, 8)
//	n, err = reader.Read(p)
//	if err != nil && err != io.EOF {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(p[:n])) // Output: , World!
//	fmt.Println(reader.Available()) // Output: 0
type ByteBufferInputStream struct {
	buffer []byte
	pos    int
}

// NewByteBufferInputStream creates a new ByteBufferInputStream with the given buffer.
//
// Parameters:
//   - buffer: A byte slice to back the input stream.
//
// Returns:
//   - A pointer to a new ByteBufferInputStream instance.
//
// Example:
//
//	buffer := []byte("Example")
//	reader := NewByteBufferInputStream(buffer)
func NewByteBufferInputStream(buffer []byte) *ByteBufferInputStream {
	return &ByteBufferInputStream{
		buffer: buffer,
		pos:    0,
	}
}

// Read reads up to len(p) bytes into p from the buffer. It returns the number of bytes read and an error, if any.
// The function will read as many bytes as possible but will not exceed the length of the provided slice.
//
// Parameters:
//   - p: A byte slice into which data will be read.
//
// Returns:
//   - The number of bytes read.
//   - An error, if any (e.g., io.EOF if the end of the buffer is reached).
//
// Example:
//
//	p := make([]byte, 10)
//	n, err := reader.Read(p)
//	if err != nil && err != io.EOF {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(p[:n])) // Output depends on buffer contents
func (b *ByteBufferInputStream) Read(p []byte) (int, error) {
	if b.pos >= len(b.buffer) {
		return 0, io.EOF
	}

	bytesToRead := len(p)
	remainingBytes := len(b.buffer) - b.pos
	if bytesToRead > remainingBytes {
		bytesToRead = remainingBytes
	}

	// Copy data from buffer to p
	copy(p, b.buffer[b.pos:b.pos+bytesToRead])
	b.pos += bytesToRead

	// Return the number of bytes read
	return bytesToRead, nil
}

// Available returns the number of bytes available to read from the buffer.
// It gives the remaining number of bytes that can be read before reaching the end of the buffer.
//
// Returns:
//   - The number of bytes available to read.
//
// Example:
//
//	available := reader.Available()
//	fmt.Println(available) // Output: number of bytes remaining in the buffer
func (b *ByteBufferInputStream) Available() int {
	if b.pos >= len(b.buffer) {
		return 0
	}
	return len(b.buffer) - b.pos
}
