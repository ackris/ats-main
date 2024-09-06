package common

import (
	"testing"
)

func TestNewTopicPartitionReplica(t *testing.T) {
	// Test valid creation
	tpr := NewTopicPartitionReplica("test-topic", 0, 1)
	if tpr.Topic() != "test-topic" || tpr.Partition() != 0 || tpr.BrokerID() != 1 {
		t.Errorf("Expected TopicPartitionReplica with topic 'test-topic', partition 0, brokerID 1, got: %v", tpr)
	}

	// Test empty topic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when creating TopicPartitionReplica with empty topic")
		}
	}()
	NewTopicPartitionReplica("", 0, 1)
}

func TesPartitionReplicaHash(t *testing.T) {
	tpr := NewTopicPartitionReplica("test-topic", 0, 1)
	hash1 := tpr.Hash()
	hash2 := tpr.Hash() // Should be cached

	if hash1 != hash2 {
		t.Errorf("Hash should be the same on subsequent calls, got %d and %d", hash1, hash2)
	}
}

func TestPartitionReplicaEquals(t *testing.T) {
	tpr1 := NewTopicPartitionReplica("test-topic", 0, 1)
	tpr2 := NewTopicPartitionReplica("test-topic", 0, 1)
	tpr3 := NewTopicPartitionReplica("test-topic", 1, 1)

	if !tpr1.Equals(tpr2) {
		t.Errorf("Expected tpr1 to equal tpr2")
	}

	if tpr1.Equals(tpr3) {
		t.Errorf("Expected tpr1 not to equal tpr3")
	}

	if tpr1.Equals(nil) {
		t.Errorf("Expected tpr1 not to equal nil")
	}
}

func TestPartitionReplicaString(t *testing.T) {
	tpr := NewTopicPartitionReplica("test-topic", 0, 1)
	expected := "test-topic-0-1"
	if tpr.String() != expected {
		t.Errorf("Expected string representation '%s', got '%s'", expected, tpr.String())
	}
}

func TestMultipleInstances(t *testing.T) {
	tpr1 := NewTopicPartitionReplica("topic1", 0, 1)
	tpr2 := NewTopicPartitionReplica("topic1", 0, 1)
	tpr3 := NewTopicPartitionReplica("topic2", 1, 2)

	if !tpr1.Equals(tpr2) {
		t.Errorf("Expected tpr1 to equal tpr2")
	}

	if tpr1.Equals(tpr3) {
		t.Errorf("Expected tpr1 not to equal tpr3")
	}
}
