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
	"errors"
	"strings"
)

// ConsumerGroupState represents the state of a consumer group in a messaging system.
// It is an enumeration of possible states that a consumer group can be in during its lifecycle.
type ConsumerGroupState int

// Define constants for each consumer group state.
// These constants are used to represent the various states of a consumer group.
const (
	Unknown             ConsumerGroupState = iota // The state is unknown.
	PreparingRebalance                            // The consumer group is preparing for a rebalance.
	CompletingRebalance                           // The consumer group is completing a rebalance.
	Stable                                        // The consumer group is in a stable state.
	Dead                                          // The consumer group is dead and not functioning.
	Empty                                         // The consumer group has no active members.
	Assigning                                     // The consumer group is assigning partitions.
	Reconciling                                   // The consumer group is reconciling its state.
)

// String representations of the states.
// This array maps each ConsumerGroupState to its corresponding string representation.
var consumerGroupStateNames = [...]string{
	"Unknown",
	"PreparingRebalance",
	"CompletingRebalance",
	"Stable",
	"Dead",
	"Empty",
	"Assigning",
	"Reconciling",
}

// String returns the string representation of the ConsumerGroupState.
// It returns "InvalidState" if the state is out of range.
func (c ConsumerGroupState) String() string {
	if c < Unknown || c > Reconciling {
		return "InvalidState"
	}
	return consumerGroupStateNames[c]
}

// consumerGroupStateMap for case-insensitive lookups.
// This map allows for easy conversion from string representations to ConsumerGroupState values.
var consumerGroupStateMap = map[string]ConsumerGroupState{
	"UNKNOWN":             Unknown,
	"PREPARINGREBALANCE":  PreparingRebalance,
	"COMPLETINGREBALANCE": CompletingRebalance,
	"STABLE":              Stable,
	"DEAD":                Dead,
	"EMPTY":               Empty,
	"ASSIGNING":           Assigning,
	"RECONCILING":         Reconciling,
}

// ParseConsumerGroupState converts a string to a ConsumerGroupState.
// It returns an error if the provided state name is not recognized.
//
// Example usage:
//
//	state, err := ParseConsumerGroupState("STABLE")
//	if err != nil {
//	    fmt.Println(err)
//	} else {
//	    fmt.Println("Parsed state:", state) // Output: Parsed state: Stable
//	}
func ParseConsumerGroupState(name string) (ConsumerGroupState, error) {
	name = strings.ToUpper(name)
	if state, exists := consumerGroupStateMap[name]; exists {
		return state, nil
	}
	return Unknown, errors.New("unrecognized consumer group state: " + name)
}
