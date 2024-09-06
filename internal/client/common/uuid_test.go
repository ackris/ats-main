package common

import (
	"math/big"
	"testing"
)

func TestNewUuid(t *testing.T) {
	mostSigBits := big.NewInt(0).SetUint64(0x123456789abcdef0)
	leastSigBits := big.NewInt(0).SetUint64(0xfedcba9876543210)
	u := NewUuid(mostSigBits.Int64(), leastSigBits.Int64())
	if u.mostSignificantBits != mostSigBits.Int64() || u.leastSignificantBits != leastSigBits.Int64() {
		t.Errorf("NewUuid() = %v; want mostSignificantBits = %d and leastSignificantBits = %d", u, mostSigBits.Int64(), leastSigBits.Int64())
	}
}

func TestRandomUuid(t *testing.T) {
	u := RandomUuid()
	if _, reserved := Reserved[u]; reserved {
		t.Errorf("RandomUuid() returned a reserved UUID: %v", u)
	}
	if startsWithDash(u.String()) {
		t.Errorf("RandomUuid() returned a UUID that starts with a dash: %v", u)
	}
}

func TestFromString(t *testing.T) {
	uuidStr := "0Aw-Aw0AQNaAQNaAQNaAQNaAQ" // Base64 representation of a UUID
	u, err := FromString(uuidStr)
	if err != nil {
		t.Fatalf("FromString() error = %v; want no error", err)
	}
	expected := NewUuid(int64(0), int64(1)) // Expected UUID
	if u != expected {
		t.Errorf("FromString() = %v; want %v", u, expected)
	}

	// Test invalid input
	_, err = FromString("invalid-uuid")
	if err == nil {
		t.Error("FromString() did not return error for invalid UUID string")
	}
}

func TestUuidString(t *testing.T) {
	mostSigBits := big.NewInt(0).SetUint64(0x123456789abcdef0)
	leastSigBits := big.NewInt(0).SetUint64(0xfedcba9876543210)
	u := NewUuid(mostSigBits.Int64(), leastSigBits.Int64())
	expected := "EjRWeWsFUhQ" // Base64 encoding of the UUID
	if s := u.String(); s != expected {
		t.Errorf("String() = %v; want %v", s, expected)
	}
}

func TestGetBytes(t *testing.T) {
	mostSigBits := big.NewInt(0).SetUint64(0x123456789abcdef0)
	leastSigBits := big.NewInt(0).SetUint64(0xfedcba9876543210)
	u := NewUuid(mostSigBits.Int64(), leastSigBits.Int64())
	expected := []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
	if b := u.getBytes(); !equalBytes(b, expected) {
		t.Errorf("getBytes() = %v; want %v", b, expected)
	}
}

func TestCompare(t *testing.T) {
	mostSigBits1 := big.NewInt(0).SetUint64(0x123456789abcdef0)
	leastSigBits1 := big.NewInt(0).SetUint64(0xfedcba9876543210)
	mostSigBits2 := big.NewInt(0).SetUint64(0x123456789abcdef0)
	leastSigBits2 := big.NewInt(0).SetUint64(0xfedcba9876543211)
	mostSigBits3 := big.NewInt(0).SetUint64(0x123456789abcdef1)
	leastSigBits3 := big.NewInt(0).SetUint64(0xfedcba9876543210)
	u1 := NewUuid(mostSigBits1.Int64(), leastSigBits1.Int64())
	u2 := NewUuid(mostSigBits2.Int64(), leastSigBits2.Int64())
	u3 := NewUuid(mostSigBits3.Int64(), leastSigBits3.Int64())

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

func TestToArray(t *testing.T) {
	mostSigBits1 := big.NewInt(0).SetUint64(0)
	leastSigBits1 := big.NewInt(0).SetUint64(1)
	mostSigBits2 := big.NewInt(0).SetUint64(0)
	leastSigBits2 := big.NewInt(0).SetUint64(2)
	uuids := []Uuid{
		NewUuid(mostSigBits1.Int64(), leastSigBits1.Int64()),
		NewUuid(mostSigBits2.Int64(), leastSigBits2.Int64()),
	}
	expected := [][2]int64{{0, 1}, {0, 2}}
	if arr := ToArray(uuids); !equalArrays(arr, expected) {
		t.Errorf("ToArray() = %v; want %v", arr, expected)
	}
}

func TestToList(t *testing.T) {
	arr := [][2]int64{{0, 1}, {0, 2}}
	expected := []Uuid{
		NewUuid(int64(0), int64(1)),
		NewUuid(int64(0), int64(2)),
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

func equalArrays(a, b [][2]int64) bool {
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
