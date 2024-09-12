// Copyright 2024 Atomstate Technologies Private Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package utils

import (
	"testing"
)

func TestNewExponentialBackoff_Valid(t *testing.T) {
	eb, err := NewExponentialBackoff(1000, 2, 32000, 0.5)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if eb == nil {
		t.Fatal("Expected non-nil ExponentialBackoff instance")
	}
}

func TestNewExponentialBackoff_InvalidParameters(t *testing.T) {
	tests := []struct {
		initialInterval, multiplier, maxInterval int64
		jitter                                   float64
		expectedError                            string
	}{
		{-1, 2, 10000, 0.5, "initialInterval must be greater than 0"},
		{1000, 1, 10000, 0.5, "multiplier must be greater than 1"},
		{1000, 2, 500, 0.5, "maxInterval must be greater than or equal to initialInterval"},
		{1000, 2, 10000, -0.5, "jitter must be non-negative"},
	}

	for _, tt := range tests {
		_, err := NewExponentialBackoff(tt.initialInterval, tt.multiplier, tt.maxInterval, tt.jitter)
		if err == nil || err.Error() != tt.expectedError {
			t.Errorf("Expected error %v, got %v", tt.expectedError, err)
		}
	}
}

func TestBackoff_Calculation(t *testing.T) {
	eb, _ := NewExponentialBackoff(1000, 2, 32000, 0.5)

	tests := []struct {
		attempts int64
		min      int64
		max      int64
	}{
		{0, 1000, 1000},
		{1, 2000, 32000},
		{2, 4000, 32000},
		{5, 32000, 32000},
	}

	for _, tt := range tests {
		result := eb.Backoff(tt.attempts)
		if result < tt.min || result > tt.max {
			t.Errorf("Backoff(%d) = %d; want between %d and %d", tt.attempts, result, tt.min, tt.max)
		}
	}
}

func TestJitter_Application(t *testing.T) {
	eb, _ := NewExponentialBackoff(1000, 2, 32000, 0.5)

	for i := 0; i < 100; i++ {
		result := eb.Backoff(1)
		if result < 1000 || result > 32000 {
			t.Errorf("Backoff with jitter applied resulted in %d; want between 1000 and 32000", result)
		}
	}
}

func TestToString(t *testing.T) {
	eb, _ := NewExponentialBackoff(1000, 2, 32000, 0.5)
	expected := "ExponentialBackoff{multiplier=2, expMax=5.000000, initialInterval=1000, jitter=0.500000}"

	if got := eb.String(); got != expected {
		t.Errorf("String() = %v; want %v", got, expected)
	}
}
