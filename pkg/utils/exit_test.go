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
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
)

// CaptureOutput captures the output of a function that writes to stdout.
func CaptureOutput(f func()) string {
	// Create a temporary file to redirect stdout
	tempFile, err := os.CreateTemp("", "stdout")
	if err != nil {
		panic(err) // Handle error appropriately in production code
	}
	defer os.Remove(tempFile.Name()) // Clean up the temp file

	original := os.Stdout                   // Keep backup of the real stdout
	defer func() { os.Stdout = original }() // Restore original stdout
	os.Stdout = tempFile                    // Redirect stdout to the temp file

	f() // Call the function that writes to stdout

	// Restore original stdout
	os.Stdout = original

	// Read the output from the temp file
	tempFile.Seek(0, 0) // Go to the beginning of the file
	output, err := io.ReadAll(tempFile)
	if err != nil {
		panic(err) // Handle error appropriately in production code
	}

	return string(output) // Return the captured output
}

// TestExit tests the Exit function.
func TestExit(t *testing.T) {
	// Capture the output
	exitProcedure = func(code int, message string) {
		fmt.Printf("Exit called with code: %d, message: %s\n", code, message)
	}

	expectedMessage := "Test exit message"
	output := CaptureOutput(func() {
		Exit(1, expectedMessage)
	})

	expectedOutput := fmt.Sprintf("Exit called with code: %d, message: %s\n", 1, expectedMessage)
	if output != expectedOutput {
		t.Errorf("Unexpected output: got %v, want %v", output, expectedOutput)
	}
}

// TestHalt tests the Halt function.
func TestHalt(t *testing.T) {
	// Capture the output
	haltProcedure = func(code int, message string) {
		fmt.Printf("Halt called with code: %d, message: %s\n", code, message)
	}

	expectedMessage := "Test halt message"
	output := CaptureOutput(func() {
		Halt(2, expectedMessage)
	})

	expectedOutput := fmt.Sprintf("Halt called with code: %d, message: %s\n", 2, expectedMessage)
	if output != expectedOutput {
		t.Errorf("Unexpected output: got %v, want %v", output, expectedOutput)
	}
}

// TestAddShutdownHook tests the AddShutdownHook function.
func TestAddShutdownHook(t *testing.T) {
	var mu sync.Mutex
	called := false

	// Save the original hook adder
	originalAdder := hookAdder
	defer func() { hookAdder = originalAdder }() // Restore original adder

	// Override the hookAdder to directly call the function for testing
	hookAdder = func(name string, f func()) {
		f() // Call the function immediately
	}

	// Add a shutdown hook
	AddShutdownHook("testHook", func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})

	// Trigger the shutdown hook
	hookAdder("testHook", func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})

	// Check if the shutdown hook was called
	mu.Lock()
	defer mu.Unlock()
	if !called {
		t.Error("Shutdown hook was not called")
	}
}

// TestSetExitProcedure tests the SetExitProcedure function.
func TestSetExitProcedure(t *testing.T) {
	originalProcedure := exitProcedure
	defer func() { exitProcedure = originalProcedure }() // Restore original exit procedure

	newProcedure := func(code int, message string) {
		fmt.Printf("New exit called with code: %d, message: %s\n", code, message)
	}
	SetExitProcedure(newProcedure)

	// Call the new procedure to verify it works
	output := CaptureOutput(func() {
		exitProcedure(1, "Test message")
	})

	expectedOutput := "New exit called with code: 1, message: Test message\n"
	if output != expectedOutput {
		t.Errorf("Unexpected output: got %v, want %v", output, expectedOutput)
	}
}

// TestSetHaltProcedure tests the SetHaltProcedure function.
func TestSetHaltProcedure(t *testing.T) {
	originalProcedure := haltProcedure
	defer func() { haltProcedure = originalProcedure }() // Restore original halt procedure

	newProcedure := func(code int, message string) {
		fmt.Printf("New halt called with code: %d, message: %s\n", code, message)
	}
	SetHaltProcedure(newProcedure)

	// Call the new procedure to verify it works
	output := CaptureOutput(func() {
		haltProcedure(2, "Test message")
	})

	expectedOutput := "New halt called with code: 2, message: Test message\n"
	if output != expectedOutput {
		t.Errorf("Unexpected output: got %v, want %v", output, expectedOutput)
	}
}

// TestResetExitProcedure tests the ResetExitProcedure function.
func TestResetExitProcedure(t *testing.T) {
	// Save the original exit procedure
	originalProcedure := exitProcedure
	defer func() { exitProcedure = originalProcedure }() // Restore original exit procedure

	// Set to a no-op
	SetExitProcedure(func(code int, message string) {
		// No-op
	})
	ResetExitProcedure()

	// Check if the exit procedure is reset to the original behavior
	output := CaptureOutput(func() {
		originalProcedure(1, "Test message")
	})

	expectedOutput := fmt.Sprintf("Exit called with code: %d, message: %s\n", 1, "Test message")
	if output != expectedOutput {
		t.Error("Exit procedure was not reset correctly")
	}
}

// TestResetHaltProcedure tests the ResetHaltProcedure function.
func TestResetHaltProcedure(t *testing.T) {
	// Save the original halt procedure
	originalProcedure := haltProcedure
	defer func() { haltProcedure = originalProcedure }() // Restore original halt procedure

	// Set to a no-op
	SetHaltProcedure(func(code int, message string) {
		// No-op
	})
	ResetHaltProcedure()

	// Check if the halt procedure is reset to the original behavior
	output := CaptureOutput(func() {
		originalProcedure(2, "Test message")
	})

	expectedOutput := fmt.Sprintf("Halt called with code: %d, message: %s\n", 2, "Test message")
	if output != expectedOutput {
		t.Error("Halt procedure was not reset correctly")
	}
}

// TestSetShutdownHookAdder tests the SetShutdownHookAdder function.
func TestSetShutdownHookAdder(t *testing.T) {
	originalAdder := hookAdder
	defer func() { hookAdder = originalAdder }() // Restore original adder

	newAdder := func(name string, f func()) {
		// Call the function directly for testing
		f()
	}
	SetShutdownHookAdder(newAdder)

	// Ensure the adder is set correctly
	called := false
	AddShutdownHook("testHook", func() {
		called = true
	})

	// Trigger the shutdown hook
	hookAdder("testHook", func() {
		called = true
	})

	if !called {
		t.Error("Shutdown hook was not called")
	}
}
