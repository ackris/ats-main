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
		expectError    bool
	}{
		{"user1", "host1", OpRead, ALLOW, false},
		{"user2", "host2", OpWrite, DENY, false},
		{"", "host3", OpRead, ALLOW, true},           // Empty principal
		{"user4", "", OpRead, ALLOW, true},           // Empty host
		{"user5", "host5", OpAny, ALLOW, true},       // Operation ANY
		{"user6", "host6", OpUnknown, UNKNOWN, true}, // PermissionType UNKNOWN
	}

	for _, test := range tests {
		_, err := NewAccessControlEntry(test.principal, test.host, test.operation, test.permissionType)
		if (err != nil) != test.expectError {
			t.Errorf("Expected error: %v, got: %v for principal: %s, host: %s", test.expectError, (err != nil), test.principal, test.host)
		}
	}
}

// TestAccessControlEntryMethods tests the methods of AccessControlEntry.
func TestAccessControlEntryMethods(t *testing.T) {
	// Create valid AccessControlEntry
	ace, err := NewAccessControlEntry("user1", "host1", OpRead, ALLOW)
	if err != nil {
		t.Fatalf("Failed to create AccessControlEntry: %v", err)
	}

	t.Run("principal", func(t *testing.T) {
		if principal := ace.Principal(); principal != "user1" {
			t.Errorf("Expected principal to be 'user1', got '%s'", principal)
		}
	})

	t.Run("host", func(t *testing.T) {
		if host := ace.Host(); host != "host1" {
			t.Errorf("Expected host to be 'host1', got '%s'", host)
		}
	})

	t.Run("operation", func(t *testing.T) {
		if operation := ace.Operation(); operation != OpRead {
			t.Errorf("Expected operation to be OpRead, got %d", operation)
		}
	})

	t.Run("permission_type", func(t *testing.T) {
		if permissionType := ace.PermissionType(); permissionType != ALLOW {
			t.Errorf("Expected permissionType to be ALLOW, got %d", permissionType)
		}
	})

	// Testing IsUnknown with valid data
	t.Run("is_unknown", func(t *testing.T) {
		if ace.IsUnknown() {
			t.Error("Expected IsUnknown to return false")
		}
	})

	// Test case for known UNKNOWN value should be handled in a different context or test
	t.Run("cannot_create_with_UNKNOWN_permissionType", func(t *testing.T) {
		_, err := NewAccessControlEntry("user2", "host2", OpRead, UNKNOWN)
		if err == nil {
			t.Error("Expected error when creating AccessControlEntry with UNKNOWN permissionType, but got none")
		}
	})
}

// TestAccessControlEntryString tests the string representation of AccessControlEntry.
func TestAccessControlEntryString(t *testing.T) {
	ace, err := NewAccessControlEntry("user1", "host1", OpRead, ALLOW)
	if err != nil {
		t.Fatalf("Failed to create AccessControlEntry: %v", err)
	}

	expected := "AccessControlEntry{Principal: user1, Host: host1, Operation: 3, PermissionType: 3}"
	if ace.String() != expected {
		t.Errorf("Expected string representation to be '%s', got '%s'", expected, ace.String())
	}
}

// TestAccessControlEntryEquals tests the Equals method.
func TestAccessControlEntryEquals(t *testing.T) {
	ace1, err := NewAccessControlEntry("user1", "host1", OpRead, ALLOW)
	if err != nil {
		t.Fatalf("Failed to create AccessControlEntry: %v", err)
	}

	ace2, err := NewAccessControlEntry("user1", "host1", OpRead, ALLOW) // Same parameters
	if err != nil {
		t.Fatalf("Failed to create AccessControlEntry: %v", err)
	}

	ace3, err := NewAccessControlEntry("user2", "host2", OpWrite, DENY)
	if err != nil {
		t.Fatalf("Failed to create AccessControlEntry: %v", err)
	}

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
	ace, err := NewAccessControlEntry("user1", "host1", OpRead, ALLOW)
	if err != nil {
		t.Fatalf("Failed to create AccessControlEntry: %v", err)
	}

	hashCode := ace.HashCode()
	if hashCode == 0 {
		t.Error("Expected HashCode to be non-zero")
	}
}
