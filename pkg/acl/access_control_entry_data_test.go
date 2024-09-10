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

import "testing"

func TestNewAccessControlEntryData(t *testing.T) {
	entry := NewAccessControlEntryData("user1", "localhost", OpRead, ALLOW)

	if entry.Principal() != "user1" {
		t.Errorf("Expected principal user1, got %s", entry.Principal())
	}
	if entry.Host() != "localhost" {
		t.Errorf("Expected host localhost, got %s", entry.Host())
	}
	if entry.Operation() != OpRead {
		t.Errorf("Expected operation OpRead, got %d", entry.Operation())
	}
	if entry.PermissionType() != ALLOW {
		t.Errorf("Expected permission type ALLOW, got %d", entry.PermissionType())
	}
}

func TestFindIndefiniteField(t *testing.T) {
	tests := []struct {
		name     string
		entry    *AccessControlEntryData
		expected string
	}{
		{"Nil Principal", NewAccessControlEntryData("", "localhost", OpRead, ALLOW), "Principal is NULL"},
		{"Nil Host", NewAccessControlEntryData("user1", "", OpRead, ALLOW), "Host is NULL"},
		{"Any Operation", NewAccessControlEntryData("user1", "localhost", OpAny, ALLOW), "Operation is ANY"},
		{"Unknown Operation", NewAccessControlEntryData("user1", "localhost", OpUnknown, ALLOW), "Operation is UNKNOWN"},
		{"Any Permission Type", NewAccessControlEntryData("user1", "localhost", OpRead, ANY), "Permission type is ANY"},
		{"Unknown Permission Type", NewAccessControlEntryData("user1", "localhost", OpRead, UNKNOWN), "Permission type is UNKNOWN"},
		{"All Fields Valid", NewAccessControlEntryData("user1", "localhost", OpRead, ALLOW), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.entry.FindIndefiniteField()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestIsUnknown(t *testing.T) {
	tests := []struct {
		entry    *AccessControlEntryData
		expected bool
	}{
		{NewAccessControlEntryData("user1", "localhost", OpUnknown, ALLOW), true},
		{NewAccessControlEntryData("user1", "localhost", OpRead, UNKNOWN), true},
		{NewAccessControlEntryData("user1", "localhost", OpRead, ALLOW), false},
	}

	for _, tt := range tests {
		result := tt.entry.IsUnknown()
		if result != tt.expected {
			t.Errorf("Expected %v, got %v", tt.expected, result)
		}
	}
}

func TestString(t *testing.T) {
	entry := NewAccessControlEntryData("user1", "localhost", OpRead, ALLOW)
	expected := "(principal=user1, host=localhost, operation=3, permissionType=3)"
	if entry.String() != expected {
		t.Errorf("Expected %s, got %s", expected, entry.String())
	}
}

func TestEquals(t *testing.T) {
	entry1 := NewAccessControlEntryData("user1", "localhost", OpRead, ALLOW)
	entry2 := NewAccessControlEntryData("user1", "localhost", OpRead, ALLOW)
	entry3 := NewAccessControlEntryData("user2", "localhost", OpRead, ALLOW)

	if !entry1.Equals(entry2) {
		t.Errorf("Expected entries to be equal")
	}
	if entry1.Equals(entry3) {
		t.Errorf("Expected entries to be different")
	}
	if entry1.Equals(nil) {
		t.Errorf("Expected entry1 to not equal nil")
	}
}

func TestHashCode(t *testing.T) {
	entry := NewAccessControlEntryData("user1", "localhost", OpRead, ALLOW)
	expectedHash := entry.HashCode() // Store the expected hash
	if entry.HashCode() != expectedHash {
		t.Errorf("Hash code should remain consistent; got %d", entry.HashCode())
	}

	// Test with different entries
	entry2 := NewAccessControlEntryData("user2", "localhost", OpRead, ALLOW)
	if entry.HashCode() == entry2.HashCode() {
		t.Errorf("Expected different hash codes for different entries")
	}
}

func TestHashCodeConsistency(t *testing.T) {
	entry := NewAccessControlEntryData("user1", "localhost", OpRead, ALLOW)
	hash1 := entry.HashCode()
	hash2 := entry.HashCode()
	if hash1 != hash2 {
		t.Errorf("Hash code should be consistent; got %d and %d", hash1, hash2)
	}
}
