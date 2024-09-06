package common

import "testing"

// TestNewTopicPartition verifies that the NewTopicPartition function correctly initializes the TopicPartition.
func TestNewTopicPartition(t *testing.T) {
	tp := NewTopicPartition("topicA", 1)
	if tp.Topic() != "topicA" {
		t.Errorf("Expected topic 'topicA', got '%s'", tp.Topic())
	}
	if tp.Partition() != 1 {
		t.Errorf("Expected partition 1, got %d", tp.Partition())
	}
}

// TestPartition verifies that the Partition method returns the correct partition number.
func TestPartition(t *testing.T) {
	tp := NewTopicPartition("topicB", 2)
	if tp.Partition() != 2 {
		t.Errorf("Expected partition 2, got %d", tp.Partition())
	}
}

// TestTopic verifies that the Topic method returns the correct topic name.
func TestTopic(t *testing.T) {
	tp := NewTopicPartition("topicC", 3)
	if tp.Topic() != "topicC" {
		t.Errorf("Expected topic 'topicC', got '%s'", tp.Topic())
	}
}

// TestHashCode verifies that the HashCode method computes the correct hash code.
func TestHashCode(t *testing.T) {
	tp1 := NewTopicPartition("topicD", 4)
	tp2 := NewTopicPartition("topicD", 4)

	// Ensure that hash codes for identical TopicPartitions are the same
	if tp1.HashCode() != tp2.HashCode() {
		t.Errorf("Expected same hash codes, got %d and %d", tp1.HashCode(), tp2.HashCode())
	}

	// Check hash code for different TopicPartitions
	tp3 := NewTopicPartition("topicE", 5)
	if tp1.HashCode() == tp3.HashCode() {
		t.Errorf("Expected different hash codes, got %d and %d", tp1.HashCode(), tp3.HashCode())
	}
}

// TestEquals verifies that the Equals method correctly identifies equal and non-equal TopicPartitions.
func TestEquals(t *testing.T) {
	tp1 := NewTopicPartition("topicF", 6)
	tp2 := NewTopicPartition("topicF", 6)
	tp3 := NewTopicPartition("topicF", 7)
	tp4 := NewTopicPartition("topicG", 6)

	if !tp1.Equals(tp2) {
		t.Errorf("Expected tp1 to be equal to tp2")
	}
	if tp1.Equals(tp3) {
		t.Errorf("Expected tp1 to not be equal to tp3")
	}
	if tp1.Equals(tp4) {
		t.Errorf("Expected tp1 to not be equal to tp4")
	}
	if tp1.Equals(nil) {
		t.Errorf("Expected tp1 to not be equal to nil")
	}
}

// TestString verifies that the String method returns the correct string representation.
func TestString(t *testing.T) {
	tp := NewTopicPartition("topicH", 8)
	expected := "topicH-8"
	if tp.String() != expected {
		t.Errorf("Expected string '%s', got '%s'", expected, tp.String())
	}
}
