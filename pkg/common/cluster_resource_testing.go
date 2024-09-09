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

package common

import "testing"

func TestNewClusterResource(t *testing.T) {
	t.Run("with valid cluster ID", func(t *testing.T) {
		cr := NewClusterResource("my-cluster-id")
		if cr == nil {
			t.Errorf("NewClusterResource returned nil")
		} else if cr.clusterID != "my-cluster-id" {
			t.Errorf("Expected cluster ID 'my-cluster-id', got '%s'", cr.clusterID)
		}
	})

	t.Run("with empty cluster ID", func(t *testing.T) {
		cr := NewClusterResource("")
		if cr == nil {
			t.Errorf("NewClusterResource returned nil")
		} else if cr.clusterID != "" {
			t.Errorf("Expected empty cluster ID, got '%s'", cr.clusterID)
		}
	})
}

func TestClusterResource_ClusterID(t *testing.T) {
	cr := NewClusterResource("my-cluster-id")
	id := cr.ClusterID()
	if id != "my-cluster-id" {
		t.Errorf("Expected cluster ID 'my-cluster-id', got '%s'", id)
	}
}

func TestClusterResource_String(t *testing.T) {
	cr := NewClusterResource("my-cluster-id")
	str := cr.String()
	expected := "ClusterResource(ClusterID=my-cluster-id)"
	if str != expected {
		t.Errorf("Expected string '%s', got '%s'", expected, str)
	}
}

func TestClusterResource_Equals(t *testing.T) {
	t.Run("equal cluster IDs", func(t *testing.T) {
		cr1 := NewClusterResource("my-cluster-id")
		cr2 := NewClusterResource("my-cluster-id")
		if !cr1.Equals(cr2) {
			t.Errorf("Expected ClusterResource instances to be equal")
		}
	})

	t.Run("different cluster IDs", func(t *testing.T) {
		cr1 := NewClusterResource("my-cluster-id")
		cr2 := NewClusterResource("another-cluster-id")
		if cr1.Equals(cr2) {
			t.Errorf("Expected ClusterResource instances to be different")
		}
	})

	t.Run("nil comparison", func(t *testing.T) {
		cr := NewClusterResource("my-cluster-id")
		if cr.Equals(nil) {
			t.Errorf("Expected Equals to return false for nil comparison")
		}
	})
}

func TestClusterResource_Hash(t *testing.T) {
	cr1 := NewClusterResource("my-cluster-id")
	cr2 := NewClusterResource("my-cluster-id")
	hash1 := cr1.Hash()
	hash2 := cr2.Hash()
	if hash1 != hash2 {
		t.Errorf("Expected same cluster IDs to have the same hash, got %d and %d", hash1, hash2)
	}

	cr3 := NewClusterResource("another-cluster-id")
	hash3 := cr3.Hash()
	if hash3 == hash1 {
		t.Errorf("Expected different cluster IDs to have different hashes, got %d and %d", hash3, hash1)
	}
}
