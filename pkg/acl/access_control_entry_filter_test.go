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

func TestNewAccessControlEntryFilter(t *testing.T) {
	t.Run("valid_parameters", func(t *testing.T) {
		filter := NewAccessControlEntryFilter("user", "host", OpRead, ALLOW)
		if filter == nil {
			t.Errorf("NewAccessControlEntryFilter returned nil")
		}
	})

	t.Run("invalid_operation", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("NewAccessControlEntryFilter did not panic with invalid operation")
			}
		}()
		NewAccessControlEntryFilter("user", "host", OpUnknown, ALLOW)
	})

	t.Run("invalid_permission_type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("NewAccessControlEntryFilter did not panic with invalid permission type")
			}
		}()
		NewAccessControlEntryFilter("user", "host", OpRead, UNKNOWN)
	})
}

func TestAccessControlEntryFilter(t *testing.T) {
	filter := NewAccessControlEntryFilter("user", "host", OpRead, ALLOW)

	t.Run("principal", func(t *testing.T) {
		if filter.Principal() != "user" {
			t.Errorf("Principal() returned %s, expected %s", filter.Principal(), "user")
		}
	})

	t.Run("host", func(t *testing.T) {
		if filter.Host() != "host" {
			t.Errorf("Host() returned %s, expected %s", filter.Host(), "host")
		}
	})

	t.Run("operation", func(t *testing.T) {
		if filter.Operation() != OpRead {
			t.Errorf("Operation() returned %v, expected %v", filter.Operation(), OpRead)
		}
	})

	t.Run("permission_type", func(t *testing.T) {
		if filter.PermissionType() != ALLOW {
			t.Errorf("PermissionType() returned %v, expected %v", filter.PermissionType(), ALLOW)
		}
	})

	t.Run("is_unknown", func(t *testing.T) {
		if filter.IsUnknown() {
			t.Errorf("IsUnknown() returned true, expected false")
		}
	})

	t.Run("matches_at_most_one", func(t *testing.T) {
		if !filter.MatchesAtMostOne() {
			t.Errorf("MatchesAtMostOne() returned false, expected true")
		}
	})

	t.Run("find_indefinite_field", func(t *testing.T) {
		if field := filter.FindIndefiniteField(); field != "" {
			t.Errorf("FindIndefiniteField() returned %s, expected empty string", field)
		}
	})

	t.Run("string_representation", func(t *testing.T) {
		expected := "(principal=user, host=host, operation=3, permissionType=3)" // Adjust these values based on your AclOperation and AclPermissionType values
		if str := filter.String(); str != expected {
			t.Errorf("String() returned %s, expected %s", str, expected)
		}
	})

	t.Run("equals", func(t *testing.T) {
		sameFilter := NewAccessControlEntryFilter("user", "host", OpRead, ALLOW)
		if !filter.Equals(sameFilter) {
			t.Errorf("Equals() returned false, expected true")
		}
	})

	t.Run("hash_code", func(t *testing.T) {
		sameFilter := NewAccessControlEntryFilter("user", "host", OpRead, ALLOW)
		if filter.HashCode() != sameFilter.HashCode() {
			t.Errorf("HashCode() returned different values for equal filters")
		}
	})
}

func TestAccessControlEntryFilter_Matches(t *testing.T) {
	filter := NewAccessControlEntryFilter("user", "host", OpRead, ALLOW)

	t.Run("matches_entry", func(t *testing.T) {
		entry := NewAccessControlEntry("user", "host", OpRead, ALLOW)
		if !filter.Matches(entry) {
			t.Errorf("Matches() returned false, expected true")
		}
	})

	t.Run("does_not_match_entry", func(t *testing.T) {
		entry := NewAccessControlEntry("user", "host", OpWrite, ALLOW)
		if filter.Matches(entry) {
			t.Errorf("Matches() returned true, expected false")
		}
	})

	t.Run("does_not_match_nil_entry", func(t *testing.T) {
		if filter.Matches(nil) {
			t.Errorf("Matches() returned true when matching against nil")
		}
	})
}

func TestAccessControlEntryFilter_Nil(t *testing.T) {
	var filter *AccessControlEntryFilter

	t.Run("principal_nil", func(t *testing.T) {
		if filter.Principal() != "" {
			t.Errorf("Principal() returned %s, expected empty string", filter.Principal())
		}
	})

	t.Run("host_nil", func(t *testing.T) {
		if filter.Host() != "" {
			t.Errorf("Host() returned %s, expected empty string", filter.Host())
		}
	})

	t.Run("operation_nil", func(t *testing.T) {
		if filter.Operation() != OpUnknown {
			t.Errorf("Operation() returned %v, expected %v", filter.Operation(), OpUnknown)
		}
	})

	t.Run("permission_type_nil", func(t *testing.T) {
		if filter.PermissionType() != UNKNOWN {
			t.Errorf("PermissionType() returned %v, expected %v", filter.PermissionType(), UNKNOWN)
		}
	})

	t.Run("is_unknown_nil", func(t *testing.T) {
		if !filter.IsUnknown() {
			t.Errorf("IsUnknown() returned false, expected true")
		}
	})

	t.Run("find_indefinite_field_nil", func(t *testing.T) {
		if field := filter.FindIndefiniteField(); field != "" {
			t.Errorf("FindIndefiniteField() returned %s, expected empty string", field)
		}
	})

	t.Run("string_representation_nil", func(t *testing.T) {
		expected := "nil AccessControlEntryFilter"
		if str := filter.String(); str != expected {
			t.Errorf("String() returned %s, expected %s", str, expected)
		}
	})

	t.Run("equals_nil", func(t *testing.T) {
		if filter.Equals(nil) {
			t.Errorf("Equals() returned true when comparing with nil")
		}
	})

	t.Run("hash_code_nil", func(t *testing.T) {
		if filter.HashCode() != 0 {
			t.Errorf("HashCode() returned %d, expected 0", filter.HashCode())
		}
	})
}

func BenchmarkAccessControlEntryFilter_Matches(b *testing.B) {
	filter := NewAccessControlEntryFilter("user", "host", OpRead, ALLOW)
	entry := NewAccessControlEntry("user", "host", OpRead, ALLOW)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Matches(entry)
	}
}
