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
)

// TopicPartitionReplica represents a topic name, partition number, and broker ID of the replica.
type TopicPartitionReplica struct {
	topic     string
	partition int
	brokerID  int
	hash      uint32
	isHashSet bool // Flag to check if hash is already computed
}

// NewTopicPartitionReplica creates a new TopicPartitionReplica instance.
//
// Example:
//
//	tpr := NewTopicPartitionReplica("example-topic", 0, 1)
//
// Parameters:
//
//	topic (string): The name of the topic.
//	partition (int): The partition number.
//	brokerID (int): The broker ID.
//
// Returns:
//
//	*TopicPartitionReplica: A pointer to the created TopicPartitionReplica instance.
//
// Panics:
//
//	If the topic is an empty string.
func NewTopicPartitionReplica(topic string, partition, brokerID int) *TopicPartitionReplica {
	if topic == "" {
		panic("topic cannot be empty")
	}
	return &TopicPartitionReplica{
		topic:     topic,
		partition: partition,
		brokerID:  brokerID,
	}
}

// Topic returns the topic name.
//
// Example:
//
//	fmt.Println(tpr.Topic()) // Output: example-topic
//
// Returns:
//
//	string: The topic name.
func (tpr *TopicPartitionReplica) Topic() string {
	return tpr.topic
}

// Partition returns the partition number.
//
// Example:
//
//	fmt.Println(tpr.Partition()) // Output: 0
//
// Returns:
//
//	int: The partition number.
func (tpr *TopicPartitionReplica) Partition() int {
	return tpr.partition
}

// BrokerID returns the broker ID.
//
// Example:
//
//	fmt.Println(tpr.BrokerID()) // Output: 1
//
// Returns:
//
//	int: The broker ID.
func (tpr *TopicPartitionReplica) BrokerID() int {
	return tpr.brokerID
}

// Hash returns the hash code for the TopicPartitionReplica.
//
// Example:
//
//	hash := tpr.Hash()
//	fmt.Println(hash) // Output: 1234567890
//
// Returns:
//
//	uint32: The hash code for the TopicPartitionReplica.
func (tpr *TopicPartitionReplica) Hash() uint32 {
	if !tpr.isHashSet {
		h := fnv.New32a()
		h.Write([]byte(tpr.topic))
		h.Write([]byte(fmt.Sprintf("%d", tpr.partition)))
		h.Write([]byte(fmt.Sprintf("%d", tpr.brokerID)))
		tpr.hash = h.Sum32()
		tpr.isHashSet = true
	}
	return tpr.hash
}

// Equals checks if two TopicPartitionReplica instances are equal.
// Example:
//
//	tpr1 := NewTopicPartitionReplica("example-topic", 0, 1)
//	tpr2 := NewTopicPartitionReplica("example-topic", 0, 1)
//	if tpr1.Equals(tpr2) {
//	  fmt.Println("TopicPartitionReplica instances are equal")
//	}
//
// Parameters:
//
//	other (*TopicPartitionReplica): The other TopicPartitionReplica instance to compare with.
//
// Returns:
//
//	bool: True if the TopicPartitionReplica instances are equal, false otherwise.
func (tpr *TopicPartitionReplica) Equals(other *TopicPartitionReplica) bool {
	return other != nil && tpr.partition == other.partition && tpr.brokerID == other.brokerID && tpr.topic == other.topic
}

// String returns a string representation of the TopicPartitionReplica.
//
// Example:
//
//	fmt.Println(tpr.String()) // Output: example-topic-0-1
//
// Returns:
//
//	string: The string representation of the TopicPartitionReplica.
func (tpr *TopicPartitionReplica) String() string {
	return fmt.Sprintf("%s-%d-%d", tpr.topic, tpr.partition, tpr.brokerID)
}
