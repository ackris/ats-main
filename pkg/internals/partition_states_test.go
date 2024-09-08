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

package internals

import (
	"testing"

	"github.com/ackris/ats-main/pkg/common"
)

func TestPartitionStates(t *testing.T) {
	// Create a new PartitionStates instance
	ps := NewPartitionStates[string]()

	// Test adding and updating states
	tp1 := common.NewTopicPartition("topic1", 1)
	tp2 := common.NewTopicPartition("topic2", 2)
	ps.Update(tp1, "state1")
	ps.Update(tp2, "state2")

	// Check size
	if got := ps.Size(); got != 2 {
		t.Errorf("Size() = %v; want 2", got)
	}

	// Check values
	if state, exists := ps.StateValue(tp1); !exists || state != "state1" {
		t.Errorf("StateValue(tp1) = %v, %v; want %v, true", state, exists, "state1")
	}

	if state, exists := ps.StateValue(tp2); !exists || state != "state2" {
		t.Errorf("StateValue(tp2) = %v, %v; want %v, true", state, exists, "state2")
	}

	// Move tp1 to end
	ps.MoveToEnd(tp1)
	if got := ps.Size(); got != 2 {
		t.Errorf("Size() = %v; want 2", got)
	}

	// Test PartitionSet
	set := ps.PartitionSet()
	if len(set) != 2 {
		t.Errorf("PartitionSet() = %v; want length 2", len(set))
	}

	// Test PartitionStateMap
	stateMap := ps.PartitionStateMap()
	if len(stateMap) != 2 {
		t.Errorf("PartitionStateMap() = %v; want length 2", len(stateMap))
	}

	// Test PartitionStateValues
	values := ps.PartitionStateValues()
	if len(values) != 2 {
		t.Errorf("PartitionStateValues() = %v; want length 2", len(values))
	}

	// Test Contains
	if !ps.Contains(tp1) || !ps.Contains(tp2) {
		t.Errorf("Contains(tp1) or Contains(tp2) = false; want true")
	}

	// Test Remove
	ps.Remove(tp1)
	if ps.Contains(tp1) {
		t.Errorf("Contains(tp1) after Remove() = true; want false")
	}

	// Test Clear
	ps.Clear()
	if got := ps.Size(); got != 0 {
		t.Errorf("Size() after Clear() = %v; want 0", got)
	}

	// Test Set with new data
	ps.Set(map[*common.TopicPartition]string{
		tp1: "newstate1",
	})
	if state, exists := ps.StateValue(tp1); !exists || state != "newstate1" {
		t.Errorf("StateValue(tp1) after Set() = %v, %v; want %v, true", state, exists, "newstate1")
	}
}

func TestPartitionState(t *testing.T) {
	tp1 := common.NewTopicPartition("topic1", 1)
	tp2 := common.NewTopicPartition("topic2", 2)

	ps1 := NewPartitionState(tp1, "state1")
	ps2 := NewPartitionState(tp1, "state1")
	ps3 := NewPartitionState(tp2, "state2")

	// Test Equals
	if !ps1.Equals(ps2) {
		t.Errorf("Equals(ps1, ps2) = false; want true")
	}

	if ps1.Equals(ps3) {
		t.Errorf("Equals(ps1, ps3) = true; want false")
	}

	// Test String representation
	expectedString := "PartitionState(topic1-1=state1)"
	if got := ps1.String(); got != expectedString {
		t.Errorf("String() = %v; want %v", got, expectedString)
	}
}

func TestNewPartitionStates(t *testing.T) {
	ps := NewPartitionStates[string]()
	if ps.Size() != 0 {
		t.Errorf("Expected size 0, got %d", ps.Size())
	}
}

func TestUpdate(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp := common.NewTopicPartition("topicA", 1)
	ps.Update(tp, "state1")

	if state, exists := ps.StateValue(tp); !exists || state != "state1" {
		t.Errorf("Expected state 'state1', got %v", state)
	}
}

func TestRemove(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp := common.NewTopicPartition("topicB", 2)
	ps.Update(tp, "state2")
	ps.Remove(tp)

	if _, exists := ps.StateValue(tp); exists {
		t.Errorf("Partition should be removed")
	}
}

func TestMoveToEnd(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp1 := common.NewTopicPartition("topicC", 3)
	tp2 := common.NewTopicPartition("topicD", 4)
	ps.Update(tp1, "state3")
	ps.Update(tp2, "state4")

	ps.MoveToEnd(tp1)

	values := ps.PartitionStateValues()
	if len(values) != 2 {
		t.Errorf("Expected 2 states, got %d", len(values))
	}

	if values[0] != "state4" || values[1] != "state3" {
		t.Errorf("Expected states to be ['state4', 'state3'], got %v", values)
	}
}

