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
)

// AccessControlEntryFilter represents a filter that matches access control entries.
type AccessControlEntryFilter struct {
	data *AccessControlEntryData
}

// Any represents a filter that matches any access control entry.
var Any = NewAccessControlEntryFilter("", "", OpAny, ANY)

// NewAccessControlEntryFilter creates an instance of an access control entry filter with the provided parameters.
// It validates the operation and permission type to ensure they are not UNKNOWN.
//
// Example usage:
//
//	filter := NewAccessControlEntryFilter("user1", "localhost", OpRead, ALLOW)
func NewAccessControlEntryFilter(principal, host string, operation AclOperation, permissionType AclPermissionType) *AccessControlEntryFilter {
	if operation == OpUnknown {
		panic("operation cannot be unknown")
	}
	if permissionType == UNKNOWN {
		panic("permissionType cannot be unknown")
	}
	return &AccessControlEntryFilter{
		data: NewAccessControlEntryData(principal, host, operation, permissionType),
	}
}

// Principal returns the principal or an empty string if none.
//
// Example usage:
//
//	principal := filter.Principal()
func (f *AccessControlEntryFilter) Principal() string {
	if f == nil || f.data == nil {
		return ""
	}
	return f.data.principal
}

// Host returns the host or an empty string if none. The value "*" means any host.
//
// Example usage:
//
//	host := filter.Host()
func (f *AccessControlEntryFilter) Host() string {
	if f == nil || f.data == nil {
		return ""
	}
	return f.data.host
}

// Operation returns the AclOperation.
//
// Example usage:
//
//	operation := filter.Operation()
func (f *AccessControlEntryFilter) Operation() AclOperation {
	if f == nil || f.data == nil {
		return OpUnknown // Return a known state if data is nil
	}
	return f.data.operation
}

// PermissionType returns the AclPermissionType.
//
// Example usage:
//
//	permissionType := filter.PermissionType()
func (f *AccessControlEntryFilter) PermissionType() AclPermissionType {
	if f == nil || f.data == nil {
		return UNKNOWN // Return a known state if data is nil
	}
	return f.data.permissionType
}

// IsUnknown returns true if there are any UNKNOWN components.
//
// Example usage:
//
//	isUnknown := filter.IsUnknown()
func (f *AccessControlEntryFilter) IsUnknown() bool {
	if f == nil || f.data == nil {
		return true // If data is nil, consider it unknown
	}
	return f.data.IsUnknown()
}

// Matches returns true if this filter matches the given AccessControlEntry.
//
// Example usage:
//
//	entry := NewAccessControlEntry("user1", "localhost", OpRead, ALLOW)
//	matches := filter.Matches(entry)
func (f *AccessControlEntryFilter) Matches(other *AccessControlEntry) bool {
	if other == nil {
		return false // Cannot match against a nil entry
	}
	if f.Principal() != "" && f.Principal() != other.Principal() {
		return false
	}
	if f.Host() != "" && f.Host() != other.Host() {
		return false
	}
	if f.Operation() != OpAny && f.Operation() != other.Operation() {
		return false
	}
	return f.PermissionType() == ANY || f.PermissionType() == other.PermissionType()
}

// MatchesAtMostOne returns true if this filter could only match one ACE.
//
// Example usage:
//
//	atMostOne := filter.MatchesAtMostOne()
func (f *AccessControlEntryFilter) MatchesAtMostOne() bool {
	return f.FindIndefiniteField() == ""
}

// FindIndefiniteField returns a string describing an ANY or UNKNOWN field, or an empty string
// if there is no such field.
//
// Example usage:
//
//	indefiniteField := filter.FindIndefiniteField()
func (f *AccessControlEntryFilter) FindIndefiniteField() string {
	if f == nil || f.data == nil {
		return "" // If data is nil, return empty
	}
	return f.data.FindIndefiniteField()
}

// String returns a string representation of the AccessControlEntryFilter.
//
// Example usage:
//
//	str := filter.String()
func (f *AccessControlEntryFilter) String() string {
	if f == nil {
		return "nil AccessControlEntryFilter" // Provide a clear indication of nil
	}
	if f.data == nil {
		return "AccessControlEntryFilter with nil data"
	}
	return fmt.Sprintf("%v", f.data)
}

// Equals returns true if this filter is equal to another filter.
//
// Example usage:
//
//	isEqual := filter.Equals(otherFilter)
func (f *AccessControlEntryFilter) Equals(other *AccessControlEntryFilter) bool {
	if other == nil {
		return false // Cannot be equal to a nil filter
	}
	if f == nil {
		return false // This filter is nil
	}
	return f.data.Equals(other.data)
}

// HashCode returns a hash code for the filter.
//
// Example usage:
//
//	hash := filter.HashCode()
func (f *AccessControlEntryFilter) HashCode() int {
	if f == nil || f.data == nil {
		return 0 // Return a known hash code for nil
	}
	return f.data.HashCode()
}
