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
// It encapsulates the essential information required to identify and describe a metric.
//
// Example:
//
//	tags := map[string]string{"key1": "value1", "key2": "value2"}
//	metricName := NewMetricName("metric1", "group1", "description1", tags)
//	fmt.Println(metricName.GetName())        // Output: metric1
//	fmt.Println(metricName.GetGroup())       // Output: group1
//	fmt.Println(metricName.GetDescription()) // Output: description1
//	fmt.Println(metricName.GetTags())        // Output: map[key1:value1 key2:value2]
//	fmt.Println(metricName.Hash())           // Output: <some hash value>
type MetricName struct {
	Name        string
	Group       string
	Description string
	Tags        map[string]string
}

// NewMetricName creates and returns a new MetricName instance with the given parameters.
//
// Parameters:
//   - name: The name of the metric.
//   - group: The group to which this metric belongs.
//   - description: A human-readable description of the metric (optional).
//   - tags: Additional key-value pairs associated with the metric (optional).
//
// Returns:
//   - A pointer to a new MetricName instance.
//
// Example:
//
//	tags := map[string]string{"key1": "value1"}
//	metricName := NewMetricName("metric2", "group2", "description2", tags)
//	fmt.Println(metricName.GetName()) // Output: metric2
func NewMetricName(name, group, description string, tags map[string]string) *MetricName {
	return &MetricName{
		Name:        name,
		Group:       group,
		Description: description,
		Tags:        tags,
	}
}

// GetName returns the name of the metric.
//
// Returns:
//   - The name of the metric.
//
// Example:
//
//	metricName := NewMetricName("metric3", "group3", "", nil)
//	fmt.Println(metricName.GetName()) // Output: metric3
func (m *MetricName) GetName() string {
	return m.Name
}

// GetGroup returns the group of the metric.
//
// Returns:
//   - The group of the metric.
//
// Example:
//
//	metricName := NewMetricName("", "group4", "", nil)
//	fmt.Println(metricName.GetGroup()) // Output: group4
func (m *MetricName) GetGroup() string {
	return m.Group
}

// GetTags returns the tags associated with the metric.
//
// Returns:
//   - A map of tags associated with the metric.
//
// Example:
//
//	tags := map[string]string{"key1": "value1"}
//	metricName := NewMetricName("", "", "", tags)
//	fmt.Println(metricName.GetTags()) // Output: map[key1:value1]
func (m *MetricName) GetTags() map[string]string {
	return m.Tags
}

// GetDescription returns the description of the metric.
//
// Returns:
//   - The description of the metric.
//
// Example:
//
//	metricName := NewMetricName("", "", "This is a description", nil)
//	fmt.Println(metricName.GetDescription()) // Output: This is a description
func (m *MetricName) GetDescription() string {
	return m.Description
}

// Hash returns a hash code for the MetricName based on its fields.
//
// Returns:
//   - A uint32 hash code.
//
// Example:
//
//	metricName := NewMetricName("metric5", "group5", "description5", nil)
//	fmt.Println(metricName.Hash()) // Output: <some hash value>
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
// It includes methods to get the MetricName and MetricValue.
//
// Example:
//
//	type MyMetric struct {
//	    name  *MetricName
//	    value interface{}
//	}
//
//	func (m *MyMetric) GetMetricName() *MetricName {
//	    return m.name
//	}
//
//	func (m *MyMetric) GetMetricValue() interface{} {
//	    return m.value
//	}
//
//	metricName := NewMetricName("metric6", "group6", "description6", nil)
//	metricValue := 100
//	myMetric := &MyMetric{name: metricName, value: metricValue}
//
//	fmt.Println(myMetric.GetMetricName().GetName()) // Output: metric6
//	fmt.Println(myMetric.GetMetricValue())          // Output: 100
type Metric interface {
	// GetMetricName returns the MetricName of the metric.
	//
	// Returns:
	//   - The MetricName of the metric.
	GetMetricName() *MetricName

	// GetMetricValue returns the value of the metric.
	//
	// Returns:
	//   - The value of the metric.
	GetMetricValue() interface{}
}
