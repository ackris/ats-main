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
	"sync"
	"testing"
)

// TestNewTopicCollection tests the creation of a CombinedTopicCollection.
func TestNewTopicCollection(t *testing.T) {
	tests := []struct {
		name      string
		ids       []Uuid
		names     []string
		expectErr bool
	}{
		{"ValidCollection", []Uuid{NewUuid(1, 1), NewUuid(1, 2)}, []string{"topic1", "topic2"}, false},
		{"MismatchedLengths", []Uuid{NewUuid(1, 1)}, []string{"topic1", "topic2"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTopicCollection(tt.ids, tt.names)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// TestNewTopic tests the creation of a Topic.
func TestNewTopic(t *testing.T) {
	tests := []struct {
		testName  string // Renamed from 'name' to 'testName' to avoid conflict
		id        Uuid
		topicName string // Renamed from 'name' to 'topicName' to clarify its purpose
		expectErr bool
	}{
		{"ValidTopic", NewUuid(1, 1), "topic1", false},
		{"EmptyID", ZeroUUID, "topic1", true},
		{"EmptyName", NewUuid(1, 1), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := NewTopic(tt.id, tt.topicName)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// TestTopicIdCollection tests the TopicIdCollection functionality.
func TestTopicIdCollection(t *testing.T) {
	ids := []Uuid{NewUuid(1, 1), NewUuid(1, 2)}
	idCollection := NewTopicIdCollection(ids)

	// Test reading IDs
	readIds := idCollection.TopicIds()
	if len(readIds) != len(ids) {
		t.Errorf("expected %d IDs, got %d", len(ids), len(readIds))
	}

	// Test immutability
	idCollection.TopicIds()[0] = NewUuid(2, 2) // Attempt to modify
	readIds = idCollection.TopicIds()
	if readIds[0].Compare(NewUuid(1, 1)) != 0 {
		t.Errorf("expected first ID to be unchanged")
	}
}

// TestTopicNameCollection tests the TopicNameCollection functionality.
func TestTopicNameCollection(t *testing.T) {
	names := []string{"topic1", "topic2"}
	nameCollection := NewTopicNameCollection(names)

	// Test reading names
	readNames := nameCollection.TopicNames()
	if len(readNames) != len(names) {
		t.Errorf("expected %d names, got %d", len(names), len(readNames))
	}

	// Test immutability
	nameCollection.TopicNames()[0] = "topic3" // Attempt to modify
	readNames = nameCollection.TopicNames()
	if readNames[0] != "topic1" {
		t.Errorf("expected first name to be unchanged")
	}
}

// TestConcurrentAccess tests concurrent access to the collections.
func TestConcurrentAccess(t *testing.T) {
	ids := []Uuid{NewUuid(1, 1), NewUuid(1, 2)}
	idCollection := NewTopicIdCollection(ids)

	// Use a WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		readIds := idCollection.TopicIds()
		if len(readIds) != len(ids) {
			t.Errorf("expected %d IDs, got %d", len(ids), len(readIds))
		}
	}()

	go func() {
		defer wg.Done()
		readIds := idCollection.TopicIds()
		if len(readIds) != len(ids) {
			t.Errorf("expected %d IDs, got %d", len(ids), len(readIds))
		}
	}()

	wg.Wait() // Wait for both goroutines to finish
}
