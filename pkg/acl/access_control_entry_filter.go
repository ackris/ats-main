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
func (f *AccessControlEntryFilter) Principal() string {
	if f.data == nil {
		return ""
	}
	return f.data.principal
}

// Host returns the host or an empty string if none. The value "*" means any host.
func (f *AccessControlEntryFilter) Host() string {
	if f.data == nil {
		return ""
	}
	return f.data.host
}

// Operation returns the AclOperation.
func (f *AccessControlEntryFilter) Operation() AclOperation {
	if f.data == nil {
		return OpUnknown // Return a known state if data is nil
	}
	return f.data.operation
}

// PermissionType returns the AclPermissionType.
func (f *AccessControlEntryFilter) PermissionType() AclPermissionType {
	if f.data == nil {
		return UNKNOWN // Return a known state if data is nil
	}
	return f.data.permissionType
}

// IsUnknown returns true if there are any UNKNOWN components.
func (f *AccessControlEntryFilter) IsUnknown() bool {
	if f.data == nil {
		return true // If data is nil, consider it unknown
	}
	return f.data.IsUnknown()
}

// Matches returns true if this filter matches the given AccessControlEntry.
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
func (f *AccessControlEntryFilter) MatchesAtMostOne() bool {
	return f.FindIndefiniteField() == ""
}

// FindIndefiniteField returns a string describing an ANY or UNKNOWN field, or an empty string if there is no such field.
func (f *AccessControlEntryFilter) FindIndefiniteField() string {
	if f.data == nil {
		return "" // If data is nil, return empty
	}
	return f.data.FindIndefiniteField()
}

// String returns a string representation of the AccessControlEntryFilter.
func (f *AccessControlEntryFilter) String() string {
	if f.data == nil {
		return "nil AccessControlEntryFilter" // Provide a clear indication of nil
	}
	return fmt.Sprintf("%v", f.data)
}

// Equals returns true if this filter is equal to another filter.
func (f *AccessControlEntryFilter) Equals(other *AccessControlEntryFilter) bool {
	if other == nil {
		return false // Cannot be equal to a nil filter
	}
	return f.data.Equals(other.data)
}

// HashCode returns a hash code for the filter.
func (f *AccessControlEntryFilter) HashCode() int {
	if f.data == nil {
		return 0 // Return a known hash code for nil
	}
	return f.data.HashCode()
}
