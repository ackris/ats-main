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

// AccessControlEntry represents an access control entry.
type AccessControlEntry struct {
	data *AccessControlEntryData
}

// NewAccessControlEntry creates an instance of an access control entry with the provided parameters.
//
// Parameters:
//   - principal: The principal for this entry. Must not be empty.
//   - host: The host for this entry. Must not be empty.
//   - operation: The AclOperation for this entry. Must not be AclOperation.ANY.
//   - permissionType: The AclPermissionType for this entry. Must not be AclPermissionType.ANY.
//
// Returns:
//   - A pointer to the created AccessControlEntry instance.
//
// Panics:
//   - If principal is empty.
//   - If host is empty.
//   - If operation is AclOperation.ANY.
//   - If permissionType is AclPermissionType.ANY.
//
// Example:
//
//	ace := NewAccessControlEntry("user1", "host1", OpRead, ALLOW)
func NewAccessControlEntry(
	principal, host string,
	operation AclOperation,
	permissionType AclPermissionType) *AccessControlEntry {
	if principal == "" {
		panic("principal cannot be empty")
	}
	if host == "" {
		panic("host cannot be empty")
	}
	if operation == OpAny {
		panic("operation must not be ANY")
	}
	if permissionType == ANY {
		panic("permissionType must not be ANY")
	}

	return &AccessControlEntry{
		data: NewAccessControlEntryData(principal, host, operation, permissionType),
	}
}

// Principal returns the principal for this entry.
//
// Returns:
//   - The principal for this entry.
//
// Example:
//
//	principal := ace.Principal()
func (ace *AccessControlEntry) Principal() string {
	return ace.data.principal
}

// Host returns the host or `*` for all hosts.
//
// Returns:
//   - The host for this entry, or `*` if it represents all hosts.
//
// Example:
//
//	host := ace.Host()
func (ace *AccessControlEntry) Host() string {
	return ace.data.host
}

// Operation returns the AclOperation. This method will never return AclOperation.ANY.
//
// Returns:
//   - The AclOperation for this entry. It will never return AclOperation.ANY.
//
// Example:
//
//	operation := ace.Operation()
func (ace *AccessControlEntry) Operation() AclOperation {
	return ace.data.operation
}

// PermissionType returns the AclPermissionType. This method will never return AclPermissionType.ANY.
//
// Returns:
//   - The AclPermissionType for this entry. It will never return AclPermissionType.ANY.
//
// Example:
//
//	permissionType := ace.PermissionType()
func (ace *AccessControlEntry) PermissionType() AclPermissionType {
	return ace.data.permissionType
}

// ToFilter creates a filter which matches only this AccessControlEntry.
//
// Returns:
//   - A new AccessControlEntryFilter that matches only this AccessControlEntry.
//
// Example:
//
//	filter := ace.ToFilter()
func (ace *AccessControlEntry) ToFilter() *AccessControlEntryFilter {
	return NewAccessControlEntryFilter(ace.Principal(), ace.Host(), ace.Operation(), ace.PermissionType())
}

// IsUnknown returns true if this AccessControlEntry has any UNKNOWN components.
//
// Returns:
//   - true if either the operation or permissionType is UNKNOWN, false otherwise.
//
// Example:
//
//	isUnknown := ace.IsUnknown()
func (ace *AccessControlEntry) IsUnknown() bool {
	return ace.data.operation == OpUnknown || ace.data.permissionType == UNKNOWN
}

// String returns a string representation of the AccessControlEntry.
//
// Returns:
//   - A string representation of the AccessControlEntry in the format:
//     AccessControlEntry{Principal: <principal>, Host: <host>, Operation: <operation>, PermissionType: <permissionType>}
//
// Example:
//
//	aceString := ace.String()
func (ace *AccessControlEntry) String() string {
	return fmt.Sprintf("AccessControlEntry{Principal: %s, Host: %s, Operation: %d, PermissionType: %d}",
		ace.Principal(), ace.Host(), ace.Operation(), ace.PermissionType())
}

// Equals checks if two AccessControlEntry instances are equal.
//
// Parameters:
//   - other: The AccessControlEntry instance to compare with.
//
// Returns:
//   - true if the AccessControlEntry instances have the same principal, host, operation, and permissionType, false otherwise.
//
// Example:
//
//	equal := ace.Equals(otherEntry)
func (ace *AccessControlEntry) Equals(other *AccessControlEntry) bool {
	if other == nil {
		return false
	}
	data1 := ace.data
	data2 := other.data
	return data1.principal == data2.principal &&
		data1.host == data2.host &&
		data1.operation == data2.operation &&
		data1.permissionType == data2.permissionType
}

// HashCode returns a hash code for the AccessControlEntry.
//
// Returns:
//   - A hash code for the AccessControlEntry based on its principal, host, operation, and permissionType.
//   - If the data is nil, it returns 0.
//
// Example:
//
//	hashCode := ace.HashCode()
func (ace *AccessControlEntry) HashCode() int {
	if ace.data == nil {
		return 0 // Handle the case where data might be nil
	}
	return ace.data.HashCode()
}
