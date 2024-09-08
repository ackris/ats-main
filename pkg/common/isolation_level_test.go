package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIsolationLevel(t *testing.T) {
	tests := []struct {
		id      byte
		wantErr bool
	}{
		{0, false},  // ReadUncommitted
		{1, false},  // ReadCommitted
		{2, true},   // Invalid id
		{255, true}, // Invalid id
	}

	for _, tt := range tests {
		t.Run("Testing NewIsolationLevel with id: "+string(tt.id), func(t *testing.T) {
			il, err := NewIsolationLevel(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, il)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, il)
				assert.Equal(t, tt.id, il.ID())
			}
		})
	}
}

func TestIsolationLevel_String(t *testing.T) {
	tests := []struct {
		il   *IsolationLevel
		want string
	}{
		{&IsolationLevel{id: 0}, "read_uncommitted"},
		{&IsolationLevel{id: 1}, "read_committed"},
		{&IsolationLevel{id: 2}, "unknown"}, // Testing edge case
	}

	for _, tt := range tests {
		t.Run("Testing String method", func(t *testing.T) {
			assert.Equal(t, tt.want, tt.il.String())
		})
	}
}

func TestIsolationLevel_ID(t *testing.T) {
	tests := []struct {
		il   *IsolationLevel
		want byte
	}{
		{&IsolationLevel{id: 0}, 0},
		{&IsolationLevel{id: 1}, 1},
	}

	for _, tt := range tests {
		t.Run("Testing ID method", func(t *testing.T) {
			assert.Equal(t, tt.want, tt.il.ID())
		})
	}
}

func TestForID(t *testing.T) {
	tests := []struct {
		id      byte
		wantErr bool
	}{
		{0, false}, // ReadUncommitted
		{1, false}, // ReadCommitted
		{2, true},  // Invalid id
	}

	for _, tt := range tests {
		t.Run("Testing ForID with id: "+string(tt.id), func(t *testing.T) {
			il, err := ForID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, il)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, il)
				assert.Equal(t, tt.id, il.ID())
			}
		})
	}
}
