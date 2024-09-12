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
	"hash/crc32"
	"hash/crc64"
	"testing"
)

// TestUpdateChecksum tests the UpdateChecksum function.
func TestUpdateChecksum(t *testing.T) {
	buffer := []byte("test data")

	// CRC32 test
	crc32Checksum := crc32.New(crc32.IEEETable)
	expectedCRC32 := crc32.Checksum(buffer, crc32.IEEETable)
	UpdateChecksum(crc32Checksum, buffer, 0, len(buffer))
	if got := crc32Checksum.Sum32(); got != expectedCRC32 {
		t.Errorf("UpdateChecksum for CRC32 = %v; want %v", got, expectedCRC32)
	}

	// CRC64 test
	crc64Checksum := crc64.New(crc64.MakeTable(crc64.ECMA))
	expectedCRC64 := crc64.Update(0, crc64.MakeTable(crc64.ECMA), buffer)
	UpdateChecksum(crc64Checksum, buffer, 0, len(buffer))
	if got := crc64Checksum.Sum64(); got != expectedCRC64 {
		t.Errorf("UpdateChecksum for CRC64 = %v; want %v", got, expectedCRC64)
	}
}

// TestUpdateInt tests the UpdateInt function.
func TestUpdateInt(t *testing.T) {
	// CRC32 test
	crc32Checksum := crc32.New(crc32.IEEETable)
	for _, input := range []int{0, 1, -1, 123456789, -123456789} {
		crc32Checksum.Reset()
		UpdateInt(crc32Checksum, input)
		expectedCRC32 := crc32.Checksum(intToBytes(input), crc32.IEEETable)
		if got := crc32Checksum.Sum32(); got != expectedCRC32 {
			t.Errorf("UpdateInt for CRC32 with input %d = %v; want %v", input, got, expectedCRC32)
		}
	}
}

// TestUpdateLong tests the UpdateLong function.
func TestUpdateLong(t *testing.T) {
	// CRC32 test
	crc32Checksum := crc32.New(crc32.IEEETable)
	for _, input := range []int64{0, 1, -1, 1234567890123456789, -1234567890123456789} {
		crc32Checksum.Reset()
		UpdateLong(crc32Checksum, input)
		expectedCRC32 := crc32.Checksum(longToBytes(input), crc32.IEEETable)
		if got := crc32Checksum.Sum32(); got != expectedCRC32 {
			t.Errorf("UpdateLong for CRC32 with input %d = %v; want %v", input, got, expectedCRC32)
		}
	}
}

// TestCompute tests the Compute function for CRC32C.
func TestCompute(t *testing.T) {
	data := []byte("Hello, CRC32C!")
	crc32c := NewCRC32C()
	expectedChecksum := crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))

	checksum, err := crc32c.Compute(data)
	if err != nil {
		t.Fatalf("Compute failed: %v", err)
	}
	if checksum != expectedChecksum {
		t.Errorf("Compute CRC32C = %v; want %v", checksum, expectedChecksum)
	}
}

// TestComputeWithOffset tests the ComputeWithOffset function for CRC32C.
func TestComputeWithOffset(t *testing.T) {
	data := []byte("Hello, CRC32C!")
	crc32c := NewCRC32C()
	expectedChecksum := crc32.Checksum(data[7:13], crc32.MakeTable(crc32.Castagnoli)) // "CRC32C"

	checksum, err := crc32c.ComputeWithOffset(data, 7, 6) // "CRC32"
	if err != nil {
		t.Fatalf("ComputeWithOffset failed: %v", err)
	}
	if checksum != expectedChecksum {
		t.Errorf("ComputeWithOffset CRC32C = %v; want %v", checksum, expectedChecksum)
	}

	// Test with invalid offset
	_, err = crc32c.ComputeWithOffset(data, -1, 6)
	if err == nil {
		t.Error("Expected error for negative offset, got nil")
	}

	// Test with out of bounds length
	_, err = crc32c.ComputeWithOffset(data, 7, 10)
	if err == nil {
		t.Error("Expected error for out of bounds length, got nil")
	}
}

// TestReset tests the Reset method of the CRC32C struct.
func TestReset(t *testing.T) {
	crc32c := NewCRC32C()
	data := []byte("Hello, CRC32C!")

	// Compute initial checksum
	checksum1, err := crc32c.Compute(data)
	if err != nil {
		t.Fatalf("Compute failed: %v", err)
	}

	// Reset the CRC32C instance
	crc32c.Reset()

	// Compute checksum again after reset
	checksum2, err := crc32c.Compute(data)
	if err != nil {
		t.Fatalf("Compute after reset failed: %v", err)
	}

	// The two checksums should be the same
	if checksum1 != checksum2 {
		t.Errorf("Checksum after reset = %v; want %v", checksum2, checksum1)
	}
}

// intToBytes converts an int to a 4-byte slice.
func intToBytes(value int) []byte {
	var bytes [4]byte
	bytes[0] = byte(value >> 24)
	bytes[1] = byte(value >> 16)
	bytes[2] = byte(value >> 8)
	bytes[3] = byte(value)
	return bytes[:]
}

// longToBytes converts an int64 to an 8-byte slice.
func longToBytes(value int64) []byte {
	var bytes [8]byte
	bytes[0] = byte(value >> 56)
	bytes[1] = byte(value >> 48)
	bytes[2] = byte(value >> 40)
	bytes[3] = byte(value >> 32)
	bytes[4] = byte(value >> 24)
	bytes[5] = byte(value >> 16)
	bytes[6] = byte(value >> 8)
	bytes[7] = byte(value)
	return bytes[:]
}
