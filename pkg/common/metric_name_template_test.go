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

// TestNewMetricNameTemplate_ValidInput tests the creation of MetricNameTemplate with valid inputs.
func TestNewMetricNameTemplate_ValidInput(t *testing.T) {
	name := "metric1"
	group := "group1"
	description := "A test metric"
	tags := []string{"tag1", "tag2", "tag3"}

	template, err := NewMetricNameTemplate(name, group, description, tags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if template.Name != name {
		t.Errorf("expected name %s, got %s", name, template.Name)
	}
	if template.Group != group {
		t.Errorf("expected group %s, got %s", group, template.Group)
	}
	if template.Description != description {
		t.Errorf("expected description %s, got %s", description, template.Description)
	}
	expectedTags := []string{"tag1", "tag2", "tag3"}
	if !equalStringSlices(template.Tags, expectedTags) {
		t.Errorf("expected tags %v, got %v", expectedTags, template.Tags)
	}
}

// TestNewMetricNameTemplate_DuplicateTags tests that duplicate tags are removed while preserving order.
func TestNewMetricNameTemplate_DuplicateTags(t *testing.T) {
	name := "metric1"
	group := "group1"
	description := "A test metric"
	tags := []string{"tag1", "tag2", "tag1", "tag3"}

	template, err := NewMetricNameTemplate(name, group, description, tags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedTags := []string{"tag1", "tag2", "tag3"}
	if !equalStringSlices(template.Tags, expectedTags) {
		t.Errorf("expected tags %v, got %v", expectedTags, template.Tags)
	}
}

// TestNewMetricNameTemplate_EmptyTags tests the behavior when an empty tag list is provided.
func TestNewMetricNameTemplate_EmptyTags(t *testing.T) {
	name := "metric1"
	group := "group1"
	description := "A test metric"
	tags := []string{}

	template, err := NewMetricNameTemplate(name, group, description, tags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(template.Tags) != 0 {
		t.Errorf("expected no tags, got %v", template.Tags)
	}
}

// TestEquals tests the Equals method for MetricNameTemplate.
func TestEqualsMetricTemplate(t *testing.T) {
	t1, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", []string{"tag1", "tag2"})
	t2, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", []string{"tag1", "tag2"})
	t3, _ := NewMetricNameTemplate("metric2", "group1", "A different metric", []string{"tag1", "tag2"})
	t4, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", []string{"tag1"})

	if !t1.Equals(t2) {
		t.Errorf("expected t1 to be equal to t2")
	}
	if t1.Equals(t3) {
		t.Errorf("expected t1 to be not equal to t3")
	}
	if t1.Equals(t4) {
		t.Errorf("expected t1 to be not equal to t4")
	}
}

// TestHashCode tests the HashCode method for MetricNameTemplate.
func TestHashCodeMetricTemplate(t *testing.T) {
	t1, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", []string{"tag1", "tag2"})
	t2, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", []string{"tag1", "tag2"})

	if t1.HashCode() != t2.HashCode() {
		t.Errorf("expected hash codes to be equal for t1 and t2")
	}

	t3, _ := NewMetricNameTemplate("metric2", "group1", "A different metric", []string{"tag1", "tag2"})
	if t1.HashCode() == t3.HashCode() {
		t.Errorf("expected hash codes to be different for t1 and t3")
	}
}

// TestNewMetricNameTemplate_InvalidInput tests error handling for invalid input.
func TestNewMetricNameTemplate_InvalidInput(t *testing.T) {
	_, err := NewMetricNameTemplate("", "group1", "A test metric", []string{"tag1"})
	if err == nil {
		t.Errorf("expected error for empty name")
	}
	_, err = NewMetricNameTemplate("metric1", "", "A test metric", []string{"tag1"})
	if err == nil {
		t.Errorf("expected error for empty group")
	}
	_, err = NewMetricNameTemplate("metric1", "group1", "", []string{"tag1"})
	if err == nil {
		t.Errorf("expected error for empty description")
	}
}

// Helper function to compare slices of strings
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
