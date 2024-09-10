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

package acl

import (
	"fmt"
	"testing"

	"github.com/ackris/ats-main/pkg/resource"
	"github.com/stretchr/testify/assert"
)

// TestNewAclBindingFilter tests the NewAclBindingFilter function
func TestNewAclBindingFilter(t *testing.T) {
	// Create real instances for patternFilter and entryFilter
	patternFilter, err := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
	assert.NoError(t, err)

	entryFilter := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)

	// Create an AclBindingFilter instance
	filter, err := NewAclBindingFilter(patternFilter, entryFilter)
	assert.NoError(t, err)
	assert.NotNil(t, filter)

	// Test invalid cases
	_, err = NewAclBindingFilter(nil, entryFilter)
	assert.Error(t, err)
	assert.Equal(t, "patternFilter cannot be nil", err.Error())

	_, err = NewAclBindingFilter(patternFilter, nil)
	assert.Error(t, err)
	assert.Equal(t, "entryFilter cannot be nil", err.Error())
}

// TestIsUnknown tests the IsUnknown method
func TestAclIsUnknown(t *testing.T) {
	// Create a valid pattern filter with a known resource type
	patternFilter, err := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
	if err != nil {
		t.Fatalf("Error creating ResourcePatternFilter: %v", err)
	}

	// Create a valid access control entry filter
	entryFilter := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)

	// Create the ACL binding filter
	filter, err := NewAclBindingFilter(patternFilter, entryFilter)
	if err != nil {
		t.Fatalf("Error creating AclBindingFilter: %v", err)
	}

	assert.False(t, filter.IsUnknown(), "Expected filter to be known")

	// Create a pattern filter that matches any resource pattern
	anyPatternFilter := resource.AnyResourcePatternFilter()

	// Print debug information
	fmt.Printf("anyPatternFilter: %v\n", anyPatternFilter)
	fmt.Printf("entryFilter: %v\n", entryFilter)

	// Create the ACL binding filter with the "any" pattern filter
	unknownFilter, err := NewAclBindingFilter(anyPatternFilter, entryFilter)
	if err != nil {
		t.Fatalf("Error creating AclBindingFilter with anyPatternFilter: %v", err)
	}

	// Print debug information
	fmt.Printf("unknownFilter: %v\n", unknownFilter)

	assert.False(t, unknownFilter.IsUnknown(), "Expected filter to be known")
}

// TestPatternFilter tests the PatternFilter method
func TestPatternFilter(t *testing.T) {
	patternFilter, _ := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
	entryFilter := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)
	filter, _ := NewAclBindingFilter(patternFilter, entryFilter)

	assert.Equal(t, patternFilter, filter.PatternFilter())
}

// TestEntryFilter tests the EntryFilter method
func TestEntryFilter(t *testing.T) {
	patternFilter, _ := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
	entryFilter := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)
	filter, _ := NewAclBindingFilter(patternFilter, entryFilter)

	assert.Equal(t, entryFilter, filter.EntryFilter())
}

// TestString tests the String method
func TestACLString(t *testing.T) {
	patternFilter, _ := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
	entryFilter := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)
	filter, _ := NewAclBindingFilter(patternFilter, entryFilter)

	expected := "(patternFilter=" + patternFilter.String() + ", entryFilter=" + entryFilter.String() + ")"
	assert.Equal(t, expected, filter.String())
}

// TestEquals tests the Equals method
func TestACLEquals(t *testing.T) {
	patternFilter1, _ := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
	entryFilter1 := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)
	filter1, _ := NewAclBindingFilter(patternFilter1, entryFilter1)

	patternFilter2, _ := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
	entryFilter2 := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)
	filter2, _ := NewAclBindingFilter(patternFilter2, entryFilter2)

	assert.True(t, filter1.Equals(filter2))

	// Modify one filter to create inequality
	filter2, _ = NewAclBindingFilter(patternFilter1, NewAccessControlEntryFilter("user2", "localhost", OpRead, ALLOW))
	assert.False(t, filter1.Equals(filter2))
}

// TestMatchesAtMostOne tests the MatchesAtMostOne method
func TestMatchesAtMostOne(t *testing.T) {
	patternFilter, _ := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
	entryFilter := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)
	filter, _ := NewAclBindingFilter(patternFilter, entryFilter)

	assert.True(t, filter.MatchesAtMostOne())

	// Modify patternFilter or entryFilter to create indefinite fields
	indefinitePatternFilter := resource.AnyResourcePatternFilter()
	indefiniteFilter, _ := NewAclBindingFilter(indefinitePatternFilter, entryFilter)
	assert.False(t, indefiniteFilter.MatchesAtMostOne())
}

// TestFindIndefiniteField tests the FindIndefiniteField method
func TestACLFindIndefiniteField(t *testing.T) {
	patternFilter, _ := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
	entryFilter := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)
	filter, _ := NewAclBindingFilter(patternFilter, entryFilter)

	assert.Equal(t, "", filter.FindIndefiniteField())

	// Modify patternFilter to create an indefinite field
	indefinitePatternFilter := resource.AnyResourcePatternFilter()
	indefiniteFilter, _ := NewAclBindingFilter(indefinitePatternFilter, entryFilter)
	assert.Equal(t, "Resource type is ALL_RESOURCES.", indefiniteFilter.FindIndefiniteField())
}

// TestMatches tests the Matches method
func TestMatches(t *testing.T) {
	// Create a valid pattern filter
	patternFilter, err := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
	if err != nil {
		t.Fatalf("Error creating ResourcePatternFilter: %v", err)
	}

	// Create a valid access control entry filter
	entryFilter := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)

	// Create the ACL binding filter
	filter, err := NewAclBindingFilter(patternFilter, entryFilter)
	if err != nil {
		t.Fatalf("Error creating AclBindingFilter: %v", err)
	}

	// Create a matching binding
	pattern, err := resource.NewResourcePattern(resource.TOPIC, "payments.", resource.PREFIXED)
	if err != nil {
		t.Fatalf("Error creating ResourcePattern: %v", err)
	}

	entry, err := NewAccessControlEntry("user1", "localhost", OpRead, ALLOW)
	if err != nil {
		t.Fatalf("Error creating AccessControlEntry: %v", err)
	}

	binding, err := NewAclBinding(pattern, entry)
	if err != nil {
		t.Fatalf("Error creating AclBinding: %v", err)
	}

	assert.True(t, filter.Matches(binding))

	// Create a non-matching binding
	nonMatchingPattern, err := resource.NewResourcePattern(resource.TOPIC, "otherpayments.", resource.PREFIXED)
	if err != nil {
		t.Fatalf("Error creating non-matching ResourcePattern: %v", err)
	}

	nonMatchingEntry, err := NewAccessControlEntry("user1", "localhost", OpRead, DENY)
	if err != nil {
		t.Fatalf("Error creating non-matching AccessControlEntry: %v", err)
	}

	nonMatchingBinding, err := NewAclBinding(nonMatchingPattern, nonMatchingEntry)
	if err != nil {
		t.Fatalf("Error creating non-matching AclBinding: %v", err)
	}

	assert.False(t, filter.Matches(nonMatchingBinding))
}

// TestHashCode tests the HashCode method
func TestACLHashCode(t *testing.T) {
	patternFilter, _ := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
	entryFilter := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)
	filter, _ := NewAclBindingFilter(patternFilter, entryFilter)

	expectedHash := (patternFilter.HashCode() << 32) | uint64(entryFilter.HashCode())
	assert.Equal(t, expectedHash, filter.HashCode())
}
