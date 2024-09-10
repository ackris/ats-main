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

package resource

// PatternType represents the type of resource pattern.
type PatternType byte

// Define the PatternType constants.
const (
	// UNKNOWN represents any PatternType which this client cannot understand, perhaps because this client is too old.
	UNKNOWN PatternType = iota // 0

	// ANY represents a filter that matches any resource pattern type.
	ANY // 1

	// MATCH represents a filter that will perform pattern matching.
	// e.g. Given a filter of `ResourcePatternFilter(TOPIC, "payments.received", MATCH)`, the filter match
	// any `ResourcePattern` that matches topic 'payments.received'. This might include:
	// - A Literal pattern with the same type and name, e.g. `ResourcePattern(TOPIC, "payments.received", LITERAL)`
	// - A Wildcard pattern with the same type, e.g. `ResourcePattern(TOPIC, "*", LITERAL)`
	// - A Prefixed pattern with the same type and where the name is a matching prefix, e.g. `ResourcePattern(TOPIC, "payments.", PREFIXED)`
	MATCH // 2

	// LITERAL represents a literal resource name.
	// A literal name defines the full name of a resource, e.g. topic with name 'foo', or group with name 'bob'.
	// The special wildcard character `*` can be used to represent a resource with any name.
	LITERAL // 3

	// PREFIXED represents a prefixed resource name.
	// A prefixed name defines a prefix for a resource, e.g. topics with names that start with 'foo'.
	PREFIXED // 4
)

// patternTypeCount represents the total number of PatternType constants.
// const patternTypeCount = iota

// String returns the string representation of the PatternType.
func (p PatternType) String() string {
	names := [...]string{
		"UNKNOWN",
		"ANY",
		"MATCH",
		"LITERAL",
		"PREFIXED",
	}
	if p < UNKNOWN || p > PREFIXED {
		return "UNKNOWN"
	}
	return names[p]
}

// Code returns the byte code of the PatternType.
func (p PatternType) Code() byte {
	return byte(p)
}

// IsUnknown checks if the PatternType is UNKNOWN.
func (p PatternType) IsUnknown() bool {
	return p == UNKNOWN
}

// IsSpecific checks if the PatternType is a specific type.
func (p PatternType) IsSpecific() bool {
	return p == LITERAL || p == PREFIXED // Only LITERAL and PREFIXED are specific types
}

// PatternTypeFromCode returns the PatternType for the given code, or UNKNOWN if not found.
func PatternTypeFromCode(code byte) PatternType {
	if code >= byte(UNKNOWN) && code <= byte(PREFIXED) {
		return PatternType(code)
	}
	return UNKNOWN
}

// PatternTypeFromString returns the PatternType for the given string, or UNKNOWN if not found.
func PatternTypeFromString(name string) PatternType {
	switch name {
	case "UNKNOWN":
		return UNKNOWN
	case "ANY":
		return ANY
	case "MATCH":
		return MATCH
	case "LITERAL":
		return LITERAL
	case "PREFIXED":
		return PREFIXED
	default:
		return UNKNOWN
	}
}
