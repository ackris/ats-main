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

import "fmt"

// AccessControlEntryData holds the data for access control entries.
type AccessControlEntryData struct {
	principal      string
	host           string
	operation      AclOperation
	permissionType AclPermissionType
}

// NewAccessControlEntryData creates a new AccessControlEntryData instance.
//
// Parameters:
// - principal: The principal (user or entity) associated with the entry.
// - host: The host associated with the entry.
// - operation: The ACL operation that this entry allows or denies.
// - permissionType: The type of permission (ALLOW, DENY, etc.) associated with the entry.
//
// Returns:
// A pointer to a new AccessControlEntryData instance.
//
// Example usage:
//
//	entry := NewAccessControlEntryData("user1", "localhost", OpRead, ALLOW)
func NewAccessControlEntryData(
	principal, host string,
	operation AclOperation,
	permissionType AclPermissionType) *AccessControlEntryData {
	return &AccessControlEntryData{
		principal:      principal,
		host:           host,
		operation:      operation,
		permissionType: permissionType,
	}
}

// Principal returns the principal associated with the entry.
//
// Returns:
// The principal as a string.
//
// Example usage:
//
//	principal := entry.Principal()
func (a *AccessControlEntryData) Principal() string {
	return a.principal
}

// Host returns the host associated with the entry.
//
// Returns:
// The host as a string.
//
// Example usage:
//
//	host := entry.Host()
func (a *AccessControlEntryData) Host() string {
	return a.host
}

// Operation returns the ACL operation associated with the entry.
//
// Returns:
// The ACL operation as an AclOperation type.
//
// Example usage:
//
//	operation := entry.Operation()
func (a *AccessControlEntryData) Operation() AclOperation {
	return a.operation
}

// PermissionType returns the permission type associated with the entry.
//
// Returns:
// The permission type as an AclPermissionType type.
//
// Example usage:
//
//	permissionType := entry.PermissionType()
func (a *AccessControlEntryData) PermissionType() AclPermissionType {
	return a.permissionType
}

// FindIndefiniteField checks for indefinite fields and returns a corresponding message.
//
// Returns:
// A string message indicating which field is indefinite, or an empty string if all fields are defined.
//
// Example usage:
//
//	message := entry.FindIndefiniteField()
func (a *AccessControlEntryData) FindIndefiniteField() string {
	switch {
	case a.principal == "":
		return "Principal is NULL"
	case a.host == "":
		return "Host is NULL"
	case a.operation == OpAny:
		return "Operation is ANY"
	case a.operation == OpUnknown:
		return "Operation is UNKNOWN"
	case a.permissionType == ANY:
		return "Permission type is ANY"
	case a.permissionType == UNKNOWN:
		return "Permission type is UNKNOWN"
	default:
		return ""
	}
}

// IsUnknown checks if there are any UNKNOWN components in the entry.
//
// Returns:
// A boolean indicating whether the entry contains any UNKNOWN components.
//
// Example usage:
//
//	if entry.IsUnknown() {
//	    fmt.Println("Entry contains unknown components.")
//	}
func (a *AccessControlEntryData) IsUnknown() bool {
	return a.operation == OpUnknown || a.permissionType == UNKNOWN
}

// String returns a string representation of the AccessControlEntryData.
//
// Returns:
// A formatted string that represents the AccessControlEntryData instance.
//
// Example usage:
//
//	fmt.Println(entry.String())
func (a *AccessControlEntryData) String() string {
	principalStr := "<any>"
	if a.principal != "" {
		principalStr = a.principal
	}
	hostStr := "<any>"
	if a.host != "" {
		hostStr = a.host
	}
	return fmt.Sprintf("(principal=%s, host=%s, operation=%d, permissionType=%d)",
		principalStr,
		hostStr,
		a.operation,
		a.permissionType)
}

// Equals checks if two AccessControlEntryData instances are equal.
//
// Parameters:
// - other: A pointer to another AccessControlEntryData instance to compare against.
//
// Returns:
// A boolean indicating whether the two instances are equal.
//
// Example usage:
//
//	if entry.Equals(otherEntry) {
//	    fmt.Println("Entries are equal.")
//	}
func (a *AccessControlEntryData) Equals(other *AccessControlEntryData) bool {
	if other == nil {
		return false
	}
	return a.principal == other.principal &&
		a.host == other.host &&
		a.operation == other.operation &&
		a.permissionType == other.permissionType
}

// HashCode returns a hash code for the AccessControlEntryData.
//
// Returns:
// An integer representing the hash code of the instance.
//
// Example usage:
//
//	hash := entry.HashCode()
func (a *AccessControlEntryData) HashCode() int {
	hash := 0
	for _, char := range a.principal {
		hash += int(char)
	}
	for _, char := range a.host {
		hash += int(char)
	}
	hash += int(a.operation) + int(a.permissionType)
	return hash
}
