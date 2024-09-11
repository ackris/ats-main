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

import "hash"

// UpdateChecksum updates the given checksum with the contents of the buffer slice.
// It handles buffer slices efficiently without modifying the buffer's position or limit.
//
// Parameters:
//
//	checksum: The hash.Hash instance to be updated. Common implementations include crc32 and crc64.
//	buffer: The byte slice containing data to be used for checksum calculation.
//	offset: The start position within the buffer slice to begin updating the checksum.
//	length: The number of bytes to be used from the buffer slice.
//
// If offset or length is out of bounds, this function does nothing.
func UpdateChecksum(checksum hash.Hash, buffer []byte, offset, length int) {
	if offset < 0 || length < 0 || offset+length > len(buffer) {
		// If parameters are out of bounds, do nothing
		return
	}

	// Update the checksum with the specified slice of the buffer
	checksum.Write(buffer[offset : offset+length])
}

// UpdateInt updates the given checksum with the 4-byte representation of an integer value.
// The integer is split into four bytes in big-endian order and written to the checksum.
//
// Parameters:
//
//	checksum: The hash.Hash instance to be updated. Common implementations include crc32 and crc64.
//	value: The integer value to be converted to bytes and used for checksum calculation.
//
// Example:
//
//	crc32Checksum := crc32.New(crc32.IEEETable)
//	UpdateInt(crc32Checksum, 1234)
//	fmt.Println(crc32Checksum.Sum32()) // prints the CRC32 checksum of the integer value 1234
func UpdateInt(checksum hash.Hash, value int) {
	// Use a fixed-size byte array for efficient memory allocation
	var bytes [4]byte
	bytes[0] = byte(value >> 24)
	bytes[1] = byte(value >> 16)
	bytes[2] = byte(value >> 8)
	bytes[3] = byte(value)
	checksum.Write(bytes[:])
}

// UpdateLong updates the given checksum with the 8-byte representation of a 64-bit integer (long) value.
// The long is split into eight bytes in big-endian order and written to the checksum.
//
// Parameters:
//
//	checksum: The hash.Hash instance to be updated. Common implementations include crc32 and crc64.
//	value: The long value to be converted to bytes and used for checksum calculation.
//
// Example:
//
//	crc64Checksum := crc64.New(crc64.MakeTable(crc64.ECMA))
//	UpdateLong(crc64Checksum, 1234567890123456789)
//	fmt.Println(crc64Checksum.Sum64()) // prints the CRC64 checksum of the long value 1234567890123456789
func UpdateLong(checksum hash.Hash, value int64) {
	// Use a fixed-size byte array for efficient memory allocation
	var bytes [8]byte
	bytes[0] = byte(value >> 56)
	bytes[1] = byte(value >> 48)
	bytes[2] = byte(value >> 40)
	bytes[3] = byte(value >> 32)
	bytes[4] = byte(value >> 24)
	bytes[5] = byte(value >> 16)
	bytes[6] = byte(value >> 8)
	bytes[7] = byte(value)
	checksum.Write(bytes[:])
}
