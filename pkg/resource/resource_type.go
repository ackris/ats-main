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

package resource

import (
	"errors"
	"strings"
)

// ResourceType represents a type of resource which an ACL can be applied to.
type ResourceType byte

const (
	// UNRECOGNIZED represents any ResourceType which this client cannot understand.
	UNRECOGNIZED ResourceType = iota

	// ALL_RESOURCES matches any ResourceType.
	ALL_RESOURCES

	// TOPIC represents a Kafka topic.
	TOPIC

	// GROUP represents a consumer group.
	GROUP

	// CLUSTER represents the cluster as a whole.
	CLUSTER

	// TRANSACTIONAL_ID represents a transactional ID.
	TRANSACTIONAL_ID

	// DELEGATION_TOKEN represents a token ID.
	DELEGATION_TOKEN

	// USER represents a user principal.
	USER

	resourceTypeCount // Keep this at the end for the number of resource types
)

// resourceTypeNames maps ResourceType values to their string representations.
var resourceTypeNames = [resourceTypeCount]string{
	"UNRECOGNIZED",
	"ALL_RESOURCES",
	"TOPIC",
	"GROUP",
	"CLUSTER",
	"TRANSACTIONAL_ID",
	"DELEGATION_TOKEN",
	"USER",
}

// FromString parses the given string as an ACL resource type.
//
// Example usage:
//
//	resType, err := FromString("topic")
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Printf("Parsed ResourceType: %d\n", resType)
//	}
//
// This will output:
//
//	Parsed ResourceType: 2
func FromString(str string) (ResourceType, error) {
	if str == "" {
		return UNRECOGNIZED, errors.New("input string is empty")
	}

	// Normalize the input string to upper case for comparison
	str = strings.ToUpper(str)

	for i := range resourceTypeNames {
		if resourceTypeNames[i] == str {
			return ResourceType(i), nil
		}
	}
	return UNRECOGNIZED, errors.New("unknown resource type")
}

// FromCode returns the ResourceType with the provided code or UNRECOGNIZED if one cannot be found.
//
// Example usage:
//
//	code := byte(3)
//	resTypeFromCode := FromCode(code)
//	fmt.Printf("ResourceType from code %d: %d\n", code, resTypeFromCode)
//
// This will output:
//
//	ResourceType from code 3: 3
func FromCode(code byte) ResourceType {
	if code < byte(resourceTypeCount) {
		return ResourceType(code)
	}
	return UNRECOGNIZED
}

// Code returns the code of this resource type.
func (r ResourceType) Code() byte {
	return byte(r)
}

// IsUnrecognized returns whether this resource type is UNRECOGNIZED.
func (r ResourceType) IsUnrecognized() bool {
	return r >= resourceTypeCount
}

// String returns the string representation of the ResourceType.
func (r ResourceType) String() string {
	if int(r) < len(resourceTypeNames) {
		return resourceTypeNames[r]
	}
	return "UNKNOWN"
}
