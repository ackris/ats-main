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
	"errors"
	"fmt"
)

// TopicIdPartition represents a universally unique identifier combined with a topic and its partition.
//
// It contains:
// - `topicId`: A UUID representing the unique identifier for the topic.
// - `topicPartition`: A pointer to a `TopicPartition` struct which holds the topic name and partition id.
type TopicIdPartition struct {
	topicId        Uuid
	topicPartition *TopicPartition // Changed to pointer
}

// NewTopicIdPartition creates a new TopicIdPartition instance.
//
// This function returns a new `TopicIdPartition` with the provided `topicId` and `topicPartition`.
//
// Parameters:
// - `topicId`: A UUID representing the unique identifier for the topic. Must not be an empty UUID.
// - `topicPartition`: A pointer to a `TopicPartition` struct. Must not be nil.
//
// Returns:
// - A pointer to the newly created `TopicIdPartition`.
// - An error if `topicId` is an empty UUID or `topicPartition` is nil.
//
// Errors:
// - "topicId cannot be nil" if `topicId` is an empty UUID.
// - "topicPartition cannot be nil" if `topicPartition` is nil.
//
// Example:
//
//	topicId := Uuid("123e4567-e89b-12d3-a456-426614174000")
//	topicPartition := &TopicPartition{topic: "exampleTopic", partition: 1}
//	tip, err := NewTopicIdPartition(topicId, topicPartition)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(tip.String()) // Output: 123e4567-e89b-12d3-a456-426614174000:exampleTopic-1
func NewTopicIdPartition(topicId Uuid, topicPartition *TopicPartition) (*TopicIdPartition, error) {
	if (topicId == Uuid{}) {
		return nil, errors.New("topicId cannot be nil")
	}
	if topicPartition == nil {
		return nil, errors.New("topicPartition cannot be nil")
	}
	return &TopicIdPartition{
		topicId:        topicId,
		topicPartition: topicPartition,
	}, nil
}

// NewTopicIdPartitionWithParams creates a new TopicIdPartition instance using topic name and partition id.
//
// This function constructs a `TopicIdPartition` from the provided `topicId`, topic name, and partition id.
//
// Parameters:
// - `topicId`: A UUID representing the unique identifier for the topic. Must not be an empty UUID.
// - `partition`: An integer representing the partition id.
// - `topic`: A string representing the topic name.
//
// Returns:
// - A pointer to the newly created `TopicIdPartition`.
// - An error if `topicId` is an empty UUID.
//
// Errors:
// - "topicId cannot be nil" if `topicId` is an empty UUID.
//
// Example:
//
//	topicId := Uuid("123e4567-e89b-12d3-a456-426614174000")
//	tip, err := NewTopicIdPartitionWithParams(topicId, 1, "exampleTopic")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(tip.String()) // Output: 123e4567-e89b-12d3-a456-426614174000:exampleTopic-1
func NewTopicIdPartitionWithParams(topicId Uuid, partition int, topic string) (*TopicIdPartition, error) {
	if (topicId == Uuid{}) {
		return nil, errors.New("topicId cannot be nil")
	}
	topicPartition := NewTopicPartition(topic, partition)
	return &TopicIdPartition{
		topicId:        topicId,
		topicPartition: topicPartition,
	}, nil
}

// TopicId returns the UUID representing this TopicIdPartition.
//
// Returns:
// - The `topicId` field.
//
// Example:
//
//	tip := &TopicIdPartition{topicId: Uuid("123e4567-e89b-12d3-a456-426614174000")}
//	fmt.Println(tip.TopicId()) // Output: 123e4567-e89b-12d3-a456-426614174000
func (tip *TopicIdPartition) TopicId() Uuid {
	return tip.topicId
}

// Topic returns the topic name associated with this TopicIdPartition.
//
// Returns:
// - The topic name as a string. Returns an empty string if the topic name is unknown.
//
// Example:
//
//	topicPartition := &TopicPartition{topic: "exampleTopic", partition: 1}
//	tip := &TopicIdPartition{topicPartition: topicPartition}
//	fmt.Println(tip.Topic()) // Output: exampleTopic
func (tip *TopicIdPartition) Topic() string {
	return tip.topicPartition.Topic()
}

// Partition returns the partition id associated with this TopicIdPartition.
//
// Returns:
// - The partition id as an integer.
//
// Example:
//
//	topicPartition := &TopicPartition{topic: "exampleTopic", partition: 1}
//	tip := &TopicIdPartition{topicPartition: topicPartition}
//	fmt.Println(tip.Partition()) // Output: 1
func (tip *TopicIdPartition) Partition() int {
	return tip.topicPartition.Partition()
}

// TopicPartition returns the TopicPartition struct representing this instance.
//
// Returns:
// - A pointer to the `TopicPartition` struct.
//
// Example:
//
//	topicPartition := &TopicPartition{topic: "exampleTopic", partition: 1}
//	tip := &TopicIdPartition{topicPartition: topicPartition}
//	fmt.Println(tip.TopicPartition()) // Output: &{exampleTopic 1}
func (tip *TopicIdPartition) TopicPartition() *TopicPartition {
	return tip.topicPartition
}

// Equals checks if two TopicIdPartition instances are equal.
//
// Two `TopicIdPartition` instances are considered equal if both their `topicId` and `topicPartition` fields are equal.
//
// Parameters:
// - `other`: A pointer to another `TopicIdPartition` instance to compare with.
//
// Returns:
// - `true` if both `topicId` and `topicPartition` fields are equal in both instances.
// - `false` otherwise.
//
// Example:
//
//	tip1 := &TopicIdPartition{topicId: Uuid("123e4567-e89b-12d3-a456-426614174000"), topicPartition: NewTopicPartition("exampleTopic", 1)}
//	tip2 := &TopicIdPartition{topicId: Uuid("123e4567-e89b-12d3-a456-426614174000"), topicPartition: NewTopicPartition("exampleTopic", 1)}
//	fmt.Println(tip1.Equals(tip2)) // Output: true
func (tip *TopicIdPartition) Equals(other *TopicIdPartition) bool {
	if other == nil {
		return false
	}
	// Compare Uuid
	if tip.topicId != other.topicId {
		return false
	}
	// Compare TopicPartition
	if tip.topicPartition.Topic() != other.topicPartition.Topic() || tip.topicPartition.Partition() != other.topicPartition.Partition() {
		return false
	}
	return true
}

// String returns a string representation of the TopicIdPartition.
//
// The format of the string is "<topicId>:<topicName>-<partitionId>", where:
// - `<topicId>` is the UUID of the topic.
// - `<topicName>` is the topic name.
// - `<partitionId>` is the partition id.
//
// Returns:
// - A string representation of the `TopicIdPartition`.
//
// Example:
//
//	topicPartition := &TopicPartition{topic: "exampleTopic", partition: 1}
//	tip := &TopicIdPartition{topicId: Uuid("123e4567-e89b-12d3-a456-426614174000"), topicPartition: topicPartition}
//	fmt.Println(tip.String()) // Output: 123e4567-e89b-12d3-a456-426614174000:exampleTopic-1
func (tip *TopicIdPartition) String() string {
	return fmt.Sprintf("%v:%s-%d", tip.topicId, tip.Topic(), tip.Partition())
}
