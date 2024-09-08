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
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected ShareGroupState
	}{
		{"UNKNOWN", UNKNOWN},
		{"Stable", STABLE},
		{"dead", DEAD},
		{"EMPTY", EMPTY},
		{"invalid", UNKNOWN},
		{"", UNKNOWN},
	}

	for _, test := range tests {
		result := Parse(test.input)
		if result != test.expected {
			t.Errorf("Parse(%q) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		state    ShareGroupState
		expected string
	}{
		{UNKNOWN, "Unknown"},
		{STABLE, "Stable"},
		{DEAD, "Dead"},
		{EMPTY, "Empty"},
		{ShareGroupState(100), "Unknown"}, // Testing out-of-range state
	}

	for _, test := range tests {
		result := test.state.String()
		if result != test.expected {
			t.Errorf("String(%v) = %q; expected %q", test.state, result, test.expected)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Parse("Stable")
	}
}

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = STABLE.String()
	}
}
