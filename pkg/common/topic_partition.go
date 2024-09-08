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
	"sync"
)

// TopicPartition represents a topic name and partition number in a messaging system.
// It is designed to be immutable and provides methods for comparison and hashing.
type TopicPartition struct {
	partition int       // The partition number within the topic
	topic     string    // The name of the topic
	hashCode  uint64    // Cached hash code for efficient hash computations
	once      sync.Once // Ensures the hash code is computed only once
}

// NewTopicPartition creates a new instance of TopicPartition with the given topic name and partition number.
//
// Parameters:
//
//	topic (string): The name of the topic.
//	partition (int): The partition number within the topic.
//
// Returns:
//
//	*TopicPartition: A pointer to the newly created TopicPartition instance.
//
// Example:
//
//	tp := NewTopicPartition("topicA", 1)
//	fmt.Println(tp.Topic())     // Output: topicA
//	fmt.Println(tp.Partition()) // Output: 1
func NewTopicPartition(topic string, partition int) *TopicPartition {
	return &TopicPartition{
		partition: partition,
		topic:     topic,
	}
}

// Partition returns the partition number of the TopicPartition.
//
// Returns:
//
//	int: The partition number.
//
// Example:
//
//	tp := NewTopicPartition("topicA", 1)
//	fmt.Println(tp.Partition()) // Output: 1
func (tp *TopicPartition) Partition() int {
	return tp.partition
}

// Topic returns the topic name of the TopicPartition.
//
// Returns:
//
//	string: The name of the topic.
//
// Example:
//
//	tp := NewTopicPartition("topicA", 1)
//	fmt.Println(tp.Topic()) // Output: topicA
func (tp *TopicPartition) Topic() string {
	return tp.topic
}

// HashCode calculates and returns the hash code of the TopicPartition.
// The hash code is computed once and cached for efficient future use.
//
// Returns:
//
//	uint64: The hash code of the TopicPartition.
//
// Example:
//
//	tp := NewTopicPartition("topicA", 1)
//	fmt.Println(tp.HashCode()) // Output: A hash code as uint64
func (tp *TopicPartition) HashCode() uint64 {
	tp.once.Do(func() {
		h := fnv.New64a()
		h.Write([]byte(tp.topic))
		h.Write([]byte(fmt.Sprintf("%d", tp.partition)))
		tp.hashCode = h.Sum64()
	})
	return tp.hashCode
}

// Equals compares the current TopicPartition with another TopicPartition instance.
// Two TopicPartition instances are considered equal if they have the same topic name and partition number.
//
// Parameters:
//
//	other (*TopicPartition): The other TopicPartition instance to compare with.
//
// Returns:
//
//	bool: true if both TopicPartition instances are equal, false otherwise.
//
// Example:
//
//	tp1 := NewTopicPartition("topicA", 1)
//	tp2 := NewTopicPartition("topicA", 1)
//	tp3 := NewTopicPartition("topicB", 2)
//	fmt.Println(tp1.Equals(tp2)) // Output: true
//	fmt.Println(tp1.Equals(tp3)) // Output: false
func (tp *TopicPartition) Equals(other *TopicPartition) bool {
	if other == nil {
		return false
	}
	return tp.partition == other.partition && tp.topic == other.topic
}

// String returns a string representation of the TopicPartition in the format "topic-partition".
// This is useful for debugging and logging.
//
// Returns:
//
//	string: The string representation of the TopicPartition.
//
// Example:
//
//	tp := NewTopicPartition("topicA", 1)
//	fmt.Println(tp.String()) // Output: topicA-1
func (tp *TopicPartition) String() string {
	return fmt.Sprintf("%s-%d", tp.topic, tp.partition)
}