func TestUpdateAndMoveToEnd(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp1 := common.NewTopicPartition("topicE", 5)
	tp2 := common.NewTopicPartition("topicF", 6)
	ps.Update(tp1, "state5")
	ps.Update(tp2, "state6")

	ps.UpdateAndMoveToEnd(tp1, "newState5")

	// Check if tp1 was moved to end and its state was updated
	values := ps.PartitionStateValues()
	if len(values) != 2 {
		t.Errorf("Expected 2 states, got %d", len(values))
	}

	if values[0] != "state6" || values[1] != "newState5" {
		t.Errorf("Expected states to be ['state6', 'newState5'], got %v", values)
	}

	// Check state of tp1
	if state, exists := ps.StateValue(tp1); !exists || state != "newState5" {
		t.Errorf("Expected state 'newState5', got %v", state)
	}
}

func TestPartitionSet(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp1 := common.NewTopicPartition("topicG", 7)
	tp2 := common.NewTopicPartition("topicH", 8)
	ps.Update(tp1, "state7")
	ps.Update(tp2, "state8")

	set := ps.PartitionSet()
	if len(set) != 2 {
		t.Errorf("Expected 2 partitions, got %d", len(set))
	}

	if _, exists := set[tp1]; !exists {
		t.Errorf("Partition tp1 not found in set")
	}
	if _, exists := set[tp2]; !exists {
		t.Errorf("Partition tp2 not found in set")
	}
}

func TestClear(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp := common.NewTopicPartition("topicI", 9)
	ps.Update(tp, "state9")
	ps.Clear()

	if ps.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", ps.Size())
	}
}

func TestContains(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp := common.NewTopicPartition("topicJ", 10)
	ps.Update(tp, "state10")

	if !ps.Contains(tp) {
		t.Errorf("Partition should exist")
	}
}

func TestStateIterator(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp := common.NewTopicPartition("topicK", 11)
	ps.Update(tp, "state11")

	iter := ps.StateIterator()
	var states []string
	for state := range iter {
		states = append(states, state)
	}

	if len(states) != 1 || states[0] != "state11" {
		t.Errorf("Expected state 'state11', got %v", states)
	}
}

func TestForEach(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp1 := common.NewTopicPartition("topicL", 12)
	tp2 := common.NewTopicPartition("topicM", 13)
	ps.Update(tp1, "state12")
	ps.Update(tp2, "state13")

	count := 0
	ps.ForEach(func(tp *common.TopicPartition, state string) {
		if state != "state12" && state != "state13" {
			t.Errorf("Unexpected state value: %v", state)
		}
		count++
	})

	if count != 2 {
		t.Errorf("Expected 2 calls to ForEach, got %d", count)
	}
}

func TestPartitionStateMap(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp := common.NewTopicPartition("topicN", 14)
	ps.Update(tp, "state14")

	stateMap := ps.PartitionStateMap()
	if len(stateMap) != 1 || stateMap[tp] != "state14" {
		t.Errorf("Expected map with one entry: topicN-14 -> state14")
	}
}

func TestPartitionStateValues(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp1 := common.NewTopicPartition("topicO", 15)
	tp2 := common.NewTopicPartition("topicP", 16)
	ps.Update(tp1, "state15")
	ps.Update(tp2, "state16")

	values := ps.PartitionStateValues()
	if len(values) != 2 || values[0] != "state15" || values[1] != "state16" {
		t.Errorf("Expected values ['state15', 'state16'], got %v", values)
	}
}

func TestStateValue(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp := common.NewTopicPartition("topicQ", 17)
	ps.Update(tp, "state17")

	state, exists := ps.StateValue(tp)
	if !exists || state != "state17" {
		t.Errorf("Expected state 'state17', got %v", state)
	}
}

func TestSize(t *testing.T) {
	ps := NewPartitionStates[string]()
	tp := common.NewTopicPartition("topicR", 18)
	ps.Update(tp, "state18")

	if ps.Size() != 1 {
		t.Errorf("Expected size 1, got %d", ps.Size())
	}
}

func TestSet(t *testing.T) {
	tp1 := common.NewTopicPartition("topicS", 19)
	tp2 := common.NewTopicPartition("topicT", 20)
	partitionToState := map[*common.TopicPartition]string{
		tp1: "state19",
		tp2: "state20",
	}

	ps := NewPartitionStates[string]()
	ps.Set(partitionToState)

	if size := ps.Size(); size != 2 {
		t.Errorf("Expected size 2, got %d", size)
	}

	if state, exists := ps.StateValue(tp1); !exists || state != "state19" {
		t.Errorf("Expected state 'state19', got %v", state)
	}
}
