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
func (tpr *TopicPartitionReplica) Topic() string {
	return tpr.topic
}

// Partition returns the partition number.
func (tpr *TopicPartitionReplica) Partition() int {
	return tpr.partition
}

// BrokerID returns the broker ID.
func (tpr *TopicPartitionReplica) BrokerID() int {
	return tpr.brokerID
}

// Hash returns the hash code for the TopicPartitionReplica.
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
func (tpr *TopicPartitionReplica) Equals(other *TopicPartitionReplica) bool {
	return other != nil && tpr.partition == other.partition && tpr.brokerID == other.brokerID && tpr.topic == other.topic
}

// String returns a string representation of the TopicPartitionReplica.
func (tpr *TopicPartitionReplica) String() string {
	return fmt.Sprintf("%s-%d-%d", tpr.topic, tpr.partition, tpr.brokerID)
}
