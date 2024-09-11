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
	"encoding/binary"
	"errors"
	"io"
	"math"
)

// ReadUnsignedInt reads an unsigned 32-bit integer from the byte slice at the specified index.
// The integer is read in big-endian format.
//
// Example usage:
//
//	data := []byte{0x00, 0x00, 0x01, 0x2C}
//	value := ReadUnsignedInt(data, 0)
//	fmt.Println(value) // Output: 300
//
// Parameters:
//
//	data - The byte slice containing the integer data.
//	index - The starting index in the byte slice where the integer is to be read.
//
// Returns:
//
//	uint32 - The unsigned 32-bit integer read from the byte slice.
func ReadUnsignedInt(data []byte, index int) uint32 {
	if index+4 > len(data) {
		return 0 // or handle out-of-bounds error
	}
	return binary.BigEndian.Uint32(data[index:])
}

// ReadUnsignedIntLE reads an unsigned 32-bit integer in little-endian format from the byte slice at the specified offset.
//
// Example usage:
//
//	data := []byte{0x2C, 0x01, 0x00, 0x00}
//	value := ReadUnsignedIntLE(data, 0)
//	fmt.Println(value) // Output: 300
//
// Parameters:
//
//	data - The byte slice containing the integer data.
//	offset - The starting offset in the byte slice where the integer is to be read.
//
// Returns:
//
//	uint32 - The unsigned 32-bit integer read from the byte slice.
func ReadUnsignedIntLE(data []byte, offset int) uint32 {
	if offset+4 > len(data) {
		return 0 // or handle out-of-bounds error
	}
	return binary.LittleEndian.Uint32(data[offset:])
}

// ReadIntBE reads a big-endian 32-bit signed integer from the byte slice at the specified offset.
//
// Example usage:
//
//	data := []byte{0x00, 0x00, 0x01, 0x2C}
//	value := ReadIntBE(data, 0)
//	fmt.Println(value) // Output: 300
//
// Parameters:
//
//	data - The byte slice containing the integer data.
//	offset - The starting offset in the byte slice where the integer is to be read.
//
// Returns:
//
//	int32 - The signed 32-bit integer read from the byte slice.
func ReadIntBE(data []byte, offset int) int32 {
	if offset+4 > len(data) {
		return 0 // or handle out-of-bounds error
	}
	return int32(binary.BigEndian.Uint32(data[offset:]))
}

// WriteUnsignedInt writes an unsigned 32-bit integer to the byte slice at the specified index in big-endian format.
//
// Example usage:
//
//	data := make([]byte, 4)
//	WriteUnsignedInt(data, 0, 300)
//	fmt.Println(data) // Output: [0 0 1 44]
//
// Parameters:
//
//	data - The byte slice to which the integer will be written.
//	index - The starting index in the byte slice where the integer will be written.
//	value - The unsigned 32-bit integer to write.
func WriteUnsignedInt(data []byte, index int, value uint32) {
	if index+4 <= len(data) {
		binary.BigEndian.PutUint32(data[index:], value)
	}
}

// WriteUnsignedIntLE writes an unsigned 32-bit integer to the byte slice at the specified offset in little-endian format.
//
// Example usage:
//
//	data := make([]byte, 4)
//	WriteUnsignedIntLE(data, 0, 300)
//	fmt.Println(data) // Output: [44 1 0 0]
//
// Parameters:
//
//	data - The byte slice to which the integer will be written.
//	offset - The starting offset in the byte slice where the integer will be written.
//	value - The unsigned 32-bit integer to write.
func WriteUnsignedIntLE(data []byte, offset int, value uint32) {
	if offset+4 <= len(data) {
		binary.LittleEndian.PutUint32(data[offset:], value)
	}
}

// WriteVarint writes a 32-bit signed integer using variable-length encoding to the buffer.
//
// Example usage:
//
//	buf := &bytes.Buffer{}
//	WriteVarint(300, buf)
//	fmt.Println(buf.Bytes()) // Output: [172 2]
//
// Parameters:
//
//	value - The 32-bit signed integer to encode.
//	buf - The writer where the encoded value will be written.
func WriteVarint(value uint32, buf io.Writer) {
	for value >= 0x80 {
		buf.Write([]byte{byte(value | 0x80)})
		value >>= 7
	}
	buf.Write([]byte{byte(value)})
}

