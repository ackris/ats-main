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

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestNewAtomStateError tests the creation of AtomStateError with message and cause.
func TestNewAtomStateError(t *testing.T) {
	cause := errors.New("underlying cause")
	err := NewAtomStateError("test error", cause)

	assert.NotNil(t, err)
	assert.Equal(t, "test error: underlying cause", err.Error())
	assert.Equal(t, cause, errors.Unwrap(err))
}

// TestNewAtomStateErrorWithMessage tests the creation of AtomStateError with a message only.
func TestNewAtomStateErrorWithMessage(t *testing.T) {
	err := NewAtomStateErrorWithMessage("test error")

	assert.NotNil(t, err)
	assert.Equal(t, "test error", err.Error())
	assert.Nil(t, errors.Unwrap(err))
}

// TestNewAtomStateErrorWithCause tests the creation of AtomStateError with a cause only.
func TestNewAtomStateErrorWithCause(t *testing.T) {
	cause := errors.New("underlying cause")
	err := NewAtomStateErrorWithCause(cause)

	assert.NotNil(t, err)
	assert.Equal(t, "underlying cause", err.Error())
	assert.Equal(t, cause, errors.Unwrap(err))
}

// TestErrorFormatting tests the error formatting with and without a cause.
func TestErrorFormatting(t *testing.T) {
	cause := errors.New("cause")
	errWithCause := AtomStateError{
		Message: "error with cause",
		Cause:   cause,
	}
	errWithoutCause := AtomStateError{
		Message: "error without cause",
	}

	assert.Equal(t, "error with cause: cause", errWithCause.Error())
	assert.Equal(t, "error without cause", errWithoutCause.Error())
}

// TestUnwrap tests the Unwrap method of AtomStateError.
func TestUnwrap(t *testing.T) {
	cause := errors.New("underlying cause")
	err := AtomStateError{
		Message: "error with cause",
		Cause:   cause,
	}

	assert.Equal(t, cause, err.Unwrap())
}

// TestLog tests the Log method with zap logger and verifies output.
func TestLog(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a zapcore encoder configuration
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zap.DebugLevel)

	// Create a zap logger with the custom core
	logger := zap.New(core)

	// Create an AtomStateError
	err := AtomStateError{
		Message: "test log message",
		Cause:   errors.New("log cause"),
	}

	// Log the error
	err.Log(logger)

	// Check if the log contains the expected message and fields
	logOutput := buf.String()
	assert.Contains(t, logOutput, `"msg":"test log message"`)
	assert.Contains(t, logOutput, `"err":"log cause"`)
}
