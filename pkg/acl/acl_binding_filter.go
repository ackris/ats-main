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
	"errors"
	"fmt"
	"sync"

	resource "github.com/ackris/ats-main/pkg/resource"
)

// AclBindingFilter represents a filter for matching AclBinding objects.
// It combines a ResourcePatternFilter and an AccessControlEntryFilter to determine if an AclBinding meets certain criteria.
//
// Example usage:
//
//	patternFilter, _ := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
//	entryFilter := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)
//	filter, _ := NewAclBindingFilter(patternFilter, entryFilter)
//	binding := NewAclBinding(pattern, entry)
//	matches := filter.Matches(binding) // true if the binding matches the filters
type AclBindingFilter struct {
	patternFilter *resource.ResourcePatternFilter
	entryFilter   *AccessControlEntryFilter
	hash          uint64
	once          sync.Once
}

// NewAclBindingFilter creates a new instance of AclBindingFilter.
//
// Parameters:
// - patternFilter: The resource pattern filter to use. Must not be nil.
// - entryFilter: The access control entry filter to use. Must not be nil.
//
// Returns:
// - A pointer to a new AclBindingFilter instance.
// - An error if any of the filters are nil.
//
// Example usage:
// patternFilter, _ := resource.NewResourcePatternFilter(resource.TOPIC, "payments.", resource.PREFIXED)
// entryFilter := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)
// filter, err := NewAclBindingFilter(patternFilter, entryFilter)
//
//	if err != nil {
//	    // Handle error
//	}
func NewAclBindingFilter(
	patternFilter *resource.ResourcePatternFilter,
	entryFilter *AccessControlEntryFilter) (*AclBindingFilter, error) {
	if patternFilter == nil {
		return nil, errors.New("patternFilter cannot be nil")
	}
	if entryFilter == nil {
		return nil, errors.New("entryFilter cannot be nil")
	}
	return &AclBindingFilter{patternFilter: patternFilter, entryFilter: entryFilter}, nil
}

// IsUnknown checks if the filter has any UNKNOWN components.
//
// Returns:
// - true if the patternFilter or entryFilter is unknown, false otherwise.
//
// Example usage:
//
//	unknown := filter.IsUnknown() // true if either filter component is unknown
func (f *AclBindingFilter) IsUnknown() bool {
	return f.patternFilter.IsUnknown() || f.entryFilter.IsUnknown()
}

// PatternFilter returns the resource pattern filter used by this AclBindingFilter.
//
// Returns:
// - The resource pattern filter.
//
// Example usage:
//
//	patternFilter := filter.PatternFilter()
func (f *AclBindingFilter) PatternFilter() *resource.ResourcePatternFilter {
	return f.patternFilter
}

// EntryFilter returns the access control entry filter used by this AclBindingFilter.
//
// Returns:
// - The access control entry filter.
//
// Example usage:
//
//	entryFilter := filter.EntryFilter()
func (f *AclBindingFilter) EntryFilter() *AccessControlEntryFilter {
	return f.entryFilter
}

// String returns a string representation of the AclBindingFilter.
//
// Returns:
// - A string describing the filter.
//
// Example usage:
//
//	fmt.Println(filter.String()) // Outputs the string representation of the filter
func (f *AclBindingFilter) String() string {
	return fmt.Sprintf("(patternFilter=%v, entryFilter=%v)", f.patternFilter, f.entryFilter)
}

// Equals checks if two AclBindingFilter instances are equal.
//
// Parameters:
// - other: The other AclBindingFilter to compare with.
//
// Returns:
// - true if both filters are equal, false otherwise.
//
// Example usage:
//
//	equal := filter.Equals(otherFilter) // true if filters are the same
func (f *AclBindingFilter) Equals(other *AclBindingFilter) bool {
	if other == nil {
		return false
	}
	return f.patternFilter.Equal(other.patternFilter) && f.entryFilter.Equals(other.entryFilter)
}

// MatchesAtMostOne checks if the resource and entry filters can only match one ACE.
//
// Returns:
// - true if both patternFilter and entryFilter can only match one ACE, false otherwise.
//
// Example usage:
//
//	atMostOne := filter.MatchesAtMostOne() // true if the filters are restrictive
func (f *AclBindingFilter) MatchesAtMostOne() bool {
	return f.patternFilter.MatchesAtMostOne() && f.entryFilter.MatchesAtMostOne()
}

// FindIndefiniteField returns a description of any indefinite (ANY or UNKNOWN) field in the filters,
// or an empty string if there is none.
//
// Returns:
// - A string describing any indefinite field, or an empty string if none.
//
// Example usage:
//
//	indefinite := filter.FindIndefiniteField() // Returns a description of indefinite fields
func (f *AclBindingFilter) FindIndefiniteField() string {
	if indefinite := f.patternFilter.FindIndefiniteField(); indefinite != "" {
		return indefinite
	}
	return f.entryFilter.FindIndefiniteField()
}

// Matches checks if the filter matches the given AclBinding.
//
// Parameters:
// - binding: The AclBinding to match against the filter.
//
// Returns:
// - true if the filter matches the binding's resource and entry, false otherwise.
//
// Example usage:
//
//	binding := NewAclBinding(pattern, entry)
//	matches := filter.Matches(binding) // true if the binding matches the filter
func (f *AclBindingFilter) Matches(binding *AclBinding) bool {
	return f.patternFilter.Matches(binding.Pattern()) && f.entryFilter.Matches(binding.Entry())
}

// HashCode returns a hash code for the AclBindingFilter.
//
// Returns:
// - A 64-bit unsigned integer representing the hash code.
//
// Example usage:
//
//	hashCode := filter.HashCode() // Hash code for the filter
func (f *AclBindingFilter) HashCode() uint64 {
	f.once.Do(func() {
		patternHash := f.patternFilter.HashCode()
		entryHash := uint64(f.entryFilter.HashCode())
		f.hash = (patternHash << 32) | entryHash
	})
	return f.hash
}
