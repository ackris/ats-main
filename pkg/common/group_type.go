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

package common

import "strings"

// GroupType represents the different types of groups.
type GroupType int

const (
	// UNSPECIFIED represents an unspecified or unknown group type.
	UNSPECIFIED GroupType = iota
	// CONSUMER represents a consumer group type.
	CONSUMER
	// CLASSIC represents a classic group type.
	CLASSIC
	// SHARE represents a share group type.
	SHARE
)

// stringToGroupType maps string representations to GroupType values.
var stringToGroupType = map[string]GroupType{
	"unspecified": UNSPECIFIED,
	"consumer":    CONSUMER,
	"classic":     CLASSIC,
	"share":       SHARE,
}

// String returns the string representation of the GroupType.
func (g GroupType) String() string {
	return [...]string{"Unspecified", "Consumer", "Classic", "Share"}[g]
}

// ParseGroupType converts a string to a GroupType in a case-insensitive manner.
// Returns UNSPECIFIED if the input is empty or does not match any known group type.
//
// Examples:
//
//	ParseGroupType("consumer") // Returns CONSUMER_GROUP
//	ParseGroupType("CLASSIC")  // Returns CLASSIC_GROUP
//	ParseGroupType("unknown")  // Returns UNSPECIFIED
func ParseGroupType(name string) GroupType {
	if name == "" {
		return UNSPECIFIED
	}

	// Convert to lowercase and lookup in the map
	if groupType, exists := stringToGroupType[strings.ToLower(name)]; exists {
		return groupType
	}
	return UNSPECIFIED
}
