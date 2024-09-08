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

package internals

import (
	"errors"
	"strconv"
	"sync"
)

// IdempotentCloser is a struct that represents a resource that can be closed idempotently.
// It ensures thread-safety by using a mutex to protect the closed state.
type IdempotentCloser struct {
	mu       sync.Mutex // Mutex to ensure thread-safe operations
	isClosed bool       // Indicates whether the resource is closed
}

// NewIdempotentCloser creates a new IdempotentCloser that is not yet closed.
//
// Example usage:
//
//	closer := NewIdempotentCloser()
func NewIdempotentCloser() *IdempotentCloser {
	return &IdempotentCloser{}
}

// AssertOpen checks if the IdempotentCloser is still open. If it is closed, it returns an error with the provided message.
//
// Example usage:
//
//	err := closer.AssertOpen("Resource is closed")
//	if err != nil {
//		fmt.Println(err)
//	}
func (ic *IdempotentCloser) AssertOpen(message string) error {
	ic.mu.Lock()
	defer ic.mu.Unlock()

	if ic.isClosed {
		return errors.New(message)
	}
	return nil
}

// IsClosed returns whether the IdempotentCloser is closed.
//
// Example usage:
//
//	if closer.IsClosed() {
//		fmt.Println("Resource is closed")
//	}
func (ic *IdempotentCloser) IsClosed() bool {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	return ic.isClosed
}

// Close closes the resource in a thread-safe manner.
// It accepts two optional functions: onInitialClose to be executed when the resource is initially closed,
// and onSubsequentClose to be executed if the resource was already closed.
//
// Example usage:
//
//	closer.Close(func() {
//		fmt.Println("Resource closed for the first time.")
//	}, func() {
//		fmt.Println("Resource was already closed.")
//	})
func (ic *IdempotentCloser) Close(onInitialClose func(), onSubsequentClose func()) {
	ic.mu.Lock()
	defer ic.mu.Unlock()

	if !ic.isClosed {
		ic.isClosed = true
		if onInitialClose != nil {
			onInitialClose()
		}
	} else {
		if onSubsequentClose != nil {
			onSubsequentClose()
		}
	}
}

// String returns a string representation of the IdempotentCloser.
func (ic *IdempotentCloser) String() string {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	return "IdempotentCloser{isClosed: " + strconv.FormatBool(ic.isClosed) + "}"
}
