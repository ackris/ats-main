package common

import (
	"strings"
)

// ShareGroupState represents the state of a share group.
type ShareGroupState int

const (
	UNKNOWN ShareGroupState = iota
	STABLE
	DEAD
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
func Parse(name string) ShareGroupState {
	if state, exists := nameToEnum[strings.ToUpper(name)]; exists {
		return state
	}
	return UNKNOWN
}

// String returns the string representation of the ShareGroupState.
func (s ShareGroupState) String() string {
	if int(s) >= 0 && int(s) < len(stateNames) {
		return stateNames[s]
	}
	return "Unknown"
}
