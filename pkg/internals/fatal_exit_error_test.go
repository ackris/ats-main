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
	"fmt"
	"testing"

	ex "github.com/ackris/ats-main/pkg/common" // Ensure this import is present
	"go.uber.org/zap/zaptest"
)

// TestNewFatalExitError tests the NewFatalExitError function.
func TestNewFatalExitError(t *testing.T) {
	tests := []struct {
		statusCode int
		expectErr  bool
	}{
		{0, true},  // Expect an error for status code 0
		{1, false}, // Valid status code
		{2, false}, // Valid status code
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("statusCode=%d", tt.statusCode), func(t *testing.T) {
			err := NewFatalExitError(tt.statusCode)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected an error, got nil")
				}
			} else {
				if err == nil {
					t.Errorf("expected a FatalExitError, got nil")
				} else if _, ok := err.(*FatalExitError); !ok {
					t.Errorf("expected a FatalExitError, got: %T", err)
				}
			}
		})
	}
}

// TestNewDefaultFatalExitError tests the NewDefaultFatalExitError function.
func TestNewDefaultFatalExitError(t *testing.T) {
	err := NewDefaultFatalExitError()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if fatalErr, ok := err.(*FatalExitError); !ok || fatalErr.StatusCode() != 1 {
		t.Fatalf("expected FatalExitError with status code 1, got: %v", err)
	}
}

// TestStatusCode tests the StatusCode method.
func TestStatusCode(t *testing.T) {
	tests := []struct {
		fatalErr *FatalExitError
		expected int
	}{
		{&FatalExitError{statusCode: 42}, 42},
		{nil, 0},
	}

	for _, tt := range tests {
		if code := tt.fatalErr.StatusCode(); code != tt.expected {
			t.Errorf("expected status code %d, got %d", tt.expected, code)
		}
	}
}

// TestUnwrap tests the Unwrap method.
func TestUnwrap(t *testing.T) {
	cause := errors.New("underlying cause")
	atomStateErr := ex.NewAtomStateErrorWithCause(cause) // Create the AtomStateError
	fatalErr := &FatalExitError{
		AtomStateError: atomStateErr.(ex.AtomStateError), // Type assertion to AtomStateError
	}

	if unwrapped := fatalErr.Unwrap(); unwrapped == nil || unwrapped.Error() != "underlying cause" {
		t.Errorf("expected underlying cause to be 'underlying cause', got %v", unwrapped)
	}

	var nilFatalErr *FatalExitError
	if unwrapped := nilFatalErr.Unwrap(); unwrapped != nil {
		t.Error("expected nil for unwrapped nil FatalExitError")
	}
}

// TestLog tests the Log method.
func TestLog(t *testing.T) {
	logger := zaptest.NewLogger(t)
	fatalErr := NewDefaultFatalExitError().(*FatalExitError)

	// Log the error
	fatalErr.Log(logger)

	// No assertion needed; we just ensure it doesn't panic.
}
