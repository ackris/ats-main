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
	"fmt"

	ex "github.com/ackris/ats-main/pkg/common"
	"go.uber.org/zap"
)

// FatalExitError represents a fatal error that should cause the program to exit.
// It embeds AtomStateError to provide additional error context and a status code.
//
// Example usage:
//
//	func main() {
//	    err := NewFatalExitError(2)
//	    if err != nil {
//	        logger := zap.NewExample() // Create a logger
//	        fatalErr := err.(*FatalExitError)
//	        fatalErr.Log(logger) // Log the fatal error
//	        os.Exit(fatalErr.StatusCode()) // Exit the program with the status code
//	    }
//	}
type FatalExitError struct {
	ex.AtomStateError
	statusCode int
}

// NewFatalExitError creates a new FatalExitError with the given status code.
// If the status code is 0, it returns an AtomStateError indicating the issue.
//
// Parameters:
//   - statusCode: An integer representing the exit status code. Must not be 0.
//
// Returns:
//   - An error of type FatalExitError if the statusCode is valid, or an AtomStateError if statusCode is 0.
//
// Example usage:
//
//	err := NewFatalExitError(1)
//	if err != nil {
//	    fmt.Println(err.Error()) // Output: "Fatal exit error with status code 1"
//	}
func NewFatalExitError(statusCode int) error {
	if statusCode == 0 {
		return ex.NewAtomStateErrorWithMessage("statusCode must not be 0")
	}
	return &FatalExitError{
		AtomStateError: ex.AtomStateError{
			Message: fmt.Sprintf("Fatal exit error with status code %d", statusCode),
		},
		statusCode: statusCode,
	}
}

// NewDefaultFatalExitError creates a new FatalExitError with a default status code of 1.
//
// Returns:
//   - An error of type FatalExitError with status code 1.
//
// Example usage:
//
//	err := NewDefaultFatalExitError()
//	fmt.Println(err.Error()) // Output: "Fatal exit error with status code 1"
func NewDefaultFatalExitError() error {
	return NewFatalExitError(1)
}

// StatusCode returns the status code of the FatalExitError.
//
// Returns:
//   - An integer representing the status code. If the receiver is nil, returns 0.
//
// Example usage:
//
//	fatalErr := NewDefaultFatalExitError().(*FatalExitError)
//	fmt.Println(fatalErr.StatusCode()) // Output: 1
func (e *FatalExitError) StatusCode() int {
	if e == nil {
		return 0
	}
	return e.statusCode
}

// Unwrap returns the underlying cause of the error, if any.
// This method is part of Go's error wrapping introduced in Go 1.13.
// It allows error unwrapping to retrieve the original cause of the error.
//
// Returns:
//   - An error representing the underlying cause, or nil if there is none.
//
// Example usage:
//
//	cause := errors.New("underlying cause")
//	fatalErr := &FatalExitError{
//	    AtomStateError: ex.NewAtomStateErrorWithCause(cause).(ex.AtomStateError),
//	}
//	fmt.Println(fatalErr.Unwrap()) // Output: "underlying cause"
func (e *FatalExitError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// Log logs the error using the provided logger.
// It logs the error message and the cause if present, using the zap.Logger.
//
// Parameters:
//   - logger: A pointer to a zap.Logger instance used for logging.
//
// Example usage:
//
//	logger := zap.NewExample() // Create a zap logger
//	fatalErr := NewDefaultFatalExitError().(*FatalExitError)
//	fatalErr.Log(logger) // Logs the fatal error message
func (e *FatalExitError) Log(logger *zap.Logger) {
	if e == nil {
		return
	}
	e.AtomStateError.Log(logger)
}