// WriteVarlong writes a 64-bit signed integer using variable-length encoding to the buffer.
//
// Example usage:
//
//	buf := &bytes.Buffer{}
//	WriteVarlong(300, buf)
//	fmt.Println(buf.Bytes()) // Output: [172 2]
//
// Parameters:
//
//	value - The 64-bit signed integer to encode.
//	buf - The writer where the encoded value will be written.
func WriteVarlong(value int64, buf io.Writer) {
	for (value &^ 0x7F) != 0 {
		buf.Write([]byte{byte(value&0x7F | 0x80)})
		value >>= 7
	}
	buf.Write([]byte{byte(value)})
}

// SizeOfUnsignedVarint calculates the number of bytes required to encode an unsigned 32-bit integer.
//
// Example usage:
//
//	size := SizeOfUnsignedVarint(300)
//	fmt.Println(size) // Output: 2
//
// Parameters:
//
//	value - The unsigned 32-bit integer whose size is to be calculated.
//
// Returns:
//
//	int - The number of bytes required to encode the value.
func SizeOfUnsignedVarint(value uint32) int {
	if value < 1<<7 {
		return 1
	} else if value < 1<<14 {
		return 2
	} else if value < 1<<21 {
		return 3
	} else if value < 1<<28 {
		return 4
	}
	return 5
}

// SizeOfVarint calculates the number of bytes required to encode a 32-bit signed integer.
//
// Example usage:
//
//	size := SizeOfVarint(300)
//	fmt.Println(size) // Output: 2
//
// Parameters:
//
//	value - The 32-bit signed integer whose size is to be calculated.
//
// Returns:
//
//	int - The number of bytes required to encode the value.
func SizeOfVarint(value int32) int {
	if value < 0 {
		value = (value << 1) ^ (value >> 31)
	}
	return SizeOfUnsignedVarint(uint32(value))
}

// SizeOfVarlong calculates the number of bytes required to encode a 64-bit signed integer.
//
// Example usage:
//
//	size := SizeOfVarlong(300)
//	fmt.Println(size) // Output: 2
//
// Parameters:
//
//	value - The 64-bit signed integer whose size is to be calculated.
//
// Returns:
//
//	int - The number of bytes required to encode the value.
func SizeOfVarlong(value int64) int {
	// Handle zero separately
	if value == 0 {
		return 1
	}

	// Count the number of bytes needed
	size := 0
	for value != 0 {
		size++
		value >>= 7
	}

	return size
}

// SizeOfUnsignedVarlong calculates the number of bytes required to encode an unsigned 64-bit integer.
//
// Example usage:
//
//	size := SizeOfUnsignedVarlong(300)
//	fmt.Println(size) // Output: 2
//
// Parameters:
//
//	value - The unsigned 64-bit integer whose size is to be calculated.
//
// Returns:
//
//	int - The number of bytes required to encode the value.
func SizeOfUnsignedVarlong(value uint64) int {
	if value < 1<<7 {
		return 1
	} else if value < 1<<14 {
		return 2
	} else if value < 1<<21 {
		return 3
	} else if value < 1<<28 {
		return 4
	} else if value < 1<<35 {
		return 5
	} else if value < 1<<42 {
		return 6
	} else if value < 1<<49 {
		return 7
	} else if value < 1<<56 {
		return 8
	}
	return 9
}

// ReadUnsignedVarint reads an unsigned integer from the reader in variable-length encoding.
//
// Example usage:
//
//	buf := bytes.NewReader([]byte{172, 2})
//	value, err := ReadUnsignedVarint(buf)
//	fmt.Println(value, err) // Output: 300 <nil>
//
// Parameters:
//
//	r - The reader from which the encoded value will be read.
//
// Returns:
//
//	uint32 - The decoded unsigned integer.
//	error - An error if any occurred while reading.
func ReadUnsignedVarint(r io.Reader) (uint32, error) {
	var result uint32
	var shift uint
	for i := 0; i < 5; i++ {
		b, err := readByte(r)
		if err != nil {
			return 0, err
		}
		result |= (uint32(b) & 0x7F) << shift
		if b&0x80 == 0 {
			return result, nil
		}
		shift += 7
	}
	return 0, errors.New("varint is too long")
}

