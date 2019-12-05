package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
