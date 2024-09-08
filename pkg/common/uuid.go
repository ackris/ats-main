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

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"

	"time"
)

// Constants
const (
	// UUIDSize is the size of a UUID in bytes.
	uuidSize = 16
)

// Uuid represents a 128-bit UUID.
type Uuid struct {
	mostSignificantBits  uint64 // Most significant 64 bits of the UUID
	leastSignificantBits uint64 // Least significant 64 bits of the UUID
}

// Predefined reserved UUIDs
var (
	// OneUUID is a predefined UUID with most significant bits 0 and least significant bits 1.
	OneUUID = NewUuid(0, 1)
	// MetadataTopicID is an alias for OneUUID.
	MetadataTopicID = OneUUID
	// ZeroUUID is a predefined UUID with both most and least significant bits set to 0.
	ZeroUUID = NewUuid(0, 0)
	// Reserved contains a set of predefined reserved UUIDs.
	Reserved = map[Uuid]struct{}{
		OneUUID:         {},
		MetadataTopicID: {},
		ZeroUUID:        {},
	}
)

// NewUuid creates a new Uuid from the provided 64-bit parts.
//
// Parameters:
//   - mostSigBits: The most significant 64 bits of the UUID.
//   - leastSigBits: The least significant 64 bits of the UUID.
//
// Returns:
//   - A new Uuid instance representing the combined 128-bit UUID.
//
// Example:
//
//	u := NewUuid(0x123456789abcdef0, 0xfedcba9876543210)
//	fmt.Println(u.String()) // Output: Base64 encoded UUID
func NewUuid(mostSigBits, leastSigBits uint64) Uuid {
	return Uuid{
		mostSignificantBits:  mostSigBits,
		leastSignificantBits: leastSigBits,
	}
}

// RandomUuid generates a random UUID and ensures it is not in the Reserved set
// and doesn't start with a dash in its Base64 representation.
//
// Returns:
//   - A randomly generated Uuid that is not reserved and does not start with a dash.
//
// Example:
//
//	randomUUID := RandomUuid()
//	fmt.Println(randomUUID.String()) // Output: Random Base64 encoded UUID
func RandomUuid() Uuid {
	// Create a new random source seeded with the current time
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	for {
		u := unsafeRandomUuid(random)
		if _, reserved := Reserved[u]; !reserved && !startsWithDash(u.String()) {
			return u
		}
	}
}

// unsafeRandomUuid generates a UUID using the google/uuid package with a custom random source.
// It does not check for reserved UUIDs or dash prefixes.
func unsafeRandomUuid(r *rand.Rand) Uuid {
	// Generate 16 random bytes
	b := make([]byte, uuidSize)
	r.Read(b)

	// Construct UUID from bytes
	return NewUuid(
		binary.BigEndian.Uint64(b[:8]),
		binary.BigEndian.Uint64(b[8:]),
	)
}

// startsWithDash checks if the Base64 encoded UUID string starts with a dash.
//
// Parameters:
//   - s: The Base64 encoded UUID string to check.
//
// Returns:
//   - true if the string starts with a dash, false otherwise.
//
// Example:
//
//	fmt.Println(startsWithDash("-abc")) // Output: true
//	fmt.Println(startsWithDash("abc"))   // Output: false
func startsWithDash(s string) bool {
	return len(s) > 0 && s[0] == '-'
}

// String returns the Base64 encoded representation of the UUID.
//
// Returns:
//   - A string representing the UUID in Base64 format.
//
// Example:
//
//	u := NewUuid(0x123456789abcdef0, 0xfedcba9876543210)
//	fmt.Println(u.String()) // Output: Base64 encoded UUID
func (u Uuid) String() string {
	b := u.getBytes()
	return base64.RawURLEncoding.EncodeToString(b)
}