// ReadVarint reads a signed 32-bit integer from the reader using variable-length encoding.
//
// Example usage:
//
//	buf := bytes.NewReader([]byte{172, 2})
//	value, err := ReadVarint(buf)
//	fmt.Println(value, err) // Output: 300 <nil>
//
// Parameters:
//
//	buf - The buffer from which the encoded value will be read.
//
// Returns:
//
//	int32 - The decoded signed integer.
//	error - An error if any occurred while reading.
func ReadVarint(buf *bytes.Buffer) (int32, error) {
	var value int32
	var shift uint
	for {
		b, err := buf.ReadByte()
		if err != nil {
			if err == io.EOF {
				return 0, errors.New("unexpected end of buffer")
			}
			return 0, err
		}
		value |= int32(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
		if shift >= 32 {
			return 0, errors.New("varint too big")
		}
	}
	return value, nil
}

// ReadVarlong reads a signed 64-bit integer from the reader using variable-length encoding.
//
// Example usage:
//
//	buf := bytes.NewReader([]byte{172, 2})
//	value, err := ReadVarlong(buf)
//	fmt.Println(value, err) // Output: 300 <nil>
//
// Parameters:
//
//	buf - The buffer from which the encoded value will be read.
//
// Returns:
//
//	int64 - The decoded signed integer.
//	error - An error if any occurred while reading.
func ReadVarlong(buf *bytes.Buffer) (int64, error) {
	var value int64
	var shift uint
	for {
		b, err := buf.ReadByte()
		if err != nil {
			if err == io.EOF {
				return 0, errors.New("unexpected end of buffer")
			}
			return 0, err
		}
		value |= int64(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
		if shift >= 64 {
			return 0, errors.New("varlong too big")
		}
	}
	return value, nil
}

// ReadDouble reads a 64-bit IEEE 754 floating-point number from the reader.
//
// Example usage:
//
//	buf := bytes.NewReader([]byte{0, 0, 0, 0, 0, 0, 248, 63})
//	value, err := ReadDouble(buf)
//	fmt.Println(value, err) // Output: 1 <nil>
//
// Parameters:
//
//	r - The reader from which the encoded value will be read.
//
// Returns:
//
//	float64 - The decoded floating-point number.
//	error - An error if any occurred while reading.
func ReadDouble(r io.Reader) (float64, error) {
	var buf [8]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(binary.BigEndian.Uint64(buf[:])), nil
}

// WriteUnsignedVarint writes an unsigned 32-bit integer to the buffer using variable-length encoding.
//
// Example usage:
//
//	buf := &bytes.Buffer{}
//	WriteUnsignedVarint(300, buf)
//	fmt.Println(buf.Bytes()) // Output: [172 2]
//
// Parameters:
//
//	value - The unsigned 32-bit integer to encode.
//	buf - The buffer to which the encoded data will be written.
func WriteUnsignedVarint(value uint32, buf *bytes.Buffer) {
	for {
		if value < 0x80 {
			buf.WriteByte(byte(value))
			return
		}
		buf.WriteByte(byte((value & 0x7F) | 0x80))
		value >>= 7
	}
}

// WriteDouble writes a 64-bit IEEE 754 floating-point number to the buffer.
//
// Example usage:
//
//	buf := &bytes.Buffer{}
//	WriteDouble(1.0, buf)
//	fmt.Println(buf.Bytes()) // Output: [0 0 0 0 0 0 248 63]
//
// Parameters:
//
//	value - The floating-point number to encode.
//	buf - The writer to which the encoded data will be written.
func WriteDouble(value float64, buf io.Writer) {
	var bits = math.Float64bits(value)
	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, bits)
	buf.Write(b)
}

// readByte reads a single byte from the reader.
//
// Example usage:
//
//	r := bytes.NewReader([]byte{42})
//	b, err := readByte(r)
//	fmt.Println(b, err) // Output: 42 <nil>
//
// Parameters:
//
//	r - The reader from which the byte will be read.
//
// Returns:
//
//	byte - The byte read from the reader.
//	error - An error if any occurred while reading.
func readByte(r io.Reader) (byte, error) {
	var b [1]byte
	_, err := r.Read(b[:])
	return b[0], err
}
