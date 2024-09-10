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

// TestNewResourcePattern tests the NewResourcePattern function.
func TestNewResourcePattern(t *testing.T) {
	tests := []struct {
		resourceType ResourceType
		name         string
		patternType  PatternType
		expectError  bool
	}{
		{TOPIC, "test-topic", LITERAL, false},
		{GROUP, "test-group", PREFIXED, false},
		{UNRECOGNIZED, "test-name", LITERAL, true}, // Expect error for UNRECOGNIZED
		{TOPIC, "test-topic", UNKNOWN, true},       // Expect error for UNKNOWN
	}

	for _, tt := range tests {
		rp, err := NewResourcePattern(tt.resourceType, tt.name, tt.patternType)
		if (err != nil) != tt.expectError {
			t.Errorf("NewResourcePattern(%v, %q, %v) error = %v, expectError %v", tt.resourceType, tt.name, tt.patternType, err, tt.expectError)
		}
		if !tt.expectError && rp == nil {
			t.Errorf("Expected ResourcePattern but got nil")
		}
	}
}

// TestResourcePatternMethods tests the methods of ResourcePattern.
func TestResourcePatternMethods(t *testing.T) {
	rp, err := NewResourcePattern(TOPIC, "test-topic", LITERAL)
	if err != nil {
		t.Fatalf("Failed to create ResourcePattern: %v", err)
	}

	if got := rp.ResourceType(); got != TOPIC {
		t.Errorf("ResourceType() = %v, want %v", got, TOPIC)
	}

	if got := rp.Name(); got != "test-topic" {
		t.Errorf("Name() = %q, want %q", got, "test-topic")
	}

	if got := rp.PatternType(); got != LITERAL {
		t.Errorf("PatternType() = %v, want %v", got, LITERAL)
	}

	filter, err := rp.ToFilter()
	if err != nil {
		t.Fatalf("ToFilter() error = %v", err)
	}
	if filter == nil {
		t.Error("Expected ResourcePatternFilter but got nil")
	}

	if got := rp.IsUnknown(); got {
		t.Error("IsUnknown() = true, want false")
	}

	expectedString := "ResourcePattern{resourceType=TOPIC, name=\"test-topic\", patternType=LITERAL}"
	if got := rp.String(); got != expectedString {
		t.Errorf("String() = %q, want %q", got, expectedString)
	}
}

// TestResourcePatternEquals tests the Equals method.
func TestResourcePatternEquals(t *testing.T) {
	rp1, _ := NewResourcePattern(TOPIC, "test-topic", LITERAL)
	rp2, _ := NewResourcePattern(TOPIC, "test-topic", LITERAL)
	rp3, _ := NewResourcePattern(GROUP, "test-group", PREFIXED)

	if !rp1.Equals(rp2) {
		t.Error("Expected rp1 to equal rp2")
	}

	if rp1.Equals(rp3) {
		t.Error("Expected rp1 to not equal rp3")
	}

	if rp1.Equals(nil) {
		t.Error("Expected rp1 to not equal nil")
	}
}
