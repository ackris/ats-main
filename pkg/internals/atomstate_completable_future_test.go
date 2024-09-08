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

package internals

import (
	"errors"
	"testing"
	"time"
)

func TestAtomstateCompletableFuture(t *testing.T) {
	t.Run("TestComplete", func(t *testing.T) {
		future := NewAtomstateCompletableFuture[int]()
		err := future.Complete(42)
		if err != nil {
			t.Errorf("Complete() returned an error: %v", err)
		}

		value, err := future.Get()
		if err != nil {
			t.Errorf("Get() returned an error: %v", err)
		}
		if value != 42 {
			t.Errorf("Expected value 42, got %d", value)
		}
	})

	t.Run("TestCompleteExceptionally", func(t *testing.T) {
		future := NewAtomstateCompletableFuture[int]()
		err := future.CompleteExceptionally(errors.New("test error"))
		if err != nil {
			t.Errorf("CompleteExceptionally() returned an error: %v", err)
		}

		_, err = future.Get()
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if err.Error() != "test error" {
			t.Errorf("Expected error 'test error', got '%v'", err)
		}
	})

	t.Run("TestCompleteAfterComplete", func(t *testing.T) {
		future := NewAtomstateCompletableFuture[int]()
		err := future.Complete(42)
		if err != nil {
			t.Errorf("Complete() returned an error: %v", err)
		}

		err = future.Complete(24)
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if err.Error() != "future is already complete" {
			t.Errorf("Expected error 'future is already complete', got '%v'", err)
		}
	})

	t.Run("TestCompleteAsyncAndGet", func(t *testing.T) {
		future := NewAtomstateCompletableFuture[int]()

		err := future.CompleteAsync(func() int {
			time.Sleep(100 * time.Millisecond)
			return 42
		}, 200*time.Millisecond)
		if err != nil {
			t.Errorf("CompleteAsync() returned an error: %v", err)
		}

		value, err := future.Get()
		if err != nil {
			t.Errorf("Get() returned an error: %v", err)
		}
		if value != 42 {
			t.Errorf("Expected value 42, got %d", value)
		}
	})

	t.Run("TestCompleteOnTimeout", func(t *testing.T) {
		future := NewAtomstateCompletableFuture[int]()

		future.CompleteOnTimeout(42, 100*time.Millisecond)

		value, err := future.Get()
		if err != nil {
			t.Errorf("Get() returned an error: %v", err)
		}
		if value != 42 {
			t.Errorf("Expected value 42, got %d", value)
		}
	})
}
