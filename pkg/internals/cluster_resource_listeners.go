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
	"sync"

	"go.uber.org/zap"

	cr "github.com/ackris/ats-main/pkg/common"
)

// ClusterResourceListener defines an interface for receiving updates about ClusterResource changes.
//
// Implementers of this interface should provide logic to handle updates in the OnUpdate method.
//
// Example usage:
//
// type MyListener struct {}
//
//	func (l *MyListener) OnUpdate(clusterResource *cr.ClusterResource) {
//	    // Handle the cluster resource update here
//	}
//
// myListener := &MyListener{}
// listenersManager := internals.NewClusterResourceListeners(logger)
// listenersManager.MaybeAdd(myListener)
type ClusterResourceListener interface {
	// OnUpdate is a callback method that a user can implement to get updates for ClusterResource.
	OnUpdate(clusterResource *cr.ClusterResource)
}

// ClusterResourceListeners manages a list of ClusterResourceListener instances.
//
// It allows adding listeners, either individually or in bulk, and notifies them of updates.
//
// Example usage:
//
// logger := zap.NewExample() // or use a more suitable logger configuration
// listenersManager := internals.NewClusterResourceListeners(logger)
//
// myListener := &MyListener{}
// listenersManager.MaybeAdd(myListener)
//
// anotherListener := &AnotherListener{}
// listenersManager.MaybeAddAll([]interface{}{anotherListener, &YetAnotherListener{}})
//
// clusterResource := cr.NewClusterResource("cluster-id")
// listenersManager.OnUpdate(clusterResource)
type ClusterResourceListeners struct {
	listeners []ClusterResourceListener
	mu        sync.RWMutex // ensures thread-safe access to listeners
	logger    *zap.Logger  // logger for structured logging
}

// NewClusterResourceListeners creates a new instance of ClusterResourceListeners with the provided logger.
//
// Parameters:
// - logger: A zap.Logger instance used for logging warnings and errors.
//
// Returns:
// - *ClusterResourceListeners: A new instance of ClusterResourceListeners.
//
// Example:
//
// logger := zap.NewExample()
// listenersManager := internals.NewClusterResourceListeners(logger)
func NewClusterResourceListeners(logger *zap.Logger) *ClusterResourceListeners {
	return &ClusterResourceListeners{
		listeners: make([]ClusterResourceListener, 0),
		logger:    logger,
	}
}

// MaybeAdd adds a ClusterResourceListener to the internal list if the candidate implements the interface.
//
// If the candidate does not implement ClusterResourceListener, a warning is logged.
//
// Parameters:
// - candidate: The object to add, which should implement ClusterResourceListener.
//
// Example:
//
// myListener := &MyListener{}
// listenersManager.MaybeAdd(myListener)
//
// invalidListener := "not-a-listener"
// listenersManager.MaybeAdd(invalidListener)
func (c *ClusterResourceListeners) MaybeAdd(candidate interface{}) {
	listener, ok := candidate.(ClusterResourceListener)
	if !ok {
		c.logger.Warn("Candidate does not implement ClusterResourceListener", zap.Any("candidate", candidate))
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.listeners = append(c.listeners, listener)
}

// MaybeAddAll adds all items from the list that implement ClusterResourceListener to the internal list.
//
// Each candidate that does not implement ClusterResourceListener will have a warning logged.
//
// Parameters:
// - candidateList: A slice of objects to add, which should implement ClusterResourceListener.
//
// Example:
//
// validListener1 := &MyListener{}
// validListener2 := &AnotherListener{}
// invalidListener := "not-a-listener"
//
// listenersManager.MaybeAddAll([]interface{}{validListener1, invalidListener, validListener2})
func (c *ClusterResourceListeners) MaybeAddAll(candidateList []interface{}) {
	validListeners := make([]ClusterResourceListener, 0, len(candidateList))
	for _, candidate := range candidateList {
		listener, ok := candidate.(ClusterResourceListener)
		if !ok {
			c.logger.Warn("Candidate does not implement ClusterResourceListener", zap.Any("candidate", candidate))
			continue
		}
		validListeners = append(validListeners, listener)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.listeners = append(c.listeners, validListeners...)
}

// OnUpdate sends the updated cluster metadata to all registered ClusterResourceListener instances.
//
// Parameters:
// - cluster: The updated cluster metadata to notify listeners of.
//
// Example:
//
// clusterResource := cr.NewClusterResource("cluster-id")
// listenersManager.OnUpdate(clusterResource)
func (c *ClusterResourceListeners) OnUpdate(cluster *cr.ClusterResource) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, listener := range c.listeners {
		listener.OnUpdate(cluster)
	}
}
