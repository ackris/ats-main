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
	"bytes"
	"errors"
	"fmt"
	"hash/fnv"
	"strings"
)

// Bytes represents an immutable byte array.
type Bytes struct {
	data []byte
	hash uint32
}

// NewBytes creates a new Bytes instance from the provided byte slice.
// If data is nil, it returns nil.
//
// Example:
//
//	b := NewBytes([]byte{0x01, 0x02, 0x03})
//	fmt.Println(b) // Output: [0x01 0x02 0x03]
func NewBytes(data []byte) *Bytes {
	if data == nil {
		return nil
	}
	b := &Bytes{
		data: append([]byte(nil), data...),
	}
	return b
}

// Get returns a copy of the underlying byte slice. The returned slice is a
// copy to ensure immutability.
//
// Example:
//
//	b := NewBytes([]byte{0x01, 0x02, 0x03})
//	data := b.Get()
//	fmt.Println(data) // Output: [0x01 0x02 0x03]
func (b *Bytes) Get() []byte {
	return append([]byte(nil), b.data...)
}

// HashCode returns the hash code for the Bytes object using the FNV-1a hash algorithm.
// It returns the cached hash code if it has been computed previously.
//
// Example:
//
//	b := NewBytes([]byte{0x01, 0x02, 0x03})
//	fmt.Println(b.HashCode()) // Output: (a hash value, e.g., 28273073)
func (b *Bytes) HashCode() uint32 {
	if b.hash == 0 {
		h := fnv.New32a()
		h.Write(b.data)
		b.hash = h.Sum32()
	}
	return b.hash
}

// Equals checks if the current Bytes object is equal to another Bytes object.
// Two Bytes objects are considered equal if they have the same byte data and hash code.
//
// Example:
//
//	b1 := NewBytes([]byte{0x01, 0x02, 0x03})
//	b2 := NewBytes([]byte{0x01, 0x02, 0x03})
//	fmt.Println(b1.Equals(b2)) // Output: true
func (b *Bytes) Equals(other *Bytes) bool {
	if b == nil && other == nil {
		return true
	}
	if b == nil || other == nil {
		return false
	}
	if b.HashCode() != other.HashCode() {
		return false
	}
	return bytes.Equal(b.data, other.data)
}

// CompareTo compares the current Bytes object with another Bytes object.
// It returns a negative integer, zero, or a positive integer if the current
// Bytes object is less than, equal to, or greater than the other object.
//
// Example:
//
//	b1 := NewBytes([]byte{0x01, 0x02, 0x03})
//	b2 := NewBytes([]byte{0x01, 0x02, 0x04})
//	fmt.Println(b1.CompareTo(b2)) // Output: -1
func (b *Bytes) CompareTo(other *Bytes) int {
	if other == nil {
		return 1
	}
	return bytes.Compare(b.data, other.data)
}

// String returns a printable representation of the byte array. Non-printable
// characters are hex-escaped in the format \x%02X.
//
// Example:
//
//	b := NewBytes([]byte{0x01, 0x02, 0x03, 0xFF})
//	fmt.Println(b.String()) // Output: "\x01\x02\x03\xFF"
func (b *Bytes) String() string {
	return toString(b.data)
}

// Increment increments the byte array by 1 and returns a new Bytes object.
// If incrementing causes an overflow, it returns an error.
//
// Example:
//
//	b := NewBytes([]byte{0xFF, 0xFF})
//	newB, err := Increment(b)
//	if err != nil {
//		fmt.Println(err) // Output: byte array overflow
//	} else {
//		fmt.Println(newB) // Output: [0x00 0x00 0x01]
//	}
func Increment(input *Bytes) (*Bytes, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	data := input.Get()
	result := make([]byte, len(data))
	carry := 1

	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == 0xFF && carry == 1 {
			result[i] = 0x00
		} else {
			result[i] = data[i] + byte(carry)
			carry = 0
		}
	}

	if carry != 0 {
		return nil, errors.New("byte array overflow")
	}

	return NewBytes(result), nil
}

// toString creates a printable representation of a byte slice with hex escaping.
// It formats non-printable characters as \x%02X.
//
// Example:
//
//	data := []byte{0x01, 0x02, 0x03, 0xFF}
//	fmt.Println(toString(data)) // Output: "\x01\x02\x03\xFF"
func toString(data []byte) string {
	var sb strings.Builder
	for _, b := range data {
		if b >= ' ' && b <= '~' && b != '\\' {
			sb.WriteByte(b)
		} else {
			sb.WriteString(fmt.Sprintf("\\x%02X", b))
		}
	}
	return sb.String()
}

// ByteArrayComparator defines a comparison function for byte slices with offsets and lengths.
// It returns a negative integer, zero, or a positive integer if the first byte slice is less than,
// equal to, or greater than the second byte slice, respectively.
//
// Example:
//
//	comparator := ByteArrayComparator(LexicographicByteArrayComparator)
//	result := comparator([]byte{0x01, 0x02}, []byte{0x01, 0x03}, 0, 2, 0, 2)
//	fmt.Println(result) // Output: -1
type ByteArrayComparator func(a, b []byte, offsetA, lengthA, offsetB, lengthB int) int

// LexicographicByteArrayComparator compares two byte slices lexicographically with offsets and lengths.
// It returns a negative integer, zero, or a positive integer if the first byte slice is less than,
// equal to, or greater than the second byte slice, respectively.
//
// Example:
//
//	a := []byte{0x01, 0x02, 0x03}
//	b := []byte{0x01, 0x02, 0x04}
//	result := LexicographicByteArrayComparator(a, b, 0, 3, 0, 3)
//	fmt.Println(result) // Output: -1
func LexicographicByteArrayComparator(a, b []byte, offsetA, lengthA, offsetB, lengthB int) int {
	if offsetA < 0 || offsetB < 0 || lengthA < 0 || lengthB < 0 || offsetA+lengthA > len(a) || offsetB+lengthB > len(b) {
		return 0
	}

	endA := offsetA + lengthA
	endB := offsetB + lengthB

	for i, j := offsetA, offsetB; i < endA && j < endB; i, j = i+1, j+1 {
		if a[i] != b[j] {
			return int(a[i]) - int(b[j])
		}
	}
	return lengthA - lengthB
}
