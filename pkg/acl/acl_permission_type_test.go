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

package acl

import "testing"

// Unit Tests
func TestAPTString(t *testing.T) {
	tests := []struct {
		input    AclPermissionType
		expected string
	}{
		{UNKNOWN, "UNKNOWN"},
		{ANY, "ANY"},
		{DENY, "DENY"},
		{ALLOW, "ALLOW"},
		{maxAclPermissionType, "UNKNOWN"}, // Testing out-of-bounds
	}

	for _, test := range tests {
		result := test.input.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestAPTFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected AclPermissionType
	}{
		{"unknown", UNKNOWN},
		{"ANY", ANY},
		{"deny", DENY},
		{"ALLOW", ALLOW},
		{"invalid", UNKNOWN}, // Testing invalid input
	}

	for _, test := range tests {
		result := FromString(test.input)
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestAPTCode(t *testing.T) {
	tests := []struct {
		input    AclPermissionType
		expected AclPermissionType
	}{
		{UNKNOWN, UNKNOWN},
		{ANY, ANY},
		{DENY, DENY},
		{ALLOW, ALLOW},
	}

	for _, test := range tests {
		result := test.input.Code()
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestAPTIsUnknown(t *testing.T) {
	tests := []struct {
		input    AclPermissionType
		expected bool
	}{
		{UNKNOWN, true},
		{ANY, false},
		{DENY, false},
		{ALLOW, false},
	}

	for _, test := range tests {
		result := test.input.IsUnknown()
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}
