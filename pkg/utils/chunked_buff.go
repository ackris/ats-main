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
	"io"
	"sync"
)

// ChunkedBytesStream is a custom buffered reader with configurable behavior.
// It allows reading data in chunks from an underlying io.Reader, with the ability
// to skip bytes and manage buffer memory efficiently.
type ChunkedBytesStream struct {
	r              io.Reader      // The underlying reader to read data from
	bufferSupplier BufferSupplier // Supplier for managing byte buffers
	buffer         []byte         // The buffer holding the data
	count          int            // The number of bytes currently in the buffer
	pos            int            // The current position in the buffer
	delegateSkip   bool           // Flag to indicate if skipping should be delegated to the underlying reader
	mu             sync.Mutex     // Mutex for thread-safe operations
}

// NewChunkedBytesStream creates a new ChunkedBytesStream.
//
// Parameters:
//   - r: The underlying io.Reader from which to read data.
//   - bufferSupplier: A BufferSupplier to manage buffer allocation.
//   - bufferSize: The size of the buffer to be used for reading data.
//   - delegateSkip: A boolean flag indicating whether to delegate skipping to the underlying reader.
//
// Returns:
//   - A pointer to the newly created ChunkedBytesStream.
func NewChunkedBytesStream(r io.Reader, bufferSupplier BufferSupplier, bufferSize int, delegateSkip bool) *ChunkedBytesStream {
	buf := bufferSupplier.Get(bufferSize)
	return &ChunkedBytesStream{
		r:              r,
		bufferSupplier: bufferSupplier,
		buffer:         buf,
		delegateSkip:   delegateSkip,
	}
}

// Read reads from the buffer or the underlying reader if the buffer is exhausted.
//
// Parameters:
//   - p: A byte slice to store the read data.
//
// Returns:
//   - The number of bytes read and any error encountered.
//
// Example usage:
//
//	buf := make([]byte, 1024)
//	n, err := stream.Read(buf)
//	if err != nil {
//	    // Handle error
//	}
//	// Use the data in buf
func (s *ChunkedBytesStream) Read(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.pos >= s.count {
		if err := s.fill(); err != nil {
			if err == io.EOF {
				return 0, io.EOF
			}
			return 0, err
		}
	}

	n := copy(p, s.buffer[s.pos:s.count])
	s.pos += n
	return n, nil
}

// fill reads data into the buffer from the underlying reader.
//
// Returns:
//   - An error if there was an issue reading from the underlying reader.
func (s *ChunkedBytesStream) fill() error {
	s.pos = 0
	s.count = 0
	n, err := s.r.Read(s.buffer)
	if err != nil {
		if err == io.EOF && n == 0 {
			return io.EOF
		}
		return err
	}
	s.count = n
	return nil
}

// Close closes the underlying reader and releases the buffer.
//
// Returns:
//   - An error if there was an issue closing the underlying reader.
func (s *ChunkedBytesStream) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if closer, ok := s.r.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	s.bufferSupplier.Release(s.buffer)
	s.buffer = nil
	return nil
}

// Skip skips the specified number of bytes in the stream.
//
// Parameters:
//   - n: The number of bytes to skip.
//
// Returns:
//   - The total number of bytes actually skipped and any error encountered.
//
// Example usage:
//
//	skipped, err := stream.Skip(10)
//	if err != nil {
//	    // Handle error
//	}
//	fmt.Printf("Skipped %d bytes\n", skipped)
func (s *ChunkedBytesStream) Skip(n int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if n <= 0 {
		return 0, nil
	}

	originalN := n
	remaining := n
	avail := s.count - s.pos

	if int64(avail) >= remaining {
		s.pos += int(remaining)
		return originalN, nil
	}

	// Skip in the current buffer
	s.pos = s.count           // Move to the end of the current buffer
	remaining -= int64(avail) // Decrease remaining by available bytes

	for remaining > 0 {
		if s.delegateSkip {
			// Delegate skipping to the underlying reader
			if seeker, ok := s.r.(io.Seeker); ok {
				skipped, err := seeker.Seek(remaining, io.SeekCurrent)
				if err != nil {
					return originalN - remaining, err
				}
				remaining -= skipped
				if skipped == 0 {
					return originalN - remaining, nil
				}
			} else {
				return originalN - remaining, errors.New("underlying reader does not support Seek")
			}
		} else {
			// Refill the buffer and skip
			if err := s.fill(); err != nil {
				if err == io.EOF && remaining > 0 {
					return originalN - remaining, io.EOF
				}
				return originalN - remaining, err
			}
			avail = s.count - s.pos
			if int64(avail) > remaining {
				s.pos += int(remaining)
				remaining = 0
			} else {
				remaining -= int64(avail)
				s.pos = s.count
			}
		}
	}

	// Return the total number of bytes actually skipped
	return originalN - remaining, nil
}

// Available returns an estimate of the number of bytes that can be read without blocking.
//
// Returns:
//   - The estimated number of available bytes and any error encountered.
//
// Example usage:
//
//	availableBytes, err := stream.Available()
//	if err != nil {
//	    // Handle error
//	}
//	fmt.Printf("Available bytes: %d\n", availableBytes)
func (s *ChunkedBytesStream) Available() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	avail := s.count - s.pos // Calculate available bytes in the buffer
	if s.r == nil {
		return avail, nil
	}

	// If the underlying reader supports seeking, we can check its current position
	if seeker, ok := s.r.(io.Seeker); ok {
		currPos, err := seeker.Seek(0, io.SeekCurrent)
		if err != nil {
			return 0, err
		}

		// Check for potential overflow
		if currPos > int64(int(^uint(0)>>1)) {
			return 0, errors.New("current position exceeds int limit")
		}
		// Return only the bytes available in the buffer
		return avail, nil
	}

	return avail, nil
}
