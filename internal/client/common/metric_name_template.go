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
	"fmt"
	"hash/fnv"
	"strings"
)

// MetricNameTemplate represents a template for a metric name with tags.
// It contains fields for the metric name, group, description, and tags.
// Tags are preserved in the order they are provided, with duplicates removed.
type MetricNameTemplate struct {
	Name        string
	Group       string
	Description string
	Tags        []string // Slice to maintain the order of tags
}

// NewMetricNameTemplate creates a new MetricNameTemplate with the given parameters.
// It ensures that name, group, and description are not empty and removes duplicate tags while preserving their order.
//
// Parameters:
//   - name: the name of the metric; must not be an empty string.
//   - group: the name of the group; must not be an empty string.
//   - description: the description of the metric; must not be an empty string.
//   - tags: a slice of metric tag names, which can be empty or nil. Duplicates will be removed, and the order will be preserved.
//
// Returns:
//   - A pointer to a MetricNameTemplate instance.
//   - An error if any of the name, group, or description parameters are empty.
//
// Example:
//
//	template, err := NewMetricNameTemplate("metric1", "group1", "A test metric", []string{"tag1", "tag2", "tag1"})
//	if err != nil {
//	    fmt.Println("Error creating MetricNameTemplate:", err)
//	} else {
//	    fmt.Println(template)
//	}
func NewMetricNameTemplate(name, group, description string, tags []string) (*MetricNameTemplate, error) {
	if name == "" || group == "" || description == "" {
		return nil, fmt.Errorf("name, group, and description must not be empty")
	}
	if tags == nil {
		tags = []string{}
	}
	// Create a unique ordered tags
	uniqueTags := createOrderedTags(tags)
	return &MetricNameTemplate{
		Name:        name,
		Group:       group,
		Description: description,
		Tags:        uniqueTags,
	}, nil
}

// createOrderedTags ensures uniqueness while preserving order.
// It takes a slice of tags, removes any duplicates, and maintains the order of the first occurrence.
//
// Parameters:
//   - tags: a slice of tag names.
//
// Returns:
//   - A slice of unique tag names in the original order.
//
// Example:
//
//	tags := []string{"tag1", "tag2", "tag1"}
//	uniqueTags := createOrderedTags(tags)
//	fmt.Println(uniqueTags) // Output: [tag1 tag2]
func createOrderedTags(tags []string) []string {
	tagMap := make(map[string]struct{})
	var orderedTags []string

	for _, tag := range tags {
		if _, exists := tagMap[tag]; !exists {
			tagMap[tag] = struct{}{}
			orderedTags = append(orderedTags, tag)
		}
	}

	return orderedTags
}

// GetName returns the name of the metric.
//
// Returns:
//   - The name of the metric.
//
// Example:
//
//	template, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", nil)
//	fmt.Println(template.GetName()) // Output: metric1
func (m *MetricNameTemplate) GetName() string {
	return m.Name
}

// GetGroup returns the group of the metric.
//
// Returns:
//   - The group of the metric.
//
// Example:
//
//	template, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", nil)
//	fmt.Println(template.GetGroup()) // Output: group1
func (m *MetricNameTemplate) GetGroup() string {
	return m.Group
}

// GetDescription returns the description of the metric.
//
// Returns:
//   - The description of the metric.
//
// Example:
//
//	template, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", nil)
//	fmt.Println(template.GetDescription()) // Output: A test metric
func (m *MetricNameTemplate) GetDescription() string {
	return m.Description
}

// GetTags returns the ordered list of tag names for the metric.
//
// Returns:
//   - A slice of tag names in the order they were provided, with duplicates removed.
//
// Example:
//
//	template, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", []string{"tag1", "tag2", "tag1"})
//	fmt.Println(template.GetTags()) // Output: [tag1 tag2]
func (m *MetricNameTemplate) GetTags() []string {
	return m.Tags
}

// String returns a string representation of the MetricNameTemplate.
// The string format is "name=<name>, group=<group>, tags=<tags>".
//
// Returns:
//   - A string representation of the MetricNameTemplate.
//
// Example:
//
//	template, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", []string{"tag1", "tag2"})
//	fmt.Println(template.String()) // Output: name=metric1, group=group1, tags=tag1, tag2
func (m *MetricNameTemplate) String() string {
	return fmt.Sprintf("name=%s, group=%s, tags=%s", m.Name, m.Group, strings.Join(m.Tags, ", "))
}

// Equals checks if two MetricNameTemplates are equal.
// Two MetricNameTemplates are considered equal if they have the same name, group, description, and tags in the same order.
//
// Parameters:
//   - other: another MetricNameTemplate to compare against.
//
// Returns:
//   - True if the MetricNameTemplates are equal, false otherwise.
//
// Example:
//
//	t1, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", []string{"tag1", "tag2"})
//	t2, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", []string{"tag1", "tag2"})
//	fmt.Println(t1.Equals(t2)) // Output: true
func (m *MetricNameTemplate) Equals(other *MetricNameTemplate) bool {
	if m == other {
		return true
	}
	if other == nil {
		return false
	}
	if m.Name != other.Name || m.Group != other.Group || len(m.Tags) != len(other.Tags) {
		return false
	}
	for i, tag := range m.Tags {
		if tag != other.Tags[i] {
			return false
		}
	}
	return true
}

// HashCode returns a hash code for the MetricNameTemplate.
// The hash code is computed based on the name, group, description, and tags.
//
// Returns:
//   - A hash code as a uint32.
//
// Example:
//
//	template, _ := NewMetricNameTemplate("metric1", "group1", "A test metric", []string{"tag1", "tag2"})
//	fmt.Println(template.HashCode()) // Output: some hash code
func (m *MetricNameTemplate) HashCode() uint32 {
	h := fnv.New32a()
	h.Write([]byte(m.Name))
	h.Write([]byte(m.Group))
	h.Write([]byte(m.Description))
	for _, tag := range m.Tags {
		h.Write([]byte(tag))
	}
	return h.Sum32()
}
