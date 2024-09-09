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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"

	cr "github.com/ackris/ats-main/pkg/common"
)

// MockClusterResourceListener is a mock implementation of ClusterResourceListener.
type MockClusterResourceListener struct {
	mock.Mock
}

func (m *MockClusterResourceListener) OnUpdate(clusterResource *cr.ClusterResource) {
	m.Called(clusterResource)
}

func TestMaybeAdd(t *testing.T) {
	// Initialize a zap logger for testing
	logger := zaptest.NewLogger(t)

	// Create a new ClusterResourceListeners instance
	crListeners := NewClusterResourceListeners(logger)

	// Create a mock listener
	mockListener := new(MockClusterResourceListener)

	// Add the mock listener
	crListeners.MaybeAdd(mockListener)

	// Ensure the listener was added
	crListeners.mu.RLock()
	defer crListeners.mu.RUnlock()
	assert.Len(t, crListeners.listeners, 1)
	assert.Equal(t, mockListener, crListeners.listeners[0])
}

func TestMaybeAddInvalidType(t *testing.T) {
	// Initialize a zap logger for testing
	logger := zaptest.NewLogger(t)

	// Create a new ClusterResourceListeners instance
	crListeners := NewClusterResourceListeners(logger)

	// Attempt to add an invalid type
	crListeners.MaybeAdd("invalid")

	// Ensure no listeners were added
	crListeners.mu.RLock()
	defer crListeners.mu.RUnlock()
	assert.Len(t, crListeners.listeners, 0)
}

func TestMaybeAddAll(t *testing.T) {
	// Initialize a zap logger for testing
	logger := zaptest.NewLogger(t)

	// Create a new ClusterResourceListeners instance
	crListeners := NewClusterResourceListeners(logger)

	// Create mock listeners
	mockListener1 := new(MockClusterResourceListener)
	mockListener2 := new(MockClusterResourceListener)

	// List containing valid and invalid items
	candidates := []interface{}{
		mockListener1,
		"invalid",
		mockListener2,
	}

	// Add all candidates
	crListeners.MaybeAddAll(candidates)

	// Ensure both valid listeners were added
	crListeners.mu.RLock()
	defer crListeners.mu.RUnlock()
	assert.Len(t, crListeners.listeners, 2)
	assert.Contains(t, crListeners.listeners, mockListener1)
	assert.Contains(t, crListeners.listeners, mockListener2)
}

func TestOnUpdate(t *testing.T) {
	// Initialize a zap logger for testing
	logger := zaptest.NewLogger(t)

	// Create a new ClusterResourceListeners instance
	crListeners := NewClusterResourceListeners(logger)

	// Create a mock listener and add it
	mockListener := new(MockClusterResourceListener)
	crListeners.MaybeAdd(mockListener)

	// Set up the expectation
	clusterResource := cr.NewClusterResource("test-cluster")
	mockListener.On("OnUpdate", clusterResource).Once()

	// Call OnUpdate
	crListeners.OnUpdate(clusterResource)

	// Assert that OnUpdate was called with the correct argument
	mockListener.AssertExpectations(t)
}
