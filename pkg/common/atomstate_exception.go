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

package common

import "go.uber.org/zap"

// AtomStateError represents an error in the AtomState library.
type AtomStateError struct {
	Message string
	Cause   error
}

// Error implements the error interface for AtomStateError.
// It formats the error message including the cause if present.
func (e AtomStateError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap returns the underlying cause of the error, if any.
// This method is part of Go's error wrapping introduced in Go 1.13.
// It allows error unwrapping to retrieve the original cause of the error.
func (e AtomStateError) Unwrap() error {
	return e.Cause
}

// Log logs the error using the provided logger.
// It logs the error message and the cause if present, using the zap.Logger.
// Example usage:
//
//	logger := zap.NewExample() // Create a zap logger
//	err := NewAtomStateError("an error occurred", errors.New("underlying cause"))
//	err.Log(logger) // Logs: "an error occurred: underlying cause"
func (e AtomStateError) Log(logger *zap.Logger) {
	if e.Cause != nil {
		logger.Error(e.Message, zap.Error(e.Cause))
	} else {
		logger.Error(e.Message)
	}
}

// NewAtomStateError creates a new AtomStateError with the given message and cause.
// This function returns an AtomStateError with the specified message and cause.
// Example usage:
//
//	cause := errors.New("file not found")
//	err := NewAtomStateError("failed to open file", cause)
//	fmt.Println(err.Error()) // Output: "failed to open file: file not found"
func NewAtomStateError(message string, cause error) error {
	return AtomStateError{
		Message: message,
		Cause:   cause,
	}
}

// NewAtomStateErrorWithMessage creates a new AtomStateError with the given message.
// This function returns an AtomStateError with the specified message and no cause.
// Example usage:
//
//	err := NewAtomStateErrorWithMessage("connection lost")
//	fmt.Println(err.Error()) // Output: "connection lost"
func NewAtomStateErrorWithMessage(message string) error {
	return AtomStateError{
		Message: message,
	}
}

// NewAtomStateErrorWithCause creates a new AtomStateError with the given cause.
// This function returns an AtomStateError with the specified cause and no message.
// Example usage:
//
//	cause := errors.New("timeout occurred")
//	err := NewAtomStateErrorWithCause(cause)
//	fmt.Println(err.Error()) // Output: "timeout occurred"
func NewAtomStateErrorWithCause(cause error) error {
	return AtomStateError{
		Cause: cause,
	}
}
