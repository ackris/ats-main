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
