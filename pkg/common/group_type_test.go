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

import "testing"

func TestGroupType_String(t *testing.T) {
	tests := []struct {
		groupType GroupType
		expected  string
	}{
		{UNSPECIFIED, "Unspecified"},
		{CONSUMER, "Consumer"},
		{CLASSIC, "Classic"},
		{SHARE, "Share"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if got := test.groupType.String(); got != test.expected {
				t.Errorf("GroupType.String() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestParseGroupType(t *testing.T) {
	tests := []struct {
		name     string
		expected GroupType
	}{
		{"", UNSPECIFIED},            // Empty string
		{"unknown", UNSPECIFIED},     // Unrecognized input
		{"UNSPECIFIED", UNSPECIFIED}, // Case insensitive
		{"consumer", CONSUMER},       // Valid input
		{"Consumer", CONSUMER},       // Case insensitive valid input
		{"CLASSIC", CLASSIC},         // Case insensitive valid input
		{"share", SHARE},             // Valid input
		{"Share", SHARE},             // Case insensitive valid input
		{"invalid", UNSPECIFIED},     // Invalid input
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseGroupType(test.name); got != test.expected {
				t.Errorf("ParseGroupType(%q) = %v, want %v", test.name, got, test.expected)
			}
		})
	}
}
