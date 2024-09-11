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
	"testing"
)

func TestNewBytes(t *testing.T) {
	data := []byte("hello")
	b := NewBytes(data)
	if b == nil {
		t.Fatalf("Expected non-nil Bytes, got nil")
	}
	if string(b.Get()) != "hello" {
		t.Fatalf("Expected data to be 'hello', got '%s'", string(b.Get()))
	}

	b = NewBytes(nil)
	if b != nil {
		t.Fatalf("Expected nil Bytes, got non-nil")
	}
}

func TestBytes_Get(t *testing.T) {
	data := []byte("test")
	b := NewBytes(data)
	got := b.Get()
	if string(got) != "test" {
		t.Fatalf("Expected 'test', got '%s'", string(got))
	}
}

func TestBytes_HashCode(t *testing.T) {
	data1 := []byte("test")
	data2 := []byte("test")
	data3 := []byte("different")

	b1 := NewBytes(data1)
	b2 := NewBytes(data2)
	b3 := NewBytes(data3)

	if b1.HashCode() != b2.HashCode() {
		t.Fatalf("Hash codes of equal Bytes objects should be the same")
	}
	if b1.HashCode() == b3.HashCode() {
		t.Fatalf("Hash codes of different Bytes objects should be different")
	}
}

func TestBytes_Equals(t *testing.T) {
	data1 := []byte("test")
	data2 := []byte("test")
	data3 := []byte("different")

	b1 := NewBytes(data1)
	b2 := NewBytes(data2)
	b3 := NewBytes(data3)

	if !b1.Equals(b2) {
		t.Fatalf("Bytes objects with equal data should be equal")
	}
	if b1.Equals(b3) {
		t.Fatalf("Bytes objects with different data should not be equal")
	}
	if b1.Equals(nil) {
		t.Fatalf("Bytes object should not be equal to nil")
	}
	if NewBytes(nil).Equals(b1) {
		t.Fatalf("Nil Bytes object should not be equal to non-nil")
	}
}

func TestBytes_CompareTo(t *testing.T) {
	data1 := []byte("abc")
	data2 := []byte("abc")
	data3 := []byte("abd")
	data4 := []byte("ab")

	b1 := NewBytes(data1)
	b2 := NewBytes(data2)
	b3 := NewBytes(data3)
	b4 := NewBytes(data4)

	if b1.CompareTo(b2) != 0 {
		t.Fatalf("Equal Bytes objects should compare as 0")
	}
	if b1.CompareTo(b3) >= 0 {
		t.Fatalf("Bytes object with less data should compare as less than object with more data")
	}
	if b1.CompareTo(b4) <= 0 {
		t.Fatalf("Bytes object with more data should compare as greater than object with less data")
	}
	if b1.CompareTo(nil) <= 0 {
		t.Fatalf("Bytes object should compare as greater than nil")
	}
}

func TestBytes_String(t *testing.T) {
	data := []byte("hello\x00world")
	b := NewBytes(data)
	got := b.String()
	expected := "hello\\x00world"
	if got != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, got)
	}
}

func TestIncrement(t *testing.T) {
	data := []byte{0x00, 0x01, 0xFF}
	b := NewBytes(data)

	incremented, err := Increment(b)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expected := []byte{0x00, 0x02, 0x00}
	if !bytes.Equal(incremented.Get(), expected) {
		t.Fatalf("Expected incremented data to be '%v', got '%v'", expected, incremented.Get())
	}

	data = []byte{0xFF, 0xFF, 0xFF}
	b = NewBytes(data)
	_, err = Increment(b)
	if err == nil {
		t.Fatalf("Expected error for byte array overflow, got nil")
	}
}

func TestLexicographicByteArrayComparator(t *testing.T) {
	a := []byte("abc")
	b := []byte("abc")
	c := []byte("abd")
	d := []byte("ab")

	tests := []struct {
		a, b          []byte
		offsetA, lenA int
		offsetB, lenB int
		expected      int
	}{
		{a, b, 0, 3, 0, 3, 0},
		{a, c, 0, 3, 0, 3, -1},
		{a, d, 0, 3, 0, 2, 1},
		{a, b, 0, 3, 0, 2, 1},
	}

	for _, test := range tests {
		got := LexicographicByteArrayComparator(test.a, test.b, test.offsetA, test.lenA, test.offsetB, test.lenB)
		if got != test.expected {
			t.Fatalf("Expected %d, got %d", test.expected, got)
		}
	}
}
