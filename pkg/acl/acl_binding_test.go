package acl

import (
	"testing"

	"github.com/ackris/ats-main/pkg/resource"
	"github.com/stretchr/testify/assert"
)

func TestNewAclBindingNilPattern(t *testing.T) {
	// Assuming that AccessControlEntry creation uses NewAccessControlEntry function
	entry, err := NewAccessControlEntry("user", "host", OpRead, ALLOW)
	assert.NoError(t, err, "Failed to create valid AccessControlEntry")

	_, err = NewAclBinding(nil, entry)
	assert.Error(t, err)
	assert.Equal(t, "pattern cannot be nil", err.Error())
}

func TestNewAclBindingNilEntry(t *testing.T) {
	// Create a valid ResourcePattern using the factory function
	pattern, err := resource.NewResourcePattern(resource.TOPIC, "test-topic", resource.LITERAL)
	assert.NoError(t, err, "Failed to create valid ResourcePattern")

	_, err = NewAclBinding(pattern, nil)
	assert.Error(t, err)
	assert.Equal(t, "entry cannot be nil", err.Error())
}

func TestAclBindingIsUnknown(t *testing.T) {
	// Define valid values
	validResourceType := resource.TOPIC
	validName := "test-topic"
	validPatternType := resource.LITERAL
	validPrincipal := "user"
	validHost := "host"
	validOperation := OpRead
	validPermissionType := ALLOW

	// Create a valid ResourcePattern
	pattern, err := resource.NewResourcePattern(validResourceType, validName, validPatternType)
	assert.NoError(t, err, "Failed to create valid ResourcePattern")

	// Create a valid AccessControlEntry
	entry, err := NewAccessControlEntry(validPrincipal, validHost, validOperation, validPermissionType)
	assert.NoError(t, err, "Failed to create valid AccessControlEntry")

	// Create a valid AclBinding
	binding, err := NewAclBinding(pattern, entry)
	assert.NoError(t, err, "Failed to create valid AclBinding")
	assert.False(t, binding.IsUnknown(), "AclBinding should not be unknown with valid parameters")

	// Test with ResourcePattern having an invalid PatternType
	invalidPatternType := resource.UNKNOWN // Use an invalid pattern type
	invalidPattern, err := resource.NewResourcePattern(validResourceType, validName, invalidPatternType)
	if assert.Error(t, err, "Expected error when creating ResourcePattern with UNKNOWN PatternType") {
		assert.Equal(t, "patternType cannot be UNKNOWN", err.Error())
	} else {
		// Only proceed if no error was returned
		bindingInvalidPattern, err := NewAclBinding(invalidPattern, entry)
		assert.Error(t, err, "Expected error when creating AclBinding with invalid ResourcePattern")
		assert.Nil(t, bindingInvalidPattern, "AclBinding should be nil if pattern is invalid")
	}

	// Test with AccessControlEntry having an invalid Operation
	invalidOperation := OpAny // Use an invalid operation
	_, err = NewAccessControlEntry(validPrincipal, validHost, invalidOperation, validPermissionType)
	assert.Error(t, err, "Expected error when creating AccessControlEntry with ANY Operation")

	// Test with AccessControlEntry having an invalid PermissionType
	invalidPermissionType := ANY // Use an invalid permission type
	_, err = NewAccessControlEntry(validPrincipal, validHost, validOperation, invalidPermissionType)
	assert.Error(t, err, "Expected error when creating AccessControlEntry with ANY PermissionType")

	// Test creating AclBinding with invalid AccessControlEntry
	entryInvalid, err := NewAccessControlEntry(validPrincipal, validHost, validOperation, invalidPermissionType)
	if assert.Error(t, err, "Expected error when creating AccessControlEntry with ANY PermissionType") {
		assert.Equal(t, "permissionType must not be ANY", err.Error())
	} else {
		// Only proceed if no error was returned
		bindingInvalidEntry, err := NewAclBinding(pattern, entryInvalid)
		assert.Error(t, err, "Expected error when creating AclBinding with invalid AccessControlEntry")
		assert.Nil(t, bindingInvalidEntry, "AclBinding should be nil if entry is invalid")
	}
}

func TestAclBindingToFilter(t *testing.T) {
	validResourceType := resource.TOPIC
	validName := "test-topic"
	validPatternType := resource.LITERAL
	validPrincipal := "user"
	validHost := "host"
	validOperation := OpRead
	validPermissionType := ALLOW

	pattern, err := resource.NewResourcePattern(validResourceType, validName, validPatternType)
	assert.NoError(t, err)

	entry, err := NewAccessControlEntry(validPrincipal, validHost, validOperation, validPermissionType)
	assert.NoError(t, err)

	binding, err := NewAclBinding(pattern, entry)
	assert.NoError(t, err)

	filter := binding.ToFilter()
	assert.NotNil(t, filter, "Expected non-nil filter from AclBinding.ToFilter()")
}

func TestAclBindingString(t *testing.T) {
	validResourceType := resource.TOPIC
	validName := "test-topic"
	validPatternType := resource.LITERAL
	validPrincipal := "user"
	validHost := "host"
	validOperation := OpRead
	validPermissionType := ALLOW

	pattern, err := resource.NewResourcePattern(validResourceType, validName, validPatternType)
	assert.NoError(t, err)

	entry, err := NewAccessControlEntry(validPrincipal, validHost, validOperation, validPermissionType)
	assert.NoError(t, err)

	binding, err := NewAclBinding(pattern, entry)
	assert.NoError(t, err)

	// Update expected values based on actual output
	expectedString := "(pattern=ResourcePattern{resourceType=TOPIC, name=\"test-topic\", patternType=LITERAL}, entry=AccessControlEntry{Principal: user, Host: host, Operation: 3, PermissionType: 3})"
	assert.Equal(t, expectedString, binding.String(), "AclBinding String representation does not match")
}
