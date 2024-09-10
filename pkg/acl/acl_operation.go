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
	"strings"
)

// AclOperation represents an operation which an ACL grants or denies permission to perform.
type AclOperation byte

// Define ACL operations as constants.
const (
	OpUnknown         AclOperation = iota // 0
	OpAny                                 // 1
	OpAll                                 // 2
	OpRead                                // 3
	OpWrite                               // 4
	OpCreate                              // 5
	OpDelete                              // 6
	OpAlter                               // 7
	OpDescribe                            // 8
	OpClusterAction                       // 9
	OpDescribeConfigs                     // 10
	OpAlterConfigs                        // 11
	OpIdempotentWrite                     // 12
	OpCreateTokens                        // 13
	OpDescribeTokens                      // 14
)

// operationNames maps AclOperation to their string representation.
var operationNames = map[AclOperation]string{
	OpUnknown:         "UNKNOWN",
	OpAny:             "ANY",
	OpAll:             "ALL",
	OpRead:            "READ",
	OpWrite:           "WRITE",
	OpCreate:          "CREATE",
	OpDelete:          "DELETE",
	OpAlter:           "ALTER",
	OpDescribe:        "DESCRIBE",
	OpClusterAction:   "CLUSTER_ACTION",
	OpDescribeConfigs: "DESCRIBE_CONFIGS",
	OpAlterConfigs:    "ALTER_CONFIGS",
	OpIdempotentWrite: "IDEMPOTENT_WRITE",
	OpCreateTokens:    "CREATE_TOKENS",
	OpDescribeTokens:  "DESCRIBE_TOKENS",
}

// String returns the string representation of the AclOperation.
func (op AclOperation) String() string {
	if name, exists := operationNames[op]; exists {
		return name
	}
	return "UNKNOWN"
}

// FromAOPString parses a string and returns the corresponding AclOperation.
//
// Example:
//
//	op, err := FromAOPString("READ")
//	if err != nil {
//		// Handle error
//	}
//	fmt.Println(op) // Output: READ
func FromAOPString(str string) (AclOperation, error) {
	str = strings.TrimSpace(str) // Trim whitespace for better matching
	for op, name := range operationNames {
		if strings.EqualFold(name, str) {
			return op, nil
		}
	}
	return OpUnknown, errors.New("unknown AclOperation: " + str)
}

// FromCode returns the AclOperation corresponding to the given byte code.
//
// Example:
//
//	op := FromCode(byte(OpRead))
//	fmt.Println(op) // Output: READ
func FromCode(code byte) AclOperation {
	// Ensure the code is within the valid range of defined operations
	if code < byte(OpUnknown) || code > byte(OpDescribeTokens) {
		return OpUnknown
	}
	return AclOperation(code)
}

// IsUnknown checks if the AclOperation is UNKNOWN.
//
// Example:
//
//	op := FromCode(0xFF)
//	if op.IsUnknown() {
//		fmt.Println("Unknown operation")
//	}
func (op AclOperation) IsUnknown() bool {
	return op == OpUnknown
}
