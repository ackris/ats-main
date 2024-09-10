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

import "testing"

// Unit tests
func TestFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected ResourceType
		hasError bool
	}{
		{"UNRECOGNIZED", UNRECOGNIZED, false},
		{"ALL_RESOURCES", ALL_RESOURCES, false},
		{"TOPIC", TOPIC, false},
		{"GROUP", GROUP, false},
		{"CLUSTER", CLUSTER, false},
		{"TRANSACTIONAL_ID", TRANSACTIONAL_ID, false},
		{"DELEGATION_TOKEN", DELEGATION_TOKEN, false},
		{"USER", USER, false},
		{"", UNRECOGNIZED, true},             // empty string, expect error
		{"INVALID_TYPE", UNRECOGNIZED, true}, // invalid type, expect error
	}

	for _, test := range tests {
		result, err := FromString(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("FromString(%q) error = %v, wantError %v", test.input, err, test.hasError)
		}
		if result != test.expected {
			t.Errorf("FromString(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestFromCode(t *testing.T) {
	tests := []struct {
		code     byte
		expected ResourceType
	}{
		{0, UNRECOGNIZED},
		{1, ALL_RESOURCES},
		{2, TOPIC},
		{3, GROUP},
		{4, CLUSTER},
		{5, TRANSACTIONAL_ID},
		{6, DELEGATION_TOKEN},
		{7, USER},
		{8, UNRECOGNIZED}, // out of bounds
	}

	for _, test := range tests {
		result := FromCode(test.code)
		if result != test.expected {
			t.Errorf("FromCode(%d) = %v, want %v", test.code, result, test.expected)
		}
	}
}

func TestRTCode(t *testing.T) {
	tests := []struct {
		resourceType ResourceType
		expectedCode byte
	}{
		{UNRECOGNIZED, 0},
		{ALL_RESOURCES, 1},
		{TOPIC, 2},
		{GROUP, 3},
		{CLUSTER, 4},
		{TRANSACTIONAL_ID, 5},
		{DELEGATION_TOKEN, 6},
		{USER, 7},
	}

	for _, test := range tests {
		result := test.resourceType.Code()
		if result != test.expectedCode {
			t.Errorf("Code() for %v = %d, want %d", test.resourceType, result, test.expectedCode)
		}
	}
}

func TestIsUnrecognized(t *testing.T) {
	tests := []struct {
		resourceType ResourceType
		expected     bool
	}{
		{UNRECOGNIZED, true},
		{ALL_RESOURCES, false},
		{TOPIC, false},
		{GROUP, false},
		{CLUSTER, false},
		{TRANSACTIONAL_ID, false},
		{DELEGATION_TOKEN, false},
		{USER, false},
	}

	for _, test := range tests {
		result := test.resourceType.IsUnrecognized()
		if result != test.expected {
			t.Errorf("IsUnrecognized() for %v = %v, want %v", test.resourceType, result, test.expected)
		}
	}
}
