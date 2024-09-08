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
