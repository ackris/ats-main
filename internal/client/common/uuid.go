package common

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// Constants
const (
	// Base64 encoded UUID length without padding
	base64UUIDLength = 22
	// UUID size in bytes
	uuidSize = 16
)

// Uuid represents a 128-bit UUID
type Uuid struct {
	mostSignificantBits  int64
	leastSignificantBits int64
}

// Predefined reserved UUIDs
var (
	OneUUID         = NewUuid(0, 1)
	MetadataTopicID = OneUUID
	ZeroUUID        = NewUuid(0, 0)
	Reserved        = map[Uuid]struct{}{
		OneUUID:         {},
		MetadataTopicID: {},
		ZeroUUID:        {},
	}
)

// NewUuid creates a new Uuid from the provided 64-bit parts.
func NewUuid(mostSigBits, leastSigBits int64) Uuid {
	return Uuid{
		mostSignificantBits:  mostSigBits,
		leastSignificantBits: leastSigBits,
	}
}

// RandomUuid generates a random UUID and ensures it is not in the Reserved set and doesn't start with a dash in its Base64 representation.
func RandomUuid() Uuid {
	rand.Seed(time.Now().UnixNano())
	for {
		u := unsafeRandomUuid()
		if _, reserved := Reserved[u]; !reserved && !startsWithDash(u.String()) {
			return u
		}
	}
}

// unsafeRandomUuid generates a UUID using the google/uuid package.
func unsafeRandomUuid() Uuid {
	gUUID := uuid.New()
	b := gUUID[:]
	return NewUuid(
		int64(b[0])<<56|int64(b[1])<<48|int64(b[2])<<40|int64(b[3])<<32|
			int64(b[4])<<24|int64(b[5])<<16|int64(b[6])<<8|int64(b[7]),
		int64(b[8])<<56|int64(b[9])<<48|int64(b[10])<<40|int64(b[11])<<32|
			int64(b[12])<<24|int64(b[13])<<16|int64(b[14])<<8|int64(b[15]),
	)
}

// startsWithDash checks if the Base64 encoded UUID string starts with a dash.
func startsWithDash(s string) bool {
	return len(s) > 0 && s[0] == '-'
}

// String returns the Base64 encoded representation of the UUID.
func (u Uuid) String() string {
	b := u.getBytes()
	return base64.RawURLEncoding.EncodeToString(b)
}

// FromString creates a UUID from its Base64 encoded string representation.
func FromString(s string) (Uuid, error) {
	if len(s) > base64UUIDLength {
		return Uuid{}, errors.New("input string is too long to be decoded as a base64 UUID")
	}

	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return Uuid{}, fmt.Errorf("error decoding base64 string: %w", err)
	}

	if len(b) != uuidSize {
		return Uuid{}, errors.New("decoded byte slice is not 16 bytes long")
	}

	return Uuid{
		mostSignificantBits: int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 |
			int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7]),
		leastSignificantBits: int64(b[8])<<56 | int64(b[9])<<48 | int64(b[10])<<40 | int64(b[11])<<32 |
			int64(b[12])<<24 | int64(b[13])<<16 | int64(b[14])<<8 | int64(b[15]),
	}, nil
}

func (u Uuid) getBytes() []byte {
	b := make([]byte, uuidSize)
	copy(b[0:8], []byte{
		byte(u.mostSignificantBits >> 56),
		byte(u.mostSignificantBits >> 48),
		byte(u.mostSignificantBits >> 40),
		byte(u.mostSignificantBits >> 32),
		byte(u.mostSignificantBits >> 24),
		byte(u.mostSignificantBits >> 16),
		byte(u.mostSignificantBits >> 8),
		byte(u.mostSignificantBits),
	})
	copy(b[8:16], []byte{
		byte(u.leastSignificantBits >> 56),
		byte(u.leastSignificantBits >> 48),
		byte(u.leastSignificantBits >> 40),
		byte(u.leastSignificantBits >> 32),
		byte(u.leastSignificantBits >> 24),
		byte(u.leastSignificantBits >> 16),
		byte(u.leastSignificantBits >> 8),
		byte(u.leastSignificantBits),
	})
	return b
}

// Compare compares two UUIDs lexicographically.
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

// ToArray converts a slice of Uuid to an array of [2]int64.
func ToArray(slice []Uuid) [][2]int64 {
	arr := make([][2]int64, len(slice))
	for i, u := range slice {
		arr[i] = [2]int64{u.mostSignificantBits, u.leastSignificantBits}
	}
	return arr
}

// ToList converts an array of [2]int64 to a slice of Uuid.
func ToList(arr [][2]int64) []Uuid {
	slice := make([]Uuid, len(arr))
	for i, v := range arr {
		slice[i] = NewUuid(v[0], v[1])
	}
	return slice
}
