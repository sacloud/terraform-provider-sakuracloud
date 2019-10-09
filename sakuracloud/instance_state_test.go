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
