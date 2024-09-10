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

// AclBindingFilter is a filter that can match AclBinding objects.
type AclBindingFilter struct {
	patternFilter *resource.ResourcePatternFilter
	entryFilter   *AccessControlEntryFilter
	hash          uint64
	once          sync.Once
}

// NewAclBindingFilter creates an instance of AclBindingFilter with the provided parameters.
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

// IsUnknown returns true if this filter has any UNKNOWN components.
func (f *AclBindingFilter) IsUnknown() bool {
	return f.patternFilter.IsUnknown() || f.entryFilter.IsUnknown()
}

// PatternFilter returns the resource pattern filter.
func (f *AclBindingFilter) PatternFilter() *resource.ResourcePatternFilter {
	return f.patternFilter
}

// EntryFilter returns the access control entry filter.
func (f *AclBindingFilter) EntryFilter() *AccessControlEntryFilter {
	return f.entryFilter
}

// String returns a string representation of the AclBindingFilter.
func (f *AclBindingFilter) String() string {
	return fmt.Sprintf("(patternFilter=%v, entryFilter=%v)", f.patternFilter, f.entryFilter)
}

// Equals checks if two AclBindingFilters are equal.
func (f *AclBindingFilter) Equals(other *AclBindingFilter) bool {
	if other == nil {
		return false
	}
	return f.patternFilter.Equal(other.patternFilter) && f.entryFilter.Equals(other.entryFilter)
}

// MatchesAtMostOne returns true if the resource and entry filters can only match one ACE.
func (f *AclBindingFilter) MatchesAtMostOne() bool {
	return f.patternFilter.MatchesAtMostOne() && f.entryFilter.MatchesAtMostOne()
}

// FindIndefiniteField returns a string describing an ANY or UNKNOWN field, or nil if there is no such field.
func (f *AclBindingFilter) FindIndefiniteField() string {
	if indefinite := f.patternFilter.FindIndefiniteField(); indefinite != "" {
		return indefinite
	}
	return f.entryFilter.FindIndefiniteField()
}

// Matches returns true if the resource filter matches the binding's resource and the entry filter matches the binding's entry.
func (f *AclBindingFilter) Matches(binding *AclBinding) bool {
	return f.patternFilter.Matches(binding.Pattern()) && f.entryFilter.Matches(binding.Entry())
}

// HashCode returns the hash code of the AclBindingFilter.
func (f *AclBindingFilter) HashCode() uint64 {
	f.once.Do(func() {
		patternHash := f.patternFilter.HashCode()
		entryHash := uint64(f.entryFilter.HashCode())
		f.hash = (patternHash << 32) | entryHash
	})
	return f.hash
}
