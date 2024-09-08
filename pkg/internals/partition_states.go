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
	"fmt"
	"reflect"
	"sync"

	"github.com/ackris/ats-main/pkg/common"
)

// PartitionStates represents a collection of states for various partitions.
// It supports concurrent access and maintains the size of the partition map.
type PartitionStates[S any] struct {
	mu       sync.RWMutex
	mapState map[*common.TopicPartition]S
	size     int
}

// NewPartitionStates creates a new PartitionStates instance with an empty state map.
//
// Returns:
//
//	*PartitionStates[S]: A pointer to the newly created PartitionStates instance.
//
// Example:
//
//	ps := NewPartitionStates[string]()
//	fmt.Println(ps.Size()) // Output: 0
func NewPartitionStates[S any]() *PartitionStates[S] {
	return &PartitionStates[S]{
		mapState: make(map[*common.TopicPartition]S),
	}
}

// MoveToEnd moves the given partition to the end of the order. This operation
// updates the internal ordering of partitions to ensure the specified partition
// is moved to the end.
//
// Parameters:
//
//	tp (*common.TopicPartition): The partition to be moved to the end.
//
// Example:
//
//	tp := common.NewTopicPartition("topicA", 1)
//	ps := NewPartitionStates[string]()
//	ps.Update(tp, "state1")
//	ps.MoveToEnd(tp)
//	fmt.Println(ps.PartitionStateValues()) // Output: ["state1"]
func (ps *PartitionStates[S]) MoveToEnd(tp *common.TopicPartition) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if state, ok := ps.mapState[tp]; ok {
		delete(ps.mapState, tp)
		ps.mapState[tp] = state
	}
}

// UpdateAndMoveToEnd updates the state of a partition and moves it to the end.
//
// Parameters:
//
//	tp (*common.TopicPartition): The partition whose state is being updated.
//	state (S): The new state to be assigned to the partition.
//
// Example:
//
//	tp := common.NewTopicPartition("topicB", 2)
//	ps := NewPartitionStates[string]()
//	ps.UpdateAndMoveToEnd(tp, "newState")
//	fmt.Println(ps.StateValue(tp)) // Output: newState, true
func (ps *PartitionStates[S]) UpdateAndMoveToEnd(tp *common.TopicPartition, state S) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	delete(ps.mapState, tp)
	ps.mapState[tp] = state
	ps.updateSize()
}

// Update updates the state of a partition.
//
// Parameters:
//
//	tp (*common.TopicPartition): The partition whose state is being updated.
//	state (S): The new state to be assigned to the partition.
//
// Example:
//
//	tp := common.NewTopicPartition("topicC", 3)
//	ps := NewPartitionStates[string]()
//	ps.Update(tp, "updatedState")
//	fmt.Println(ps.StateValue(tp)) // Output: updatedState, true
func (ps *PartitionStates[S]) Update(tp *common.TopicPartition, state S) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.mapState[tp] = state
	ps.updateSize()
}

// Remove removes the specified partition from the state.
//
// Parameters:
//
//	tp (*common.TopicPartition): The partition to be removed.
//
// Example:
//
//	tp := common.NewTopicPartition("topicD", 4)
//	ps := NewPartitionStates[string]()
//	ps.Remove(tp)
//	fmt.Println(ps.Contains(tp)) // Output: false
func (ps *PartitionStates[S]) Remove(tp *common.TopicPartition) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	delete(ps.mapState, tp)
	ps.updateSize()
}

// PartitionSet returns a set of all partitions currently in the state.
//
// Returns:
//
//	map[*common.TopicPartition]struct{}: A set of all partitions.
//
// Example:
//
//	tp1 := common.NewTopicPartition("topicE", 5)
//	tp2 := common.NewTopicPartition("topicF", 6)
//	ps := NewPartitionStates[string]()
//	ps.Update(tp1, "state1")
//	ps.Update(tp2, "state2")
//	set := ps.PartitionSet()
//	fmt.Println(len(set)) // Output: 2
func (ps *PartitionStates[S]) PartitionSet() map[*common.TopicPartition]struct{} {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	set := make(map[*common.TopicPartition]struct{}, len(ps.mapState))
	for tp := range ps.mapState {
		set[tp] = struct{}{}
	}
	return set
}

