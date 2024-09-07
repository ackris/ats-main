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
type MetricNameTemplate struct {
	Name        string
	Group       string
	Description string
	Tags        []string // Slice to maintain the order of tags
}

// NewMetricNameTemplate creates a new MetricNameTemplate with the given parameters.
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

// Name returns the name of the metric.
func (m *MetricNameTemplate) GetName() string {
	return m.Name
}

// Group returns the group of the metric.
func (m *MetricNameTemplate) GetGroup() string {
	return m.Group
}

// Description returns the description of the metric.
func (m *MetricNameTemplate) GetDescription() string {
	return m.Description
}

// Tags returns the ordered list of tag names for the metric.
func (m *MetricNameTemplate) GetTags() []string {
	return m.Tags
}

// String returns a string representation of the MetricNameTemplate.
func (m *MetricNameTemplate) String() string {
	return fmt.Sprintf("name=%s, group=%s, tags=%s", m.Name, m.Group, strings.Join(m.Tags, ", "))
}

// Equals checks if two MetricNameTemplates are equal.
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
