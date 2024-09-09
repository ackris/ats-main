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

import "testing"

func TestConsumerGroupStateString(t *testing.T) {
	tests := []struct {
		state    ConsumerGroupState
		expected string
	}{
		{Unknown, "Unknown"},
		{PreparingRebalance, "PreparingRebalance"},
		{CompletingRebalance, "CompletingRebalance"},
		{Stable, "Stable"},
		{Dead, "Dead"},
		{Empty, "Empty"},
		{Assigning, "Assigning"},
		{Reconciling, "Reconciling"},
		{ConsumerGroupState(100), "InvalidState"}, // Testing an invalid state
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := test.state.String()
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestParseConsumerGroupState(t *testing.T) {
	tests := []struct {
		input    string
		expected ConsumerGroupState
		wantErr  bool
	}{
		{"UNKNOWN", Unknown, false},
		{"PREPARINGREBALANCE", PreparingRebalance, false},
		{"COMPLETINGREBALANCE", CompletingRebalance, false},
		{"STABLE", Stable, false},
		{"DEAD", Dead, false},
		{"EMPTY", Empty, false},
		{"ASSIGNING", Assigning, false},
		{"RECONCILING", Reconciling, false},
		{"invalid", Unknown, true},
		{"", Unknown, true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ParseConsumerGroupState(test.input)
			if result != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
			if (err != nil) != test.wantErr {
				t.Errorf("expected error: %v, got %v", test.wantErr, err)
			}
		})
	}
}
