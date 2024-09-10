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

	"github.com/ackris/ats-main/pkg/resource"
	"go.uber.org/zap"
)

// AclBinding represents a binding between a resource pattern and an access control entry.
// It defines which resources are governed by which access control rules.
//
// Example usage:
//
//	pattern, err := resource.NewResourcePattern(resource.TOPIC, "test-topic", resource.LITERAL)
//	if err != nil {
//	    log.Fatalf("failed to create resource pattern: %v", err)
//	}
//
//	entry, err := NewAccessControlEntry("user1", "host1", OpRead, ALLOW)
//	if err != nil {
//	    log.Fatalf("failed to create access control entry: %v", err)
//	}
//
//	binding, err := NewAclBinding(pattern, entry)
//	if err != nil {
//	    log.Fatalf("failed to create ACL binding: %v", err)
//	}
//
//	fmt.Println(binding.String())
type AclBinding struct {
	pattern *resource.ResourcePattern
	entry   *AccessControlEntry
}

// NewAclBinding creates a new AclBinding instance. It requires both a valid resource pattern and access control entry.
//
// Parameters:
// - pattern: The resource pattern that defines the scope of the binding.
// - entry: The access control entry that specifies the rules for the binding.
//
// Returns:
// - A pointer to a newly created AclBinding instance, or an error if the pattern or entry is nil.
//
// Example usage:
//
//	pattern, err := resource.NewResourcePattern(resource.TOPIC, "test-topic", resource.LITERAL)
//	if err != nil {
//	    log.Fatalf("failed to create resource pattern: %v", err)
//	}
//
//	entry, err := NewAccessControlEntry("user1", "host1", OpRead, ALLOW)
//	if err != nil {
//	    log.Fatalf("failed to create access control entry: %v", err)
//	}
//
//	binding, err := NewAclBinding(pattern, entry)
//	if err != nil {
//	    log.Fatalf("failed to create ACL binding: %v", err)
//	}
func NewAclBinding(pattern *resource.ResourcePattern, entry *AccessControlEntry) (*AclBinding, error) {
	if pattern == nil {
		return nil, errors.New("pattern cannot be nil")
	}
	if entry == nil {
		return nil, errors.New("entry cannot be nil")
	}

	return &AclBinding{pattern: pattern, entry: entry}, nil
}

// IsUnknown determines if the AclBinding is considered unknown. An AclBinding is unknown if:
// - The pattern or entry is nil, or
// - The pattern or entry is marked as unknown.
//
// Returns:
// - True if the AclBinding is unknown, false otherwise.
//
// Example usage:
//
//	if binding.IsUnknown() {
//	    fmt.Println("The ACL binding is unknown.")
//	} else {
//	    fmt.Println("The ACL binding is known.")
//	}
func (ab *AclBinding) IsUnknown() bool {
	return ab.pattern == nil || ab.entry == nil || ab.pattern.IsUnknown() || ab.entry.IsUnknown()
}

// Pattern returns the resource pattern associated with the AclBinding.
//
// Returns:
// - The resource pattern of the AclBinding.
//
// Example usage:
//
//	pattern := binding.Pattern()
//	fmt.Println("Resource pattern:", pattern)
func (ab *AclBinding) Pattern() *resource.ResourcePattern {
	return ab.pattern
}

// Entry returns the access control entry associated with the AclBinding.
//
// Returns:
// - The access control entry of the AclBinding.
//
// Example usage:
//
//	entry := binding.Entry()
//	fmt.Println("Access control entry:", entry)
func (ab *AclBinding) Entry() *AccessControlEntry {
	return ab.entry
}

// ToFilter converts the AclBinding into an AclBindingFilter. If either the pattern or entry is nil,
// or if there is an error creating the pattern or entry filters, it returns nil.
//
// Returns:
// - A pointer to the created AclBindingFilter, or nil if conversion fails.
//
// Example usage:
//
//	filter := binding.ToFilter()
//	if filter != nil {
//	    fmt.Println("Successfully created ACL binding filter.")
//	} else {
//	    fmt.Println("Failed to create ACL binding filter.")
//	}
func (ab *AclBinding) ToFilter() *AclBindingFilter {
	if ab.pattern == nil || ab.entry == nil {
		return nil
	}

	patternFilter, err := ab.pattern.ToFilter()
	if err != nil {
		zap.L().Error("failed to create pattern filter", zap.Error(err))
		return nil
	}

	entryFilter := ab.entry.ToFilter()

	filter, err := NewAclBindingFilter(patternFilter, entryFilter)
	if err != nil {
		zap.L().Error("failed to create ACL binding filter", zap.Error(err))
		return nil
	}

	return filter
}

// String returns a string representation of the AclBinding. This representation includes
// the string representations of both the resource pattern and access control entry.
//
// Returns:
// - A string representing the AclBinding.
//
// Example usage:
//
//	fmt.Println("ACL Binding:", binding.String())
func (ab *AclBinding) String() string {
	return fmt.Sprintf("(pattern=%v, entry=%v)", ab.pattern.String(), ab.entry.String())
}

// Equals compares the AclBinding to another AclBinding for equality. Two AclBindings are considered
// equal if they have the same resource pattern and access control entry.
//
// Parameters:
// - other: The other AclBinding to compare to.
//
// Returns:
// - True if the AclBindings are equal, false otherwise.
//
// Example usage:
//
//	anotherBinding := /* create another AclBinding instance */
//	if binding.Equals(anotherBinding) {
//	    fmt.Println("The ACL bindings are equal.")
//	} else {
//	    fmt.Println("The ACL bindings are not equal.")
//	}
func (ab *AclBinding) Equals(other *AclBinding) bool {
	return other != nil && ab.pattern.Equals(other.pattern) && ab.entry.Equals(other.entry)
}

// Hash returns a hash code for the AclBinding. The hash code is computed based on the hash codes
// of the resource pattern and access control entry. If either is nil, the hash code is 0.
//
// Returns:
// - The hash code of the AclBinding.
//
// Example usage:
//
//	hashCode := binding.Hash()
//	fmt.Println("Hash code of ACL Binding:", hashCode)
func (ab *AclBinding) Hash() int {
	if ab.pattern == nil || ab.entry == nil {
		return 0
	}
	return int(ab.pattern.HashCode()) ^ ab.entry.HashCode()
}
