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
	"io"
	"testing"
)

func TestChunkedBytesStream(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		readSize       int
		skipSize       int64
		expectedOutput []byte
		expectedSkip   int64
	}{
		{
			name:           "Read and Skip Test",
			input:          []byte("Hello, world!"),
			readSize:       5,
			skipSize:       7,
			expectedOutput: []byte("world"),
			expectedSkip:   7,
		},
		{
			name:           "Read Entire Buffer",
			input:          []byte("Hello"),
			readSize:       5,
			skipSize:       0,
			expectedOutput: []byte("Hello"),
			expectedSkip:   0,
		},
		{
			name:           "Skip and Read",
			input:          []byte("Hello, world!"),
			readSize:       6,
			skipSize:       7,
			expectedOutput: []byte("world!"),
			expectedSkip:   7,
		},
		{
			name:           "Read After EOF",
			input:          []byte("Hello"),
			readSize:       10,
			skipSize:       0,
			expectedOutput: []byte("Hello"),
			expectedSkip:   0,
		},
		{
			name:           "Skip More Than Available",
			input:          []byte("Hello"),
			readSize:       0,
			skipSize:       10,
			expectedOutput: []byte(""),
			expectedSkip:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSupplier := &MockBufferSupplier{}
			stream := NewChunkedBytesStream(bytes.NewReader(tt.input), mockSupplier, len(tt.input), true)
			defer stream.Close()

			// Test Skip
			skipped, err := stream.Skip(tt.skipSize)
			if err != nil {
				t.Fatalf("unexpected error during Skip: %v", err)
			}
			if skipped != tt.expectedSkip {
				t.Errorf("expected %d bytes skipped, got %d", tt.expectedSkip, skipped)
			}

			// Test Read
			buf := make([]byte, tt.readSize)
			n, err := stream.Read(buf)
			if err != nil && err != io.EOF {
				t.Fatalf("unexpected error during Read: %v", err)
			}
			if n != len(tt.expectedOutput) {
				t.Errorf("expected %d bytes read, got %d", len(tt.expectedOutput), n)
			}
			if !bytes.Equal(buf[:n], tt.expectedOutput) {
				t.Errorf("expected %q, got %q", tt.expectedOutput, buf[:n])
			}
		})
	}

	// Additional test for Available method
	t.Run("Available Method Test", func(t *testing.T) {
		mockSupplier := &MockBufferSupplier{}
		stream := NewChunkedBytesStream(bytes.NewReader([]byte("Hello, world!")), mockSupplier, 10, true)
		defer stream.Close()

		// Initially, we should have 0 bytes available
		available, err := stream.Available()
		if err != nil {
			t.Fatalf("unexpected error during Available: %v", err)
		}
		if available != 0 {
			t.Errorf("expected 0 bytes available, got %d", available)
		}

		// Read some data to fill the buffer
		buf := make([]byte, 5)
		_, err = stream.Read(buf)
		if err != nil {
			t.Fatalf("unexpected error during Read: %v", err)
		}

		// Now we should have some bytes available
		available, err = stream.Available()
		if err != nil {
			t.Fatalf("unexpected error during Available: %v", err)
		}
		if available != 5 {
			t.Errorf("expected 5 bytes available, got %d", available)
		}

		// Skip some bytes
		_, err = stream.Skip(3)
		if err != nil {
			t.Fatalf("unexpected error during Skip: %v", err)
		}

		// Check available bytes again
		available, err = stream.Available()
		if err != nil {
			t.Fatalf("unexpected error during Available: %v", err)
		}
		if available != 2 {
			t.Errorf("expected 2 bytes available after skipping, got %d", available)
		}
	})
}
