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
	"testing"
)

// TestNewAccessControlEntry tests the creation of AccessControlEntry.
func TestNewAccessControlEntry(t *testing.T) {
	tests := []struct {
		principal      string
		host           string
		operation      AclOperation
		permissionType AclPermissionType
		expectPanic    bool
	}{
		{"user1", "host1", OpRead, ALLOW, false},
		{"user2", "host2", OpWrite, DENY, false},
		{"", "host3", OpRead, ALLOW, true},           // Empty principal
		{"user4", "", OpRead, ALLOW, true},           // Empty host
		{"user5", "host5", OpAny, ALLOW, true},       // Operation ANY
		{"user6", "host6", OpUnknown, UNKNOWN, true}, // PermissionType UNKNOWN
	}

	for _, test := range tests {
		if test.expectPanic {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected panic for principal: %s, host: %s", test.principal, test.host)
				}
			}()
		}
		NewAccessControlEntry(test.principal, test.host, test.operation, test.permissionType)
	}
}

// TestAccessControlEntryMethods tests the methods of AccessControlEntry.
func TestAccessControlEntryMethods(t *testing.T) {
	ace := NewAccessControlEntry("user1", "host1", OpRead, ALLOW)

	if principal := ace.Principal(); principal != "user1" {
		t.Errorf("Expected principal to be 'user1', got '%s'", principal)
	}

	if host := ace.Host(); host != "host1" {
		t.Errorf("Expected host to be 'host1', got '%s'", host)
	}

	if operation := ace.Operation(); operation != OpRead {
		t.Errorf("Expected operation to be OpRead, got %d", operation)
	}

	if permissionType := ace.PermissionType(); permissionType != ALLOW {
		t.Errorf("Expected permissionType to be ALLOW, got %d", permissionType)
	}

	// Testing IsUnknown
	if ace.IsUnknown() {
		t.Error("Expected IsUnknown to return false")
	}

	// Create a new AccessControlEntry with UNKNOWN values
	unknownAce := NewAccessControlEntry("user2", "host2", OpUnknown, UNKNOWN)
	if !unknownAce.IsUnknown() {
		t.Error("Expected IsUnknown to return true")
	}
}

// TestAccessControlEntryString tests the string representation of AccessControlEntry.
func TestAccessControlEntryString(t *testing.T) {
	ace := NewAccessControlEntry("user1", "host1", OpRead, ALLOW)
	expected := "AccessControlEntry{Principal: user1, Host: host1, Operation: 3, PermissionType: 3}"
	if ace.String() != expected {
		t.Errorf("Expected string representation to be '%s', got '%s'", expected, ace.String())
	}
}

// TestAccessControlEntryEquals tests the Equals method.
func TestAccessControlEntryEquals(t *testing.T) {
	ace1 := NewAccessControlEntry("user1", "host1", OpRead, ALLOW)
	ace2 := NewAccessControlEntry("user1", "host1", OpRead, ALLOW) // Same parameters
	ace3 := NewAccessControlEntry("user2", "host2", OpWrite, DENY)

	if !ace1.Equals(ace2) {
		t.Error("Expected ace1 to be equal to ace2")
	}

	if ace1.Equals(ace3) {
		t.Error("Expected ace1 to not be equal to ace3")
	}

	if ace1.Equals(nil) {
		t.Error("Expected ace1 to not be equal to nil")
	}
}

// TestAccessControlEntryHashCode tests the HashCode method.
func TestAccessControlEntryHashCode(t *testing.T) {
	ace := NewAccessControlEntry("user1", "host1", OpRead, ALLOW)
	hashCode := ace.HashCode()
	if hashCode == 0 {
		t.Error("Expected HashCode to be non-zero")
	}
}
