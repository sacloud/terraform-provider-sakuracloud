// Copyright 2016-2023 The sacloud/terraform-provider-sakuracloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isValidHostName(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		want     bool
	}{
		{
			name:     "empty",
			hostname: "",
			want:     true,
		},
		{
			name:     "starts with dash",
			hostname: "-example.com",
			want:     false,
		},
		{
			name:     "ends with dash",
			hostname: "example-.com",
			want:     false,
		},
		{
			name:     "with dash",
			hostname: "exa-mple.com",
			want:     true,
		},
		{
			name:     "with multiple dash",
			hostname: "exa-m-ple.com",
			want:     true,
		},
		{
			name:     "with consecutive dashes",
			hostname: "exa--mple.com",
			want:     false,
		},
		{
			name:     "ends with under bar",
			hostname: "example_.com",
			want:     false,
		},
		{
			name:     "dot+dot",
			hostname: "example..com",
			want:     false,
		},
		{
			name:     "starts with dot",
			hostname: ".example.com",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, isValidHostName(tt.hostname), "isValidHostName(%v)", tt.hostname)
		})
	}
}

func Test_isValidNameLengthBetween(t *testing.T) {
	validateFunc := isValidNameLengthBetween(3, 64)

	tests := []struct {
		name  string
		input interface{}
		want  bool
	}{
		{
			name:  "Valid string length within range",
			input: "こんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちは",
			want:  false,
		},
		{
			name:  "String too short",
			input: "こん",
			want:  true,
		},
		{
			name:  "String too long",
			input: "こんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちは",
			want:  true,
		},
		{
			name:  "Empty string",
			input: "",
			want:  true,
		},
		{
			name:  "Exact minimum length",
			input: "こんに",
			want:  false,
		},
		{
			name:  "Exact maximum length",
			input: "こんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにちはこんにち",
			want:  false,
		},
		{
			name:  "Non-string input",
			input: 12345,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errors := validateFunc(tt.input, "test_field")

			if tt.want {
				assert.NotEmpty(t, errors, "Expected validation errors but got none")
			} else {
				assert.Empty(t, errors, "Expected no validation errors but got some")
			}
		})
	}
}
