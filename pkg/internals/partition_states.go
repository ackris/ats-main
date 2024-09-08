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

// PartitionStates represents the state of partitions.
type PartitionStates[S any] struct {
	mu       sync.RWMutex
	mapState map[*common.TopicPartition]S
	size     int
}

// NewPartitionStates creates a new PartitionStates instance.
func NewPartitionStates[S any]() *PartitionStates[S] {
	return &PartitionStates[S]{
		mapState: make(map[*common.TopicPartition]S),
	}
}

// MoveToEnd moves a partition to the end of the order.
func (ps *PartitionStates[S]) MoveToEnd(tp *common.TopicPartition) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if state, ok := ps.mapState[tp]; ok {
		delete(ps.mapState, tp)
		ps.mapState[tp] = state
	}
}

// UpdateAndMoveToEnd updates the state of a partition and moves it to the end.
func (ps *PartitionStates[S]) UpdateAndMoveToEnd(tp *common.TopicPartition, state S) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	delete(ps.mapState, tp)
	ps.mapState[tp] = state
	ps.updateSize()
}

// Update updates the state of a partition.
func (ps *PartitionStates[S]) Update(tp *common.TopicPartition, state S) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.mapState[tp] = state
	ps.updateSize()
}

// Remove removes a partition from the state.
func (ps *PartitionStates[S]) Remove(tp *common.TopicPartition) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	delete(ps.mapState, tp)
	ps.updateSize()
}

// PartitionSet returns a set of all partitions.
func (ps *PartitionStates[S]) PartitionSet() map[*common.TopicPartition]struct{} {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	set := make(map[*common.TopicPartition]struct{}, len(ps.mapState))
	for tp := range ps.mapState {
		set[tp] = struct{}{}
	}
	return set
}

// Clear clears all partitions.
func (ps *PartitionStates[S]) Clear() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.mapState = make(map[*common.TopicPartition]S)
	ps.size = 0
}

// Contains checks if a partition exists.
func (ps *PartitionStates[S]) Contains(tp *common.TopicPartition) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	_, exists := ps.mapState[tp]
	return exists
}

// StateIterator returns an iterator for the partition states.
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

// ForEach applies a function to each partition and state.
func (ps *PartitionStates[S]) ForEach(f func(*common.TopicPartition, S)) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for tp, state := range ps.mapState {
		f(tp, state)
	}
}

// PartitionStateMap returns an immutable map of partition states.
func (ps *PartitionStates[S]) PartitionStateMap() map[*common.TopicPartition]S {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	copy := make(map[*common.TopicPartition]S, len(ps.mapState))
	for tp, state := range ps.mapState {
		copy[tp] = state
	}
	return copy
}

// PartitionStateValues returns the state values in order.
func (ps *PartitionStates[S]) PartitionStateValues() []S {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	values := make([]S, 0, len(ps.mapState))
	for _, state := range ps.mapState {
		values = append(values, state)
	}
	return values
}

// StateValue returns the state of a partition.
func (ps *PartitionStates[S]) StateValue(tp *common.TopicPartition) (S, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	state, exists := ps.mapState[tp]
	return state, exists
}

// Size returns the number of partitions currently being tracked.
func (ps *PartitionStates[S]) Size() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.size
}

// Set updates the builder to have the received map as its state.
func (ps *PartitionStates[S]) Set(partitionToState map[*common.TopicPartition]S) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.mapState = make(map[*common.TopicPartition]S, len(partitionToState))
	for tp, state := range partitionToState {
		ps.mapState[tp] = state
	}
	ps.updateSize()
}

func (ps *PartitionStates[S]) updateSize() {
	ps.size = len(ps.mapState)
}

// PartitionState represents a state for a topic partition.
type PartitionState[S any] struct {
	TopicPartition *common.TopicPartition
	Value          S
}

// NewPartitionState creates a new PartitionState instance.
func NewPartitionState[S any](tp *common.TopicPartition, value S) *PartitionState[S] {
	return &PartitionState[S]{
		TopicPartition: tp,
		Value:          value,
	}
}

// Equals checks if two PartitionState instances are equal.
// Uses reflection for generic comparison.
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

// String returns a string representation of the PartitionState.
func (ps *PartitionState[S]) String() string {
	return fmt.Sprintf("PartitionState(%v=%v)", ps.TopicPartition, ps.Value)
}
