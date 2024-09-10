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

// Unit tests for AclOperation
func TestAclOperation(t *testing.T) {
	t.Run("String representation", func(t *testing.T) {
		testCases := []struct {
			name     string
			op       AclOperation
			expected string
		}{
			{"OpUnknown", OpUnknown, "UNKNOWN"},
			{"OpAny", OpAny, "ANY"},
			{"OpAll", OpAll, "ALL"},
			{"OpRead", OpRead, "READ"},
			{"OpWrite", OpWrite, "WRITE"},
			{"OpCreate", OpCreate, "CREATE"},
			{"OpDelete", OpDelete, "DELETE"},
			{"OpAlter", OpAlter, "ALTER"},
			{"OpDescribe", OpDescribe, "DESCRIBE"},
			{"OpClusterAction", OpClusterAction, "CLUSTER_ACTION"},
			{"OpDescribeConfigs", OpDescribeConfigs, "DESCRIBE_CONFIGS"},
			{"OpAlterConfigs", OpAlterConfigs, "ALTER_CONFIGS"},
			{"OpIdempotentWrite", OpIdempotentWrite, "IDEMPOTENT_WRITE"},
			{"OpCreateTokens", OpCreateTokens, "CREATE_TOKENS"},
			{"OpDescribeTokens", OpDescribeTokens, "DESCRIBE_TOKENS"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.op.String() != tc.expected {
					t.Errorf("expected %s, got %s", tc.expected, tc.op.String())
				}
			})
		}
	})

	t.Run("FromString", func(t *testing.T) {
		testCases := []struct {
			name      string
			input     string
			expected  AclOperation
			expectErr bool
		}{
			{"Valid operation", "READ", OpRead, false},
			{"Case-insensitive", "wRiTe", OpWrite, false},
			{"Whitespace", "  CREATE  ", OpCreate, false},
			{"Unknown operation", "INVALID", OpUnknown, true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				op, err := FromAOPString(tc.input)
				if (err != nil) != tc.expectErr {
					t.Errorf("expected error: %v, got: %v", tc.expectErr, err)
				}
				if op != tc.expected {
					t.Errorf("expected %s, got %s", tc.expected, op)
				}
			})
		}
	})

	t.Run("FromCode", func(t *testing.T) {
		testCases := []struct {
			name     string
			code     byte
			expected AclOperation
		}{
			{"Valid code", byte(OpRead), OpRead},
			{"Out-of-range code", 0xFF, OpUnknown},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				op := FromCode(tc.code)
				if op != tc.expected {
					t.Errorf("expected %s, got %s", tc.expected, op)
				}
			})
		}
	})

	t.Run("IsUnknown", func(t *testing.T) {
		testCases := []struct {
			name     string
			op       AclOperation
			expected bool
		}{
			{"OpUnknown", OpUnknown, true},
			{"OpRead", OpRead, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.op.IsUnknown() != tc.expected {
					t.Errorf("expected %t, got %t", tc.expected, tc.op.IsUnknown())
				}
			})
		}
	})
}
