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

type AclBinding struct {
	pattern *resource.ResourcePattern
	entry   *AccessControlEntry
}

func NewAclBinding(pattern *resource.ResourcePattern, entry *AccessControlEntry) (*AclBinding, error) {
	if pattern == nil {
		return nil, errors.New("pattern cannot be nil")
	}
	if entry == nil {
		return nil, errors.New("entry cannot be nil")
	}

	return &AclBinding{pattern: pattern, entry: entry}, nil
}

func (ab *AclBinding) IsUnknown() bool {
	return ab.pattern == nil || ab.entry == nil || ab.pattern.IsUnknown() || ab.entry.IsUnknown()
}

func (ab *AclBinding) Pattern() *resource.ResourcePattern {
	return ab.pattern
}

func (ab *AclBinding) Entry() *AccessControlEntry {
	return ab.entry
}

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

func (ab *AclBinding) String() string {
	return fmt.Sprintf("(pattern=%v, entry=%v)", ab.pattern.String(), ab.entry.String())
}

func (ab *AclBinding) Equals(other *AclBinding) bool {
	return other != nil && ab.pattern.Equals(other.pattern) && ab.entry.Equals(other.entry)
}

func (ab *AclBinding) Hash() int {
	if ab.pattern == nil || ab.entry == nil {
		return 0
	}
	return int(ab.pattern.HashCode()) ^ ab.entry.HashCode()
}
