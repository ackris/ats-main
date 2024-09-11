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
	"math"
	"testing"
)

func TestWriteVarint(t *testing.T) {
	buf := new(bytes.Buffer)
	WriteVarint(300, buf)
	expected := []byte{0xAC, 0x02}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Errorf("WriteVarint() = %v; want %v", buf.Bytes(), expected)
	}
}

func TestWriteVarlong(t *testing.T) {
	buf := new(bytes.Buffer)
	WriteVarlong(300, buf)
	expected := []byte{0xAC, 0x02}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Errorf("WriteVarlong() = %v; want %v", buf.Bytes(), expected)
	}
}

func TestSizeOfUnsignedVarint(t *testing.T) {
	tests := []struct {
		value    uint32
		expected int
	}{
		{0, 1},
		{127, 1},
		{128, 2},
		{16383, 2},
		{16384, 3},
		{1<<31 - 1, 5}, // testing large value
	}

	for _, test := range tests {
		if got := SizeOfUnsignedVarint(test.value); got != test.expected {
			t.Errorf("SizeOfUnsignedVarint(%d) = %d; want %d", test.value, got, test.expected)
		}
	}
}

func TestSizeOfVarint(t *testing.T) {
	tests := []struct {
		value    int32
		expected int
	}{
		{0, 1},
		{127, 1},
		{128, 2},
		{16383, 2},
		{16384, 3},
		{1<<31 - 1, 5}, // testing large value
	}

	for _, test := range tests {
		if got := SizeOfVarint(test.value); got != test.expected {
			t.Errorf("SizeOfVarint(%d) = %d; want %d", test.value, got, test.expected)
		}
	}
}

func TestSizeOfVarlong(t *testing.T) {
	tests := []struct {
		value    int64
		expected int
	}{
		{0, 1},
		{127, 1},
		{128, 2},
		{16383, 2},
		{16384, 3},
		{1<<7 - 1, 1},
		{1<<14 - 1, 2},
		{1<<21 - 1, 3},
		{1<<28 - 1, 4},
		{1<<35 - 1, 5},
		{1<<42 - 1, 6},
		{1<<49 - 1, 7},
		{1<<56 - 1, 8},
		{1<<63 - 1, 9},
	}

	for _, test := range tests {
		if got := SizeOfVarlong(test.value); got != test.expected {
			t.Errorf("SizeOfVarlong(%d) = %d; want %d", test.value, got, test.expected)
		}
	}
}

func TestReadVarint(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0xAC, 0x02})
	value, err := ReadVarint(buf)
	if err != nil {
		t.Fatal(err)
	}
	expected := int32(300)
	if value != expected {
		t.Errorf("ReadVarint() = %v; want %v", value, expected)
	}
}

func TestReadVarlong(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0xAC, 0x02})
	value, err := ReadVarlong(buf)
	if err != nil {
		t.Fatal(err)
	}
	expected := int64(300)
	if value != expected {
		t.Errorf("ReadVarlong() = %v; want %v", value, expected)
	}
}

func TestReadDouble(t *testing.T) {
	data := []byte{0x40, 0x09, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18}
	buf := bytes.NewBuffer(data)
	value, err := ReadDouble(buf)
	if err != nil {
		t.Fatal(err)
	}
	expected := 3.141592653589793
	if math.Abs(value-expected) > 1e-9 {
		t.Errorf("ReadDouble() = %v; want %v", value, expected)
	}
}

func TestWriteDouble(t *testing.T) {
	buf := new(bytes.Buffer)
	WriteDouble(3.141592653589793, buf)
	expected := []byte{0x40, 0x09, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Errorf("WriteDouble() = %v; want %v", buf.Bytes(), expected)
	}
}
