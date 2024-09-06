package common

import (
	"testing"
)

// TestNewUuid tests the NewUuid function for correct UUID creation.
func TestNewUuid(t *testing.T) {
	mostSigBits := uint64(0x123456789abcdef0)
	leastSigBits := uint64(0xfedcba9876543210)
	u := NewUuid(mostSigBits, leastSigBits)
	if u.mostSignificantBits != mostSigBits || u.leastSignificantBits != leastSigBits {
		t.Errorf("NewUuid() = %v; want mostSignificantBits = %d and leastSignificantBits = %d", u, mostSigBits, leastSigBits)
	}
}

// TestRandomUuid tests the RandomUuid function for generating non-reserved UUIDs and valid Base64 encoding.
func TestRandomUuid(t *testing.T) {
	u := RandomUuid()
	if _, reserved := Reserved[u]; reserved {
		t.Errorf("RandomUuid() returned a reserved UUID: %v", u)
	}
	if startsWithDash(u.String()) {
		t.Errorf("RandomUuid() returned a UUID that starts with a dash: %v", u)
	}
}

// TestFromString tests the FromString function for correct Base64 UUID decoding and error handling.
func TestFromString(t *testing.T) {
	// Example valid Base64 representation of a UUID
	uuidStr := "EjRWeJq83vD-3LqYdlQyEA" // Adjust to valid Base64 string

	mostSigBits := uint64(0x123456789abcdef0)
	leastSigBits := uint64(0xfedcba9876543210)
	expected := NewUuid(mostSigBits, leastSigBits)

	u, err := FromString(uuidStr)
	if err != nil {
		t.Fatalf("FromString() error = %v; want no error", err)
	}
	if u.mostSignificantBits != expected.mostSignificantBits || u.leastSignificantBits != expected.leastSignificantBits {
		t.Errorf("FromString() = %v; want %v", u, expected)
	}

	// Test invalid input
	invalidUUIDStr := "invalid-uuid"
	_, err = FromString(invalidUUIDStr)
	if err == nil {
		t.Error("FromString() did not return error for invalid UUID string")
	}

	// Test too long input
	longUUIDStr := "EjRWeJq83vD-3LqYdlQyEA--" // Modified Base64 string with extra padding
	_, err = FromString(longUUIDStr)
	if err == nil {
		t.Error("FromString() did not return error for too long UUID string")
	}
}

// TestUuidString tests the String method for correct Base64 encoding of the UUID.
func TestUuidString(t *testing.T) {
	mostSigBits := uint64(0x123456789abcdef0)
	leastSigBits := uint64(0xfedcba9876543210)
	u := NewUuid(mostSigBits, leastSigBits)
	expected := "EjRWeJq83vD-3LqYdlQyEA" // Adjust this based on correct Base64 encoding
	if s := u.String(); s != expected {
		t.Errorf("String() = %v; want %v", s, expected)
	}
}

// TestGetBytes tests the getBytes method for correct byte slice representation of the UUID.
func TestGetBytes(t *testing.T) {
	mostSigBits := uint64(0x123456789abcdef0)
	leastSigBits := uint64(0xfedcba9876543210)
	u := NewUuid(mostSigBits, leastSigBits)
	expected := []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
	if b := u.getBytes(); !equalBytes(b, expected) {
		t.Errorf("getBytes() = %v; want %v", b, expected)
	}
}

// TestCompare tests the Compare method for correct lexicographical comparison of UUIDs.
func TestCompare(t *testing.T) {
	mostSigBits1 := uint64(0x123456789abcdef0)
	leastSigBits1 := uint64(0xfedcba9876543210)
	mostSigBits2 := uint64(0x123456789abcdef0)
	leastSigBits2 := uint64(0xfedcba9876543210) // Same as u1
	mostSigBits3 := uint64(0x123456789abcdef1)
	leastSigBits3 := uint64(0xfedcba9876543210)
	u1 := NewUuid(mostSigBits1, leastSigBits1)
	u2 := NewUuid(mostSigBits2, leastSigBits2) // Same as u1
	u3 := NewUuid(mostSigBits3, leastSigBits3)

	if cmp := u1.Compare(u2); cmp != 0 {
		t.Errorf("Compare(u1, u2) = %d; want 0", cmp)
	}
	if cmp := u1.Compare(u3); cmp != -1 {
		t.Errorf("Compare(u1, u3) = %d; want -1", cmp)
	}
	if cmp := u2.Compare(u1); cmp != 0 {
		t.Errorf("Compare(u2, u1) = %d; want 0", cmp)
	}
	if cmp := u1.Compare(u1); cmp != 0 {
		t.Errorf("Compare(u1, u1) = %d; want 0", cmp)
	}
}

// TestToArray tests the ToArray function for correct conversion of UUID slice to array of [2]uint64.
func TestToArray(t *testing.T) {
	mostSigBits1 := uint64(0)
	leastSigBits1 := uint64(1)
	mostSigBits2 := uint64(0)
	leastSigBits2 := uint64(2)
	uuids := []Uuid{
		NewUuid(mostSigBits1, leastSigBits1),
		NewUuid(mostSigBits2, leastSigBits2),
	}
	expected := [][2]uint64{{0, 1}, {0, 2}}
	if arr := ToArray(uuids); !equalArrays(arr, expected) {
		t.Errorf("ToArray() = %v; want %v", arr, expected)
	}
}

// TestToList tests the ToList function for correct conversion of array of [2]uint64 to UUID slice.
func TestToList(t *testing.T) {
	arr := [][2]uint64{{0, 1}, {0, 2}}
	expected := []Uuid{
		NewUuid(0, 1),
		NewUuid(0, 2),
	}
	if slice := ToList(arr); !equalUUIDSlices(slice, expected) {
		t.Errorf("ToList() = %v; want %v", slice, expected)
	}
}

// Helper functions for testing
func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalArrays(a, b [][2]uint64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalUUIDSlices(a, b []Uuid) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
