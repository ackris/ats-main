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
	"testing"
)

func TestNewResourcePatternFilter(t *testing.T) {
	_, err := NewResourcePatternFilter(UNRECOGNIZED, "topic1", LITERAL)
	if err == nil {
		t.Errorf("NewResourcePatternFilter should return an error for UNRECOGNIZED resourceType")
	}

	_, err = NewResourcePatternFilter(TOPIC, "topic1", UNKNOWN)
	if err == nil {
		t.Errorf("NewResourcePatternFilter should return an error for UNKNOWN patternType")
	}

	filter, err := NewResourcePatternFilter(TOPIC, "topic1", LITERAL)
	if err != nil {
		t.Errorf("NewResourcePatternFilter returned an unexpected error: %v", err)
	}
	if filter == nil {
		t.Errorf("NewResourcePatternFilter returned a nil filter")
	}
}

func TestAnyResourcePatternFilter(t *testing.T) {
	filter := AnyResourcePatternFilter()
	if filter.resourceType != ALL_RESOURCES {
		t.Errorf("AnyResourcePatternFilter should return a filter with resourceType ALL_RESOURCES")
	}
	if filter.name != "" {
		t.Errorf("AnyResourcePatternFilter should return a filter with empty name")
	}
	if filter.patternType != ANY {
		t.Errorf("AnyResourcePatternFilter should return a filter with patternType ANY")
	}
}

func TestResourcePatternFilter_Matches(t *testing.T) {
	// Test matching with LITERAL
	filter, _ := NewResourcePatternFilter(TOPIC, "payments.received", LITERAL)
	pattern, _ := NewResourcePattern(TOPIC, "payments.received", LITERAL)
	if !filter.Matches(pattern) {
		t.Errorf("Expected filter to match the pattern: filter=%s, pattern=%s", filter, pattern)
	}

	// Test matching with PREFIXED
	filter, _ = NewResourcePatternFilter(TOPIC, "payments.", PREFIXED)
	pattern, _ = NewResourcePattern(TOPIC, "payments.received", LITERAL)
	if !filter.Matches(pattern) {
		t.Errorf("Expected filter to match the prefixed pattern: filter=%s, pattern=%s", filter, pattern)
	}

	// Test wildcard matching
	filter, _ = NewResourcePatternFilter(TOPIC, "*", LITERAL)
	pattern, _ = NewResourcePattern(TOPIC, "payments.received", LITERAL)
	if !filter.Matches(pattern) {
		t.Errorf("Expected filter to match the wildcard pattern: filter=%s, pattern=%s", filter, pattern)
	}

	// Test non-matching resource type
	filter, _ = NewResourcePatternFilter(TOPIC, "payments.received", LITERAL)
	pattern, _ = NewResourcePattern(GROUP, "payments.received", LITERAL)
	if filter.Matches(pattern) {
		t.Errorf("Expected filter to not match a pattern with different resourceType: filter=%s, pattern=%s", filter, pattern)
	}

	// Test non-matching pattern type
	filter, _ = NewResourcePatternFilter(TOPIC, "payments.received", LITERAL)
	pattern, _ = NewResourcePattern(TOPIC, "payments.received", PREFIXED)
	if filter.Matches(pattern) {
		t.Errorf("Expected filter to not match a pattern with different patternType: filter=%s, pattern=%s", filter, pattern)
	}

	// Test empty name matches any pattern
	filter, _ = NewResourcePatternFilter(TOPIC, "", LITERAL)
	pattern, _ = NewResourcePattern(TOPIC, "payments.received", LITERAL)
	if !filter.Matches(pattern) {
		t.Errorf("Expected filter with empty name to match any pattern: filter=%s, pattern=%s", filter, pattern)
	}

	// Test ANY pattern type matches any pattern
	filter, _ = NewResourcePatternFilter(TOPIC, "payments.received", ANY)
	pattern, _ = NewResourcePattern(TOPIC, "payments.received", LITERAL)
	if !filter.Matches(pattern) {
		t.Errorf("Expected filter with ANY patternType to match any pattern: filter=%s, pattern=%s", filter, pattern)
	}
}

func TestResourcePatternFilter_MatchesAtMostOne(t *testing.T) {
	filter, _ := NewResourcePatternFilter(TOPIC, "payments.received", LITERAL)
	if !filter.MatchesAtMostOne() {
		t.Errorf("ResourcePatternFilter should match at most one pattern")
	}

	filter, _ = NewResourcePatternFilter(ALL_RESOURCES, "payments.received", LITERAL)
	if filter.MatchesAtMostOne() {
		t.Errorf("ResourcePatternFilter with resourceType ALL_RESOURCES should not match at most one pattern")
	}

	filter, _ = NewResourcePatternFilter(TOPIC, "", LITERAL)
	if filter.MatchesAtMostOne() {
		t.Errorf("ResourcePatternFilter with empty name should not match at most one pattern")
	}

	filter, _ = NewResourcePatternFilter(TOPIC, "payments.received", MATCH)
	if filter.MatchesAtMostOne() {
		t.Errorf("ResourcePatternFilter with patternType MATCH should not match at most one pattern")
	}
}

func TestResourcePatternFilter_String(t *testing.T) {
	filter, _ := NewResourcePatternFilter(TOPIC, "payments.received", LITERAL)
	expected := `ResourcePatternFilter{resourceType=TOPIC, name="payments.received", patternType=LITERAL}`
	if filter.String() != expected {
		t.Errorf("ResourcePatternFilter.String() returned unexpected value: %s", filter.String())
	}
}