// FromString creates a UUID from its Base64 encoded string representation.
//
// Parameters:
//   - s: The Base64 encoded string representation of the UUID.
//
// Returns:
//   - A Uuid instance if the string is valid, or an error if the string is invalid.
//
// Example:
//
//	u, err := FromString("0Aw-Aw0AQNaAQNaAQNaAQNaAQ")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//	    fmt.Println(u) // Output: {0 1}
//	}
func FromString(s string) (Uuid, error) {
	if len(s) == 0 {
		return Uuid{}, errors.New("input string is empty")
	}

	// Decode Base64 string
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return Uuid{}, fmt.Errorf("error decoding base64 string: %w", err)
	}

	// Validate the length of the decoded byte slice
	if len(b) != uuidSize {
		return Uuid{}, errors.New("decoded byte slice is not 16 bytes long")
	}

	// Create UUID from bytes
	return Uuid{
		mostSignificantBits:  binary.BigEndian.Uint64(b[:8]),
		leastSignificantBits: binary.BigEndian.Uint64(b[8:]),
	}, nil
}

// getBytes converts the UUID to a byte slice representation.
//
// Returns:
//   - A byte slice representing the UUID.
//
// Example:
//
//	u := NewUuid(0x123456789abcdef0, 0xfedcba9876543210)
//	bytes := u.getBytes()
//	fmt.Println(bytes) // Output: [byte representation of the UUID]
func (u Uuid) getBytes() []byte {
	b := make([]byte, uuidSize)
	binary.BigEndian.PutUint64(b[:8], uint64(u.mostSignificantBits))
	binary.BigEndian.PutUint64(b[8:], uint64(u.leastSignificantBits))
	return b
}

// Compare compares two UUIDs lexicographically.
//
// Parameters:
//   - other: The other Uuid to compare against.
//
// Returns:
//   - 1 if the current UUID is greater than the other,
//   - -1 if the current UUID is less than the other,
//   - 0 if they are equal.
//
// Example:
//
//	u1 := NewUuid(0x123456789abcdef0, 0xfedcba9876543210)
//	u2 := NewUuid(0x123456789abcdef0, 0xfedcba9876543211)
//	result := u1.Compare(u2)
//	fmt.Println(result) // Output: -1 (u1 is less than u2)
func (u Uuid) Compare(other Uuid) int {
	if u.mostSignificantBits > other.mostSignificantBits {
		return 1
	}
	if u.mostSignificantBits < other.mostSignificantBits {
		return -1
	}
	if u.leastSignificantBits > other.leastSignificantBits {
		return 1
	}
	if u.leastSignificantBits < other.leastSignificantBits {
		return -1
	}
	return 0
}

// ToArray converts a slice of Uuid to an array of [2]uint64.
//
// Parameters:
//   - slice: A slice of Uuid instances to convert.
//
// Returns:
//   - An array of [2]uint64 where each element contains the most and least significant bits of the UUID.
//
// Example:
//
//	uuids := []Uuid{NewUuid(0, 1), NewUuid(0, 2)}
//	arr := ToArray(uuids)
//	fmt.Println(arr) // Output: [[0 1] [0 2]]
func ToArray(slice []Uuid) [][2]uint64 {
	arr := make([][2]uint64, len(slice))
	for i, u := range slice {
		arr[i] = [2]uint64{u.mostSignificantBits, u.leastSignificantBits}
	}
	return arr
}

// ToList converts an array of [2]uint64 to a slice of Uuid.
//
// Parameters:
//   - arr: An array of [2]uint64 to convert.
//
// Returns:
//   - A slice of Uuid instances created from the provided array.
//
// Example:
//
//	arr := [][2]uint64{{0, 1}, {0, 2}}
//	uuids := ToList(arr)
//	fmt.Println(uuids[0].String()) // Output: Base64 encoded UUID of {0 1}
func ToList(arr [][2]uint64) []Uuid {
	slice := make([]Uuid, len(arr))
	for i, v := range arr {
		slice[i] = NewUuid(v[0], v[1])
	}
	return slice
}
