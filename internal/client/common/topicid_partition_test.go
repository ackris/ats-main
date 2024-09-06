package common

import (
	"fmt"
	"testing"
)

// Test for NewTopicIdPartition
func TestNewTopicIdPartition(t *testing.T) {
	validUuid := NewUuid(0x123456789abcdef0, 0xfedcba9876543210)
	tp := NewTopicPartition("test-topic", 0)

	tip, err := NewTopicIdPartition(validUuid, tp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if tip.TopicId() != validUuid {
		t.Errorf("Expected topicId %v, got %v", validUuid, tip.TopicId())
	}
	if tip.Topic() != "test-topic" {
		t.Errorf("Expected topic %s, got %s", "test-topic", tip.Topic())
	}
	if tip.Partition() != 0 {
		t.Errorf("Expected partition %d, got %d", 0, tip.Partition())
	}
}

// Test for NewTopicIdPartition with nil TopicPartition
func TestNewTopicIdPartition_NilTopicPartition(t *testing.T) {
	validUuid := NewUuid(0x123456789abcdef0, 0xfedcba9876543210)

	tip, err := NewTopicIdPartition(validUuid, nil)
	if tip != nil {
		t.Errorf("Expected nil, got %v", tip)
	}
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

// Test for NewTopicIdPartitionWithParams
func TestNewTopicIdPartitionWithParams(t *testing.T) {
	validUuid := NewUuid(0x123456789abcdef0, 0xfedcba9876543210)
	tip, err := NewTopicIdPartitionWithParams(validUuid, 0, "test-topic")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if tip.Topic() != "test-topic" {
		t.Errorf("Expected topic %s, got %s", "test-topic", tip.Topic())
	}
	if tip.Partition() != 0 {
		t.Errorf("Expected partition %d, got %d", 0, tip.Partition())
	}
}

// Test for NewTopicIdPartitionWithParams with nil Uuid
func TestNewTopicIdPartitionWithParams_NilUuid(t *testing.T) {
	tip, err := NewTopicIdPartitionWithParams(Uuid{}, 0, "test-topic")
	if tip != nil {
		t.Errorf("Expected nil, got %v", tip)
	}
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

// Test for Equals method
func TestTopicIdPartition_Equals(t *testing.T) {
	validUuid := NewUuid(0x123456789abcdef0, 0xfedcba9876543210)
	tp1 := NewTopicPartition("test-topic", 0)
	tp2 := NewTopicPartition("test-topic", 0)

	tip1, _ := NewTopicIdPartition(validUuid, tp1)
	tip2, _ := NewTopicIdPartition(validUuid, tp2)

	if !tip1.Equals(tip2) {
		t.Errorf("Expected tips to be equal, but they are not")
	}
}

// Test for String method
func TestTopicIdPartition_String(t *testing.T) {
	validUuid := NewUuid(0x123456789abcdef0, 0xfedcba9876543210)
	tp := NewTopicPartition("test-topic", 0)

	tip, _ := NewTopicIdPartition(validUuid, tp)
	expectedString := fmt.Sprintf("%s:test-topic-%d", validUuid.String(), 0)
	if tip.String() != expectedString {
		t.Errorf("Expected string %s, got %s", expectedString, tip.String())
	}
}
