package common

import (
	"errors"
	"sync"
)

// TopicCollection represents a collection of topics defined by name or ID.
type TopicCollection interface {
	TopicIds() []Uuid
	TopicNames() []string
}

// TopicIdCollection represents a collection of topics defined by their topic ID.
type TopicIdCollection struct {
	topicIds []Uuid
	mu       sync.RWMutex
}

// NewTopicIdCollection creates a new TopicIdCollection with the given topic IDs.
func NewTopicIdCollection(topicIds []Uuid) *TopicIdCollection {
	return &TopicIdCollection{
		topicIds: append([]Uuid{}, topicIds...), // Create a copy to prevent external modifications
	}
}

// TopicIds returns a read-only slice of topic IDs.
func (c *TopicIdCollection) TopicIds() []Uuid {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]Uuid{}, c.topicIds...) // Return a copy to ensure immutability
}

// TopicNameCollection represents a collection of topics defined by their topic name.
type TopicNameCollection struct {
	topicNames []string
	mu         sync.RWMutex
}

// NewTopicNameCollection creates a new TopicNameCollection with the given topic names.
func NewTopicNameCollection(topicNames []string) *TopicNameCollection {
	return &TopicNameCollection{
		topicNames: append([]string{}, topicNames...), // Create a copy to prevent external modifications
	}
}

// TopicNames returns a read-only slice of topic names.
func (c *TopicNameCollection) TopicNames() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]string{}, c.topicNames...) // Return a copy to ensure immutability
}

// NewTopic creates a new topic with the given ID and name.
func NewTopic(id Uuid, name string) (*Topic, error) {
	if id.Compare(ZeroUUID) == 0 {
		return nil, errors.New("topic ID cannot be empty")
	}
	if name == "" {
		return nil, errors.New("topic name cannot be empty")
	}
	return &Topic{
		id:   id,
		name: name,
	}, nil
}

// Topic represents a single topic with an ID and name.
type Topic struct {
	id   Uuid
	name string
}

// ID returns the topic ID.
func (t *Topic) ID() Uuid {
	return t.id
}

// Name returns the topic name.
func (t *Topic) Name() string {
	return t.name
}
