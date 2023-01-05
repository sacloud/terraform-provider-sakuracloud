// Copyright 2016-2023 terraform-provider-sakuracloud authors
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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func compareState(s *terraform.InstanceState, key, value string) error {
	actual := s.Attributes[key]
	if actual != value {
		return fmt.Errorf("expected state[%s] is %q, but %q received",
			key, value, actual)
	}
	return nil
}

func compareStateMulti(s *terraform.InstanceState, expects map[string]string) error {
	for k, v := range expects {
		err := compareState(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func stateNotEmpty(s *terraform.InstanceState, key string) error {
	if v, ok := s.Attributes[key]; !ok || v == "" {
		return fmt.Errorf("state[%s] is expected not empty", key)
	}
	return nil
}

func stateNotEmptyMulti(s *terraform.InstanceState, keys ...string) error {
	for _, key := range keys {
		if err := stateNotEmpty(s, key); err != nil {
			return err
		}
	}
	return nil
}
