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

import "strings"

// AclPermissionType represents whether an ACL grants or denies permissions.
type AclPermissionType int8

const (
	// UNKNOWN represents any AclPermissionType which this client cannot understand,
	// perhaps because this client is too old.
	UNKNOWN AclPermissionType = iota

	// ANY matches any AclPermissionType in a filter.
	ANY

	// DENY disallows access.
	DENY

	// ALLOW grants access.
	ALLOW

	// maxAclPermissionType is the maximum value of AclPermissionType
	maxAclPermissionType
)

var codeToValue = [maxAclPermissionType]string{
	"UNKNOWN",
	"ANY",
	"DENY",
	"ALLOW",
}

// String returns the string representation of the AclPermissionType.
func (apt AclPermissionType) String() string {
	if apt >= 0 && int(apt) < len(codeToValue) {
		return codeToValue[apt]
	}
	return "UNKNOWN"
}

// FromString parses the given string as an ACL permission.
func FromString(str string) AclPermissionType {
	str = strings.ToUpper(str)
	for i, s := range codeToValue {
		if s == str {
			return AclPermissionType(i)
		}
	}
	return UNKNOWN
}

// Code returns the code of this permission type.
func (apt AclPermissionType) Code() AclPermissionType {
	return apt
}

// IsUnknown returns true if this permission type is UNKNOWN.
func (apt AclPermissionType) IsUnknown() bool {
	return apt == UNKNOWN
}
