package common

import (
	"errors"
	"fmt"
)

// TopicIdPartition represents a universally unique identifier with topic id for a topic partition.
type TopicIdPartition struct {
	topicId        Uuid
	topicPartition *TopicPartition // Changed to pointer
}

// NewTopicIdPartition creates a new TopicIdPartition instance.
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

// NewTopicIdPartitionWithParams creates a new TopicIdPartition instance with topic name and partition id.
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

// TopicId returns the universally unique id representing this topic partition.
func (tip *TopicIdPartition) TopicId() Uuid {
	return tip.topicId
}

// Topic returns the topic name or an empty string if it is unknown.
func (tip *TopicIdPartition) Topic() string {
	return tip.topicPartition.Topic()
}

// Partition returns the partition id.
func (tip *TopicIdPartition) Partition() int {
	return tip.topicPartition.Partition()
}

// TopicPartition returns the TopicPartition representing this instance.
func (tip *TopicIdPartition) TopicPartition() *TopicPartition {
	return tip.topicPartition
}

// Equals checks if two TopicIdPartition instances are equal.
func (tip *TopicIdPartition) Equals(other *TopicIdPartition) bool {
	if other == nil {
		return false
	}
	return tip.topicId == other.topicId && tip.topicPartition == other.topicPartition
}

// String returns a string representation of the TopicIdPartition.
func (tip *TopicIdPartition) String() string {
	return fmt.Sprintf("%v:%s-%d", tip.topicId, tip.Topic(), tip.Partition())
}
