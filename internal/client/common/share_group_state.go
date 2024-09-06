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
