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

package utils

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// ExponentialBackoff provides an implementation of an exponential backoff algorithm with optional jitter.
//
// Exponential backoff is a strategy for handling retries in which the wait time between retries grows exponentially.
// This is particularly useful in distributed systems to avoid overwhelming a system under load.
// Jitter can be added to prevent all clients from retrying at the same time.
type ExponentialBackoff struct {
	initialInterval int64      // The initial wait time (in milliseconds) before the first retry.
	multiplier      int64      // The multiplier for each retry interval. Must be greater than 1.
	maxInterval     int64      // The maximum wait time (in milliseconds) between retries.
	jitter          float64    // The maximum jitter to apply to the backoff interval. Should be between 0.0 and 1.0.
	expMax          float64    // The maximum exponent value used to calculate the retry interval.
	rng             *rand.Rand // Random number generator used for applying jitter.
}

// NewExponentialBackoff creates a new ExponentialBackoff instance with the specified parameters.
//
// Parameters:
// - initialInterval: The initial wait time in milliseconds before the first retry. Must be greater than 0.
// - multiplier: The multiplier to apply for each retry interval. Must be greater than 1.
// - maxInterval: The maximum wait time in milliseconds between retries. Must be greater than or equal to initialInterval.
// - jitter: The maximum jitter to apply to the backoff interval. Must be non-negative.
//
// Returns:
// - A pointer to an ExponentialBackoff instance and an error if any of the parameters are invalid.
//
// Example:
//
//	eb, err := NewExponentialBackoff(1000, 2, 32000, 0.5)
//	if err != nil {
//	    log.Fatalf("Error creating ExponentialBackoff: %v", err)
//	}
func NewExponentialBackoff(initialInterval, multiplier, maxInterval int64, jitter float64) (*ExponentialBackoff, error) {
	if initialInterval <= 0 {
		return nil, errors.New("initialInterval must be greater than 0")
	}
	if multiplier <= 1 {
		return nil, errors.New("multiplier must be greater than 1")
	}
	if maxInterval < initialInterval {
		return nil, errors.New("maxInterval must be greater than or equal to initialInterval")
	}
	if jitter < 0 {
		return nil, errors.New("jitter must be non-negative")
	}

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	expMax := float64(0)
	if maxInterval > initialInterval {
		expMax = math.Log(float64(maxInterval)/float64(initialInterval)) / math.Log(float64(multiplier))
	}

	return &ExponentialBackoff{
		initialInterval: initialInterval,
		multiplier:      multiplier,
		maxInterval:     maxInterval,
		jitter:          jitter,
		expMax:          expMax,
		rng:             rng,
	}, nil
}

// InitialInterval returns the initial interval used in the backoff algorithm.
//
// Returns:
// - The initial wait time in milliseconds before the first retry.
//
// Example:
//
//	interval := eb.InitialInterval()
//	fmt.Printf("Initial interval: %d milliseconds", interval)
func (eb *ExponentialBackoff) InitialInterval() int64 {
	return eb.initialInterval
}

// Backoff calculates the retry interval based on the number of attempts.
//
// Parameters:
// - attempts: The number of retry attempts made. This value is used to determine the backoff interval.
//
// Returns:
// - The calculated backoff interval in milliseconds, which will be a value between the initial interval and the maximum interval.
//
// Example:
//
//	interval := eb.Backoff(3)
//	fmt.Printf("Backoff interval for attempt 3: %d milliseconds", interval)
func (eb *ExponentialBackoff) Backoff(attempts int64) int64 {
	if eb.expMax == 0 {
		return eb.initialInterval
	}

	exp := math.Min(float64(attempts), eb.expMax)
	term := float64(eb.initialInterval) * math.Pow(float64(eb.multiplier), exp)

	randomFactor := 1.0
	if eb.jitter > 0 {
		randomFactor = 1.0 + (2.0*eb.rng.Float64()-1.0)*eb.jitter
	}
	if randomFactor < 1.0 {
		randomFactor = 1.0
	}

	backoffValue := int64(randomFactor * term)
	if backoffValue > eb.maxInterval {
		return eb.maxInterval
	}
	return backoffValue
}

// String provides a string representation of the ExponentialBackoff instance.
//
// Returns:
// - A string describing the ExponentialBackoff instance with its parameters.
//
// Example:
//
//	fmt.Println(eb.String())
//	// Output: ExponentialBackoff{multiplier=2, expMax=5.000000, initialInterval=1000, jitter=0.500000}
func (eb *ExponentialBackoff) String() string {
	return fmt.Sprintf("ExponentialBackoff{multiplier=%d, expMax=%f, initialInterval=%d, jitter=%f}",
		eb.multiplier, eb.expMax, eb.initialInterval, eb.jitter)
}
