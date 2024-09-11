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
	"testing"
)

func TestRead(t *testing.T) {
	buffer := []byte("Hello, World!")
	reader := NewByteBufferInputStream(buffer)

	// Test reading fewer bytes than available
	p := make([]byte, 5)
	n, err := reader.Read(p)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if got := string(p[:n]); got != "Hello" {
		t.Errorf("Read() got = %v, want %v", got, "Hello")
	}

	// Test reading exactly the remaining bytes
	p = make([]byte, 8)
	n, err = reader.Read(p)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if got := string(p[:n]); got != ", World!" {
		t.Errorf("Read() got = %v, want %v", got, ", World!")
	}

	// Test reading beyond the buffer
	p = make([]byte, 10)
	n, err = reader.Read(p)
	if err != io.EOF {
		t.Errorf("Read() error = %v, want %v", err, io.EOF)
	}
	if got := string(p[:n]); got != "" {
		t.Errorf("Read() got = %v, want %v", got, "")
	}
}

func TestAvailable(t *testing.T) {
	buffer := []byte("Hello, World!")
	reader := NewByteBufferInputStream(buffer)

	if got := reader.Available(); got != 13 {
		t.Errorf("Available() = %v, want %v", got, 13)
	}

	// Read 5 bytes
	p := make([]byte, 5)
	_, err := reader.Read(p)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if got := reader.Available(); got != 8 {
		t.Errorf("Available() = %v, want %v", got, 8)
	}

	// Read remaining bytes
	p = make([]byte, 8)
	_, err = reader.Read(p)
	if err != nil && err != io.EOF {
		t.Errorf("Read() error = %v, want %v", err, io.EOF)
	}

	if got := reader.Available(); got != 0 {
		t.Errorf("Available() = %v, want %v", got, 0)
	}
}

func TestReadEmptyBuffer(t *testing.T) {
	buffer := []byte("")
	reader := NewByteBufferInputStream(buffer)

	p := make([]byte, 5)
	n, err := reader.Read(p)
	if err != io.EOF {
		t.Errorf("Read() error = %v, want %v", err, io.EOF)
	}
	if n != 0 {
		t.Errorf("Read() n = %v, want %v", n, 0)
	}

	if got := reader.Available(); got != 0 {
		t.Errorf("Available() = %v, want %v", got, 0)
	}
}

func TestReadWithBufferLargerThanAvailable(t *testing.T) {
	buffer := []byte("Short")
	reader := NewByteBufferInputStream(buffer)

	p := make([]byte, 10)
	n, err := reader.Read(p)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if got := string(p[:n]); got != "Short" {
		t.Errorf("Read() got = %v, want %v", got, "Short")
	}

	if got := reader.Available(); got != 0 {
		t.Errorf("Available() = %v, want %v", got, 0)
	}
}
