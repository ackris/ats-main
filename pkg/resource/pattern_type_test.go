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

import (
	"testing"
)

// TestString tests the String() method of PatternType.
func TestString(t *testing.T) {
	tests := []struct {
		input    PatternType
		expected string
	}{
		{UNKNOWN, "UNKNOWN"},
		{ANY, "ANY"},
		{MATCH, "MATCH"},
		{LITERAL, "LITERAL"},
		{PREFIXED, "PREFIXED"},
	}

	for _, test := range tests {
		if got := test.input.String(); got != test.expected {
			t.Errorf("PatternType(%d).String() = %v; want %v", test.input, got, test.expected)
		}
	}

	// Test an out-of-bounds value
	var invalidPatternType PatternType = 5
	if got := invalidPatternType.String(); got != "UNKNOWN" {
		t.Errorf("PatternType(5).String() = %v; want UNKNOWN", got)
	}
}

// TestCode tests the Code() method of PatternType.
func TestCode(t *testing.T) {
	tests := []struct {
		input    PatternType
		expected byte
	}{
		{UNKNOWN, 0},
		{ANY, 1},
		{MATCH, 2},
		{LITERAL, 3},
		{PREFIXED, 4},
	}

	for _, test := range tests {
		if got := test.input.Code(); got != test.expected {
			t.Errorf("PatternType(%d).Code() = %v; want %v", test.input, got, test.expected)
		}
	}
}

// TestIsUnknown tests the IsUnknown() method of PatternType.
func TestIsUnknown(t *testing.T) {
	tests := []struct {
		input    PatternType
		expected bool
	}{
		{UNKNOWN, true},
		{ANY, false},
		{MATCH, false},
		{LITERAL, false},
		{PREFIXED, false},
	}

	for _, test := range tests {
		if got := test.input.IsUnknown(); got != test.expected {
			t.Errorf("PatternType(%d).IsUnknown() = %v; want %v", test.input, got, test.expected)
		}
	}
}

// TestIsSpecific tests the IsSpecific() method of PatternType.
func TestIsSpecific(t *testing.T) {
	tests := []struct {
		input    PatternType
		expected bool
	}{
		{UNKNOWN, false},
		{ANY, false},
		{MATCH, false},
		{LITERAL, true},
		{PREFIXED, true},
	}

	for _, test := range tests {
		if got := test.input.IsSpecific(); got != test.expected {
			t.Errorf("PatternType(%d).IsSpecific() = %v; want %v", test.input, got, test.expected)
		}
	}
}

// TestPatternTypeFromCode tests the PatternTypeFromCode function.
func TestPatternTypeFromCode(t *testing.T) {
	tests := []struct {
		input    byte
		expected PatternType
	}{
		{0, UNKNOWN},
		{1, ANY},
		{2, MATCH},
		{3, LITERAL},
		{4, PREFIXED},
		{5, UNKNOWN}, // Out of bounds
	}

	for _, test := range tests {
		if got := PatternTypeFromCode(test.input); got != test.expected {
			t.Errorf("PatternTypeFromCode(%d) = %v; want %v", test.input, got, test.expected)
		}
	}
}

// TestPatternTypeFromString tests the PatternTypeFromString function.
func TestPatternTypeFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected PatternType
	}{
		{"UNKNOWN", UNKNOWN},
		{"ANY", ANY},
		{"MATCH", MATCH},
		{"LITERAL", LITERAL},
		{"PREFIXED", PREFIXED},
		{"INVALID", UNKNOWN}, // Invalid input
	}

	for _, test := range tests {
		if got := PatternTypeFromString(test.input); got != test.expected {
			t.Errorf("PatternTypeFromString(%q) = %v; want %v", test.input, got, test.expected)
		}
	}
}
