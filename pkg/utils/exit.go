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
	"os"
	"sync"
)

// ExitProcedure defines a function type for exit procedures.
type ExitProcedure func(int, string)

// ShutdownHookAdder defines a function type for adding shutdown hooks.
type ShutdownHookAdder func(string, func())

var (
	defaultExitProcedure     = func(statusCode int, message string) { os.Exit(statusCode) }
	defaultShutdownHookAdder = addShutdownHook

	exitProcedure ExitProcedure     = defaultExitProcedure
	haltProcedure ExitProcedure     = defaultExitProcedure // Using the same as exit for simplicity
	hookAdder     ShutdownHookAdder = defaultShutdownHookAdder
	mu            sync.Mutex
)

// Exit performs a graceful exit with an optional message.
//
// Example usage:
//
//	Exit(1, "Exiting with error")
//	Exit(0) // Exit with success status
func Exit(statusCode int, message ...string) {
	mu.Lock()
	defer mu.Unlock()
	exitProcedure(statusCode, formatMessage(message...))
}

// Halt performs an immediate halt with an optional message.
//
// Example usage:
//
//	Halt(1, "Halting with error")
//	Halt(0) // Halt with success status
func Halt(statusCode int, message ...string) {
	mu.Lock()
	defer mu.Unlock()
	haltProcedure(statusCode, formatMessage(message...))
}

// AddShutdownHook registers a shutdown hook.
//
// Example usage:
//
//	AddShutdownHook("cleanup", func() {
//	    // Cleanup resources
//	})
func AddShutdownHook(name string, f func()) {
	mu.Lock()
	defer mu.Unlock()
	hookAdder(name, f)
}

// SetExitProcedure allows setting a custom exit procedure for testing.
//
// Example usage:
//
//	originalProcedure := exitProcedure
//	defer func() { exitProcedure = originalProcedure }() // Restore original procedure
//
//	SetExitProcedure(func(code int, message string) {
//	    // Custom exit procedure
//	})
func SetExitProcedure(procedure ExitProcedure) {
	mu.Lock()
	defer mu.Unlock()
	if procedure == nil {
		procedure = defaultExitProcedure // Fallback to default if nil
	}
	exitProcedure = procedure
}

// SetHaltProcedure allows setting a custom halt procedure for testing.
//
// Example usage:
//
//	originalProcedure := haltProcedure
//	defer func() { haltProcedure = originalProcedure }() // Restore original procedure
//
//	SetHaltProcedure(func(code int, message string) {
//	    // Custom halt procedure
//	})
func SetHaltProcedure(procedure ExitProcedure) {
	mu.Lock()
	defer mu.Unlock()
	if procedure == nil {
		procedure = defaultExitProcedure // Fallback to default if nil
	}
	haltProcedure = procedure
}

// SetShutdownHookAdder allows setting a custom shutdown hook adder for testing.
//
// Example usage:
//
//	originalAdder := hookAdder
//	defer func() { hookAdder = originalAdder }() // Restore original adder
//
//	SetShutdownHookAdder(func(name string, f func()) {
//	    // Custom shutdown hook adder
//	})
func SetShutdownHookAdder(adder ShutdownHookAdder) {
	mu.Lock()
	defer mu.Unlock()
	if adder == nil {
		adder = defaultShutdownHookAdder // Fallback to default if nil
	}
	hookAdder = adder
}

// ResetExitProcedure resets the exit procedure to a no-op.
//
// Example usage:
//
//	ResetExitProcedure()
//	Exit(1, "This will be a no-op")
func ResetExitProcedure() {
	mu.Lock()
	defer mu.Unlock()
	exitProcedure = noopExitProcedure
}

// ResetHaltProcedure resets the halt procedure to a no-op.
//
// Example usage:
//
//	ResetHaltProcedure()
//	Halt(1, "This will be a no-op")
func ResetHaltProcedure() {
	mu.Lock()
	defer mu.Unlock()
	haltProcedure = noopHaltProcedure
}

// ResetShutdownHookAdder resets the shutdown hook adder to the default.
//
// Example usage:
//
//	ResetShutdownHookAdder()
//	AddShutdownHook("cleanup", func() {
//	    // Cleanup resources
//	})
func ResetShutdownHookAdder() {
	mu.Lock()
	defer mu.Unlock()
	hookAdder = defaultShutdownHookAdder
}

// formatMessage formats the message for exit or halt.
func formatMessage(message ...string) string {
	if len(message) > 0 {
		return message[0]
	}
	return ""
}

// addShutdownHook adds a shutdown hook that runs in a separate goroutine.
func addShutdownHook(name string, f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Shutdown hook '%s' panicked: %v\n", name, r)
			}
			if name != "" {
				fmt.Printf("Shutdown hook '%s' completed\n", name)
			}
		}()
		fmt.Printf("Executing shutdown hook '%s'...\n", name)
		f()
	}()
}

// noopExitProcedure is a no-op exit procedure for testing.
var noopExitProcedure ExitProcedure = func(statusCode int, message string) {
	fmt.Println("No-op exit called. Message:", message)
}

// noopHaltProcedure is a no-op halt procedure for testing.
var noopHaltProcedure ExitProcedure = func(statusCode int, message string) {
	fmt.Println("No-op halt called. Message:", message)
}
