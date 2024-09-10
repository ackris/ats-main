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

func TestResourceType(t *testing.T) {
	t.Run("FromString", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    string
			expected ResourceType
			hasError bool
		}{
			{"EmptyInput", "", UNRECOGNIZED, true},
			{"UnknownType", "UNKNOWN", UNRECOGNIZED, true},
			{"ValidType", "TOPIC", TOPIC, false},
			{"CaseInsensitive", "tOpIc", TOPIC, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				actual, err := FromString(tc.input)
				if tc.hasError && err == nil {
					t.Errorf("expected error, got nil")
				}
				if !tc.hasError && err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if actual != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, actual)
				}
			})
		}
	})

	t.Run("FromCode", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    byte
			expected ResourceType
		}{
			{"ValidCode", byte(TOPIC), TOPIC},
			{"UnknownCode", byte(resourceTypeCount), UNRECOGNIZED},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				actual := FromCode(tc.input)
				if actual != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, actual)
				}
			})
		}
	})

	t.Run("Code", func(t *testing.T) {
		if TOPIC.Code() != byte(TOPIC) {
			t.Errorf("expected %d, got %d", TOPIC, TOPIC.Code())
		}
	})

	t.Run("IsUnrecognized", func(t *testing.T) {
		if UNRECOGNIZED.IsUnrecognized() != true {
			t.Errorf("expected true, got false")
		}
		if TOPIC.IsUnrecognized() != false {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("String", func(t *testing.T) {
		if TOPIC.String() != "TOPIC" {
			t.Errorf("expected TOPIC, got %s", TOPIC.String())
		}
		if ResourceType(resourceTypeCount).String() != "UNKNOWN" {
			t.Errorf("expected UNKNOWN, got %s", ResourceType(resourceTypeCount).String())
		}
	})
}

func TestResource(t *testing.T) {
	t.Run("NewResource", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("NewResource with UNRECOGNIZED type did not panic")
			}
		}()
		NewResource(UNRECOGNIZED, "test")
	})

	t.Run("ResourceType", func(t *testing.T) {
		res := NewResource(TOPIC, "test-topic")
		if res.ResourceType() != TOPIC {
			t.Errorf("expected TOPIC, got %v", res.ResourceType())
		}
	})

	t.Run("Name", func(t *testing.T) {
		res := NewResource(TOPIC, "test-topic")
		if res.Name() != "test-topic" {
			t.Errorf("expected test-topic, got %s", res.Name())
		}
	})

	t.Run("IsUnknown", func(t *testing.T) {
		res := NewResource(TOPIC, "test-topic")
		if res.IsUnknown() {
			t.Errorf("expected false, got true")
		}
		// Create a resource with an unrecognized type
		resUnknown := NewResource(ResourceType(resourceTypeCount), "test")
		if !resUnknown.IsUnknown() {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("String", func(t *testing.T) {
		res := NewResource(TOPIC, "test-topic")
		expected := "(resourceType=TOPIC, name=test-topic)"
		if res.String() != expected {
			t.Errorf("expected %s, got %s", expected, res.String())
		}
	})

	t.Run("Equals", func(t *testing.T) {
		res1 := NewResource(TOPIC, "test-topic")
		res2 := NewResource(TOPIC, "test-topic")
		res3 := NewResource(TOPIC, "other-topic")
		res4 := NewResource(CLUSTER, "test-topic")

		if !res1.Equals(res2) {
			t.Errorf("expected %v to be equal to %v", res1, res2)
		}
		if res1.Equals(res3) {
			t.Errorf("expected %v to not be equal to %v", res1, res3)
		}
		if res1.Equals(res4) {
			t.Errorf("expected %v to not be equal to %v", res1, res4)
		}
		if res1.Equals(nil) {
			t.Errorf("expected %v to not be equal to nil", res1)
		}
	})
}
