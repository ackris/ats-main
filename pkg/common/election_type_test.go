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

func TestElectionType_String(t *testing.T) {
	tests := []struct {
		name string
		et   ElectionType
		want string
	}{
		{"Preferred", Preferred, "Preferred"},
		{"Unclean", Unclean, "Unclean"},
		{"Unknown", ElectionType(2), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.et.String(); got != tt.want {
				t.Errorf("ElectionType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueOf(t *testing.T) {
	tests := []struct {
		name    string
		value   byte
		want    ElectionType
		wantErr bool
	}{
		{"Preferred", 0, Preferred, false},
		{"Unclean", 1, Unclean, false},
		{"Invalid value", 2, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValueOf(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValueOf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValueOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
