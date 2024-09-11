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

func TestUnmap_Success(t *testing.T) {
	bbu := NewByteBufferUnmapper()
	buffer := make([]byte, 1024) // A buffer with some data

	err := bbu.Unmap("test buffer", &buffer)
	if err != nil {
		t.Errorf("Unmap returned an error: %v", err)
	}

	if buffer != nil {
		t.Errorf("Expected buffer to be nil, but got %v", buffer)
	}
}

func TestUnmap_BufferIsNil(t *testing.T) {
	bbu := NewByteBufferUnmapper()
	var buffer *[]byte // buffer is nil

	err := bbu.Unmap("test buffer", buffer)
	if err == nil {
		t.Errorf("Expected an error when buffer is nil, but got nil")
	} else if err.Error() != "buffer is nil" {
		t.Errorf("Expected error 'buffer is nil', but got %v", err)
	}
}

func TestUnmap_BufferIsEmpty(t *testing.T) {
	bbu := NewByteBufferUnmapper()
	emptyBuffer := []byte{} // An empty buffer

	err := bbu.Unmap("test buffer", &emptyBuffer)
	if err == nil {
		t.Errorf("Expected an error when buffer is empty, but got nil")
	} else if err.Error() != "buffer is empty" {
		t.Errorf("Expected error 'buffer is empty', but got %v", err)
	}

	// Check the length of the buffer instead of checking if it's nil
	if len(emptyBuffer) != 0 {
		t.Errorf("Expected buffer to be empty, but got length %d", len(emptyBuffer))
	}
}
