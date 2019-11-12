// Copyright 2016-2019 terraform-provider-sakuracloud authors
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

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func TestStringSliceFromState(t *testing.T) {

	expects := []struct {
		should      string
		state       *terraform.InstanceState
		key         string
		expectValue []string
	}{
		{
			should:      "empty",
			state:       &terraform.InstanceState{},
			key:         "foobar",
			expectValue: []string{},
		},
		{
			should: "slice",
			state: &terraform.InstanceState{
				Attributes: map[string]string{
					"foobar.#": "2",
					"foobar.0": "foobar.0",
					"foobar.1": "foobar.1",
				},
			},
			key:         "foobar",
			expectValue: []string{"foobar.0", "foobar.1"},
		},
	}

	for _, expect := range expects {
		t.Run("Should "+expect.should, func(t *testing.T) {
			state := StringSliceFromState(expect.state, expect.key)
			assert.EqualValues(t, expect.expectValue, state)
		})
	}

}
