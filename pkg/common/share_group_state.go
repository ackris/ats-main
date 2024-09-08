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
	"strings"
)

// ShareGroupState represents the state of a share group.
type ShareGroupState int

const (
	// UNKNOWN represents an unknown state.
	UNKNOWN ShareGroupState = iota
	// STABLE represents a stable state.
	STABLE
	// DEAD represents a dead state.
	DEAD
	// EMPTY represents an empty state.
	EMPTY
	stateCount // Used for length of stateNames and nameToEnum
)

var stateNames = [...]string{
	"Unknown",
	"Stable",
	"Dead",
	"Empty",
}

var nameToEnum = map[string]ShareGroupState{
	"UNKNOWN": UNKNOWN,
	"STABLE":  STABLE,
	"DEAD":    DEAD,
	"EMPTY":   EMPTY,
}

// Parse converts a string name to a ShareGroupState, returning UNKNOWN if not found.
//
// Example:
//
//	state := Parse("Stable")
//	fmt.Println(state) // Output: Stable
//
//	state = Parse("invalid")
//	fmt.Println(state) // Output: Unknown
func Parse(name string) ShareGroupState {
	if state, exists := nameToEnum[strings.ToUpper(name)]; exists {
		return state
	}
	return UNKNOWN
}

// String returns the string representation of the ShareGroupState.
//
// Example:
//
//	state := STABLE
//	fmt.Println(state.String()) // Output: Stable
//
//	state = ShareGroupState(100)
//	fmt.Println(state.String()) // Output: Unknown
func (s ShareGroupState) String() string {
	if int(s) >= 0 && int(s) < len(stateNames) {
		return stateNames[s]
	}
	return "Unknown"
}
