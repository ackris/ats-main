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

import "fmt"

// ElectionType represents the type of election.
// It is defined as a byte type, allowing for efficient storage and comparisons.
type ElectionType byte

// Define constants for ElectionType.
// These constants represent the two valid types of elections.
const (
	Preferred ElectionType = iota // 0
	Unclean                       // 1
)

// validElectionTypes holds the valid values for ElectionType.
var validElectionTypes = []ElectionType{Preferred, Unclean}

// String returns a string representation of the ElectionType.
// It returns "Preferred" for Preferred, "Unclean" for Unclean,
// and "Unknown" for any other value.
//
// Example usage:
//
//	var et ElectionType = Preferred
//	fmt.Println(et.String()) // Output: Preferred
func (et ElectionType) String() string {
	switch et {
	case Preferred:
		return "Preferred"
	case Unclean:
		return "Unclean"
	default:
		return "Unknown"
	}
}

// ValueOf converts a byte to an ElectionType.
// It returns the corresponding ElectionType for valid byte values (0 or 1).
// If the byte value is invalid, it returns an error.
//
// Example usage:
//
//	et, err := ValueOf(0)
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//	    fmt.Println("Election Type:", et) // Output: Election Type: Preferred
//	}
func ValueOf(value byte) (ElectionType, error) {
	et := ElectionType(value)
	if et.IsValid() {
		return et, nil
	}
	return 0, fmt.Errorf("value %d must be one of %v", value, validElectionTypes)
}

// IsValid checks if the ElectionType is valid.
// It returns true if the ElectionType is either Preferred or Unclean,
// and false otherwise.
//
// Example usage:
//
//	var et ElectionType = Unclean
//	fmt.Println(et.IsValid()) // Output: true
func (et ElectionType) IsValid() bool {
	return et >= Preferred && et <= Unclean
}