// Clear clears all partitions and resets the size.
//
// Example:
//
//	ps := NewPartitionStates[string]()
//	ps.Clear()
//	fmt.Println(ps.Size()) // Output: 0
func (ps *PartitionStates[S]) Clear() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.mapState = make(map[*common.TopicPartition]S)
	ps.size = 0
}

// Contains checks if a partition exists in the state.
//
// Parameters:
//
//	tp (*common.TopicPartition): The partition to check.
//
// Returns:
//
//	bool: true if the partition exists, false otherwise.
//
// Example:
//
//	tp := common.NewTopicPartition("topicG", 7)
//	ps := NewPartitionStates[string]()
//	ps.Update(tp, "state3")
//	fmt.Println(ps.Contains(tp)) // Output: true
func (ps *PartitionStates[S]) Contains(tp *common.TopicPartition) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	_, exists := ps.mapState[tp]
	return exists
}

// StateIterator returns a channel that iterates over the state values.
//
// Returns:
//
//	<-chan S: A channel for iterating over the state values.
//
// Example:
//
//	tp1 := common.NewTopicPartition("topicH", 8)
//	tp2 := common.NewTopicPartition("topicI", 9)
//	ps := NewPartitionStates[string]()
//	ps.Update(tp1, "state4")
//	ps.Update(tp2, "state5")
//	for state := range ps.StateIterator() {
//		fmt.Println(state) // Output: state4, state5
//	}
func (ps *PartitionStates[S]) StateIterator() <-chan S {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	ch := make(chan S)
	go func() {
		defer close(ch)
		for _, state := range ps.mapState {
			ch <- state
		}
	}()
	return ch
}

// ForEach applies a function to each partition and its corresponding state.
//
// Parameters:
//
//	f (func(*common.TopicPartition, S)): The function to apply to each partition and state.
//
// Example:
//
//	ps := NewPartitionStates[string]()
//	ps.Update(common.NewTopicPartition("topicJ", 10), "state6")
//	ps.ForEach(func(tp *common.TopicPartition, state string) {
//		fmt.Printf("Partition: %v, State: %v\n", tp.String(), state)
//	})
func (ps *PartitionStates[S]) ForEach(f func(*common.TopicPartition, S)) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for tp, state := range ps.mapState {
		f(tp, state)
	}
}

// PartitionStateMap returns an immutable map of all partition states.
//
// Returns:
//
//	map[*common.TopicPartition]S: A map of all partitions to their states.
//
// Example:
//
//	tp := common.NewTopicPartition("topicK", 11)
//	ps := NewPartitionStates[string]()
//	ps.Update(tp, "state7")
//	stateMap := ps.PartitionStateMap()
//	fmt.Println(stateMap[tp]) // Output: state7
func (ps *PartitionStates[S]) PartitionStateMap() map[*common.TopicPartition]S {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	copy := make(map[*common.TopicPartition]S, len(ps.mapState))
	for tp, state := range ps.mapState {
		copy[tp] = state
	}
	return copy
}

// PartitionStateValues returns a slice of state values in the order they were added.
//
// Returns:
//
//	[]S: A slice of state values.
//
// Example:
//
//	tp1 := common.NewTopicPartition("topicL", 12)
//	tp2 := common.NewTopicPartition("topicM", 13)
//	ps := NewPartitionStates[string]()
//	ps.Update(tp1, "state8")
//	ps.Update(tp2, "state9")
//	values := ps.PartitionStateValues()
//	fmt.Println(values) // Output: [state8, state9]
func (ps *PartitionStates[S]) PartitionStateValues() []S {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	values := make([]S, 0, len(ps.mapState))
	for _, state := range ps.mapState {
		values = append(values, state)
	}
	return values
}

