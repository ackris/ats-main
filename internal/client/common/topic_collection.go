package common

import (
	"errors"
	"sync"
)

// TopicCollection represents a collection of topics defined by name or ID.
// It provides methods to retrieve topic IDs and names.
type TopicCollection interface {
	// TopicIds returns a read-only slice of topic IDs.
	TopicIds() []Uuid
	// TopicNames returns a read-only slice of topic names.
	TopicNames() []string
}

// CombinedTopicCollection represents a collection of topics defined by both name and ID.
// It holds instances of TopicIdCollection and TopicNameCollection.
type CombinedTopicCollection struct {
	topicIds   *TopicIdCollection
	topicNames *TopicNameCollection
}

// NewTopicCollection creates a new CombinedTopicCollection with the given topic IDs and names.
//
// Parameters:
//   - ids: A slice of Uuid representing topic IDs.
//   - names: A slice of string representing topic names.
//
// Returns:
//   - A pointer to a CombinedTopicCollection instance.
//   - An error if the lengths of ids and names slices are not equal.
//
// Example:
//
//	ids := []Uuid{NewUuid(1, 1), NewUuid(1, 2)}
//	names := []string{"topic1", "topic2"}
//	topicCollection, err := NewTopicCollection(ids, names)
//	if err != nil {
//	    // Handle the error
//	    return
//	}
//	// Use the topicCollection
func NewTopicCollection(ids []Uuid, names []string) (*CombinedTopicCollection, error) {
	if len(ids) != len(names) {
		return nil, errors.New("topic IDs and names must have the same length")
	}

	return &CombinedTopicCollection{
		topicIds:   NewTopicIdCollection(ids),
		topicNames: NewTopicNameCollection(names),
	}, nil
}

// TopicIdCollection represents a collection of topics defined by their topic ID.
// It uses a sync.RWMutex to ensure thread-safety.
type TopicIdCollection struct {
	topicIds []Uuid
	mu       sync.RWMutex
}

// NewTopicIdCollection creates a new TopicIdCollection with the given topic IDs.
//
// Parameters:
//   - topicIds: A slice of Uuid representing topic IDs.
//
// Returns:
//   - A pointer to a TopicIdCollection instance.
//
// Example:
//
//	ids := []Uuid{NewUuid(1, 1), NewUuid(1, 2)}
//	idCollection := NewTopicIdCollection(ids)
//	// Use the idCollection
func NewTopicIdCollection(topicIds []Uuid) *TopicIdCollection {
	return &TopicIdCollection{
		topicIds: append([]Uuid{}, topicIds...), // Create a copy to prevent external modifications
	}
}

// TopicIds returns a read-only slice of topic IDs.
//
// Returns:
//   - A copy of the topic IDs slice.
//
// Example:
//
//	ids := idCollection.TopicIds()
//	// Use the ids slice
func (c *TopicIdCollection) TopicIds() []Uuid {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]Uuid{}, c.topicIds...) // Return a copy to ensure immutability
}

// TopicNameCollection represents a collection of topics defined by their topic name.
// It uses a sync.RWMutex to ensure thread-safety.
type TopicNameCollection struct {
	topicNames []string
	mu         sync.RWMutex
}

// NewTopicNameCollection creates a new TopicNameCollection with the given topic names.
//
// Parameters:
//   - topicNames: A slice of string representing topic names.
//
// Returns:
//   - A pointer to a TopicNameCollection instance.
//
// Example:
//
//	names := []string{"topic1", "topic2"}
//	nameCollection := NewTopicNameCollection(names)
//	// Use the nameCollection
func NewTopicNameCollection(topicNames []string) *TopicNameCollection {
	return &TopicNameCollection{
		topicNames: append([]string{}, topicNames...), // Create a copy to prevent external modifications
	}
}

// TopicNames returns a read-only slice of topic names.
//
// Returns:
//   - A copy of the topic names slice.
//
// Example:
//
//	names := nameCollection.TopicNames()
//	// Use the names slice
func (c *TopicNameCollection) TopicNames() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]string{}, c.topicNames...) // Return a copy to ensure immutability
}

// NewTopic creates a new topic with the given ID and name.
//
// Parameters:
//   - id: A Uuid representing the topic ID.
//   - name: A string representing the topic name.
//
// Returns:
//   - A pointer to a Topic instance.
//   - An error if the topic ID is empty or the topic name is empty.
//
// Example:
//
//	topic, err := NewTopic(NewUuid(1, 1), "topic1")
//	if err != nil {
//	    // Handle the error
//	    return
//	}
//	// Use the topic
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
//
// Returns:
//   - The topic ID.
//
// Example:
//
//	id := topic.ID()
//	// Use the id
func (t *Topic) ID() Uuid {
	return t.id
}

// Name returns the topic name.
//
// Returns:
//   - The topic name.
//
// Example:
//
//	name := topic.Name()
//	// Use the name
func (t *Topic) Name() string {
	return t.name
}
