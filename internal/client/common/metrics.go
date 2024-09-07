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
	"hash/fnv"
)

// MetricName represents the name, group, description, and tags of a metric.
type MetricName struct {
	Name        string
	Group       string
	Description string
	Tags        map[string]string
}

// NewMetricName creates and returns a new MetricName instance.
func NewMetricName(name, group, description string, tags map[string]string) *MetricName {
	return &MetricName{
		Name:        name,
		Group:       group,
		Description: description,
		Tags:        tags,
	}
}

// GetName returns the name of the metric.
func (m *MetricName) GetName() string {
	return m.Name
}

// GetGroup returns the group of the metric.
func (m *MetricName) GetGroup() string {
	return m.Group
}

// GetTags returns the tags associated with the metric.
func (m *MetricName) GetTags() map[string]string {
	return m.Tags
}

// GetDescription returns the description of the metric.
func (m *MetricName) GetDescription() string {
	return m.Description
}

// Hash returns a hash code for the MetricName.
func (m *MetricName) Hash() uint32 {
	h := fnv.New32a()
	h.Write([]byte(m.Name))
	h.Write([]byte(m.Group))
	for key, value := range m.Tags {
		h.Write([]byte(key))
		h.Write([]byte(value))
	}
	return h.Sum32()
}

// Metric is an interface that defines methods for interacting with metrics.
type Metric interface {
	GetMetricName() *MetricName
	GetMetricValue() interface{}
}