// StateValue returns the state of a specific partition.
//
// Parameters:
//
//	tp (*common.TopicPartition): The partition whose state is being retrieved.
//
// Returns:
//
//	S: The state of the partition if it exists, zero value otherwise.
//	bool: true if the partition exists, false otherwise.
//
// Example:
//
//	tp := common.NewTopicPartition("topicN", 14)
//	ps := NewPartitionStates[string]()
//	ps.Update(tp, "state10")
//	state, exists := ps.StateValue(tp)
//	fmt.Println(state, exists) // Output: state10 true
func (ps *PartitionStates[S]) StateValue(tp *common.TopicPartition) (S, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	state, exists := ps.mapState[tp]
	return state, exists
}

// Size returns the number of partitions currently being tracked.
//
// Returns:
//
//	int: The number of partitions.
//
// Example:
//
//	ps := NewPartitionStates[string]()
//	ps.Update(common.NewTopicPartition("topicO", 15), "state11")
//	fmt.Println(ps.Size()) // Output: 1
func (ps *PartitionStates[S]) Size() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.size
}

// Set updates the state with the provided map of partitions to states.
//
// Parameters:
//
//	partitionToState (map[*common.TopicPartition]S): The map of partitions to their states.
//
// Example:
//
//	tp1 := common.NewTopicPartition("topicP", 16)
//	tp2 := common.NewTopicPartition("topicQ", 17)
//	ps := NewPartitionStates[string]()
//	ps.Set(map[*common.TopicPartition]S{
//		tp1: "state12",
//		tp2: "state13",
//	})
//	fmt.Println(ps.Size()) // Output: 2
func (ps *PartitionStates[S]) Set(partitionToState map[*common.TopicPartition]S) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.mapState = make(map[*common.TopicPartition]S, len(partitionToState))
	for tp, state := range partitionToState {
		ps.mapState[tp] = state
	}
	ps.updateSize()
}

// updateSize updates the size of the PartitionStates instance based on the current map state.
func (ps *PartitionStates[S]) updateSize() {
	ps.size = len(ps.mapState)
}

// PartitionState represents a state for a topic partition.
type PartitionState[S any] struct {
	TopicPartition *common.TopicPartition
	Value          S
}

// NewPartitionState creates a new PartitionState instance.
//
// Parameters:
//
//	tp (*common.TopicPartition): The topic partition associated with this state.
//	value (S): The value or state associated with the partition.
//
// Returns:
//
//	*PartitionState[S]: A pointer to the newly created PartitionState instance.
//
// Example:
//
//	tp := common.NewTopicPartition("topicR", 18)
//	ps := NewPartitionState(tp, "state14")
//	fmt.Println(ps.String()) // Output: PartitionState(topicR-18=state14)
func NewPartitionState[S any](tp *common.TopicPartition, value S) *PartitionState[S] {
	return &PartitionState[S]{
		TopicPartition: tp,
		Value:          value,
	}
}

// Equals checks if two PartitionState instances are equal. It compares both the
// TopicPartition and the Value using reflection.
//
// Parameters:
//
//	other (*PartitionState[S]): The other PartitionState instance to compare with.
//
// Returns:
//
//	bool: true if both PartitionState instances are equal, false otherwise.
//
// Example:
//
//	tp1 := common.NewTopicPartition("topicS", 19)
//	ps1 := NewPartitionState(tp1, "state15")
//	ps2 := NewPartitionState(tp1, "state15")
//	fmt.Println(ps1.Equals(ps2)) // Output: true
func (ps *PartitionState[S]) Equals(other *PartitionState[S]) bool {
	if other == nil {
		return false
	}

	// Check equality of TopicPartition using Equals method
	if !ps.TopicPartition.Equals(other.TopicPartition) {
		return false
	}

	// Check equality of Value using reflection
	return reflect.DeepEqual(ps.Value, other.Value)
}

// String returns a string representation of the PartitionState in the format
// "PartitionState(topic-partition=value)".
//
// Returns:
//
//	string: The string representation of the PartitionState.
//
// Example:
//
//	tp := common.NewTopicPartition("topicT", 20)
//	ps := NewPartitionState(tp, "state16")
//	fmt.Println(ps.String()) // Output: PartitionState(topicT-20=state16)
func (ps *PartitionState[S]) String() string {
	return fmt.Sprintf("PartitionState(%v=%v)", ps.TopicPartition, ps.Value)
}
