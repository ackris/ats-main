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

// ByteBufferUnmapper provides functionality to simulate the unmapping of a byte slice.
type ByteBufferUnmapper struct{}

// NewByteBufferUnmapper creates and returns a new ByteBufferUnmapper instance.
//
// Example:
//
//	bbu := NewByteBufferUnmapper()
func NewByteBufferUnmapper() *ByteBufferUnmapper {
	return &ByteBufferUnmapper{}
}

// Unmap simulates the unmapping of the provided byte slice by setting it to nil.
// It returns an error if the buffer is nil or empty.
//
// Parameters:
//
//	resourceDescription (string): A description of the resource associated with the byte slice.
//	buf (*[]byte): A pointer to the byte slice that needs to be unmapped.
//
// Returns:
//
//	error: If an error occurs during the unmapping process.
//
// Example:
//
//	buffer := make([]byte, 1024)
//	err := bbu.Unmap("test buffer", &buffer)
//	if err != nil {
//	  // Handle error
//	}
func (bbu *ByteBufferUnmapper) Unmap(resourceDescription string, buf *[]byte) error {
	if buf == nil {
		return errors.New("buffer is nil")
	}

	if len(*buf) == 0 {
		return errors.New("buffer is empty")
	}

	// Set the buffer to nil to make it eligible for garbage collection.
	*buf = nil

	return nil
}
