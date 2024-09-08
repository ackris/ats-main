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

func TestNewMetricName(t *testing.T) {
	tags := map[string]string{"key1": "value1", "key2": "value2"}
	metricName := NewMetricName("metric1", "group1", "description1", tags)

	if metricName.Name != "metric1" {
		t.Errorf("Expected Name 'metric1', got %s", metricName.Name)
	}
	if metricName.Group != "group1" {
		t.Errorf("Expected Group 'group1', got %s", metricName.Group)
	}
	if metricName.Description != "description1" {
		t.Errorf("Expected Description 'description1', got %s", metricName.Description)
	}
	if len(metricName.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(metricName.Tags))
	}
	if metricName.Tags["key1"] != "value1" {
		t.Errorf("Expected tag 'key1' to be 'value1', got %s", metricName.Tags["key1"])
	}
}

func TestMetricNameMethods(t *testing.T) {
	tags := map[string]string{"key1": "value1"}
	metricName := NewMetricName("metric2", "group2", "description2", tags)

	if got := metricName.GetName(); got != "metric2" {
		t.Errorf("GetName() = %v; want 'metric2'", got)
	}
	if got := metricName.GetGroup(); got != "group2" {
		t.Errorf("GetGroup() = %v; want 'group2'", got)
	}
	if got := metricName.GetTags(); len(got) != 1 || got["key1"] != "value1" {
		t.Errorf("GetTags() = %v; want {'key1': 'value1'}", got)
	}
	if got := metricName.GetDescription(); got != "description2" {
		t.Errorf("GetDescription() = %v; want 'description2'", got)
	}
}

func TestMetricNameHash(t *testing.T) {
	tags := map[string]string{"key1": "value1"}
	metricName := NewMetricName("metric3", "group3", "description3", tags)
	hash := metricName.Hash()

	// Check if hash value is not zero. This is a simple check since hash values can vary.
	if hash == 0 {
		t.Errorf("Hash() = %v; want non-zero value", hash)
	}
}
