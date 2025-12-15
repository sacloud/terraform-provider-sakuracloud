// Copyright 2016-2025 terraform-provider-sakuracloud authors
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
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/api-client-go/profile"
	"github.com/sacloud/packages-go/envvar"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/defaults"
)

var testAccProviderFactories map[string]func() (*schema.Provider, error)
var ProtoV5ProviderFactories map[string]func() (tfprotov5.ProviderServer, error)

var testAccProvider *schema.Provider
var (
	testDefaultTargetZone   = "is1b"
	testDefaultAPIRetryMax  = "30"
	testDefaultAPIRateLimit = "5"
)

func init() {
	if v := os.Getenv("SAKURACLOUD_TEST_ZONE"); v != "" {
		os.Setenv("SAKURACLOUD_ZONE", v) //nolint:errcheck,gosec
	}
	testAccProvider = Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"sakuracloud": func() (*schema.Provider, error) { return testAccProvider, nil },
	}
	ProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
		"sakuracloud": func() (tfprotov5.ProviderServer, error) {
			providerServerFactory, err := ProtoV5ProviderServerFactory(context.Background())
			if err != nil {
				return nil, err
			}
			return providerServerFactory(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderSchema(t *testing.T) {
	provider := Provider()

	tests := []struct {
		fieldName  string
		oldEnvName string
		oldEnvVal  string
		newEnvName string
		newEnvVal  string
		defaultVal string
	}{
		{
			fieldName:  "profile",
			oldEnvName: "SAKURACLOUD_PROFILE",
			oldEnvVal:  "foo",
			newEnvName: "SAKURA_PROFILE",
			newEnvVal:  "bar",
			defaultVal: profile.DefaultProfileName,
		},
		{
			fieldName:  "token",
			oldEnvName: "SAKURACLOUD_ACCESS_TOKEN",
			oldEnvVal:  "foo",
			newEnvName: "SAKURA_ACCESS_TOKEN",
			newEnvVal:  "bar",
		},
		{
			fieldName:  "secret",
			oldEnvName: "SAKURACLOUD_ACCESS_TOKEN_SECRET",
			oldEnvVal:  "foo",
			newEnvName: "SAKURA_ACCESS_TOKEN_SECRET",
			newEnvVal:  "bar",
		},
		{
			fieldName:  "zone",
			oldEnvName: "SAKURACLOUD_ZONE",
			oldEnvVal:  "foo",
			newEnvName: "SAKURA_ZONE",
			newEnvVal:  "bar",
			defaultVal: defaults.Zone,
		},
		{
			fieldName:  "default_zone",
			oldEnvName: "SAKURACLOUD_DEFAULT_ZONE",
			oldEnvVal:  "foo",
			newEnvName: "SAKURA_DEFAULT_ZONE",
			newEnvVal:  "bar",
		},
		{
			fieldName:  "accept_language",
			oldEnvName: "SAKURACLOUD_ACCEPT_LANGUAGE",
			oldEnvVal:  "foo",
			newEnvName: "SAKURA_ACCEPT_LANGUAGE",
			newEnvVal:  "bar",
		},
		{
			fieldName:  "api_root_url",
			oldEnvName: "SAKURACLOUD_API_ROOT_URL",
			oldEnvVal:  "foo",
			newEnvName: "SAKURA_API_ROOT_URL",
			newEnvVal:  "bar",
		},
	}
	for _, tt := range tests {
		cleanup := func() {
			os.Unsetenv(tt.oldEnvName) //nolint:errcheck
			os.Unsetenv(tt.newEnvName) //nolint:errcheck
		}

		field, ok := provider.Schema[tt.fieldName]
		if !ok {
			t.Fatalf("field %s not found in provider schema", tt.fieldName)
		}

		t.Run("only old env is set", func(t *testing.T) {
			cleanup()

			if err := os.Setenv(tt.oldEnvName, tt.oldEnvVal); err != nil {
				t.Fatalf("setting env failed: %s", err)
			}

			value, err := field.DefaultValue()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(value, tt.oldEnvVal) {
				t.Errorf("expected %v, got %v", tt.oldEnvVal, value)
			}
		})
		t.Run("only new env is set", func(t *testing.T) {
			cleanup()

			if err := os.Setenv(tt.newEnvName, tt.newEnvVal); err != nil {
				t.Fatalf("setting env failed: %s", err)
			}

			value, err := field.DefaultValue()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(value, tt.newEnvVal) {
				t.Errorf("expected %v, got %v", tt.newEnvVal, value)
			}
		})
		t.Run("both of old and new are set", func(t *testing.T) {
			cleanup()

			if err := os.Setenv(tt.oldEnvName, tt.oldEnvVal); err != nil {
				t.Fatalf("setting env failed: %s", err)
			}
			if err := os.Setenv(tt.newEnvName, tt.newEnvVal); err != nil {
				t.Fatalf("setting env failed: %s", err)
			}

			value, err := field.DefaultValue()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(value, tt.newEnvVal) {
				t.Errorf("expected %v, got %v", tt.newEnvVal, value)
			}
		})
		t.Run("both of old and new are empty", func(t *testing.T) {
			cleanup()

			value, err := field.DefaultValue()
			if err != nil {
				t.Fatal(err)
			}

			if value == nil {
				value = ""
			}

			if !reflect.DeepEqual(value, tt.defaultVal) {
				t.Errorf("expected %v, got %v", tt.newEnvVal, value)
			}
		})
	}
}

func testAccPreCheck(t *testing.T) {
	requiredEnvs := [][]string{
		{"SAKURA_ACCESS_TOKEN", "SAKURACLOUD_ACCESS_TOKEN"},
		{"SAKURA_ACCESS_TOKEN_SECRET", "SAKURACLOUD_ACCESS_TOKEN_SECRET"},
	}

	if isFakeModeEnabled() {
		for _, envs := range requiredEnvs {
			for _, env := range envs {
				if err := os.Setenv(env, "dummy"); err != nil {
					t.Fatalf("setting up dummy environment variables is failed: %s", err)
				}
			}
		}
	} else {
		for _, envs := range requiredEnvs {
			if v := envvar.StringFromEnvMulti(envs, ""); v == "" {
				t.Fatalf("%s must be set for acceptance tests", envs)
			}
		}
	}

	if v := envvar.StringFromEnvMulti([]string{"SAKURA_ZONE", "SAKURACLOUD_ZONE"}, ""); v == "" {
		os.Setenv("SAKURAD_ZONE", testDefaultTargetZone) //nolint:errcheck,gosec
	}

	if v := envvar.StringFromEnvMulti([]string{"SAKURA_RETRY_MAX", "SAKURACLOUD_RETRY_MAX"}, ""); v == "" {
		os.Setenv("SAKURA_RETRY_MAX", testDefaultAPIRetryMax) //nolint:errcheck,gosec
	}

	if v := envvar.StringFromEnvMulti([]string{"SAKURA_RATE_LIMIT", "SAKURACLOUD_RATE_LIMIT"}, ""); v == "" {
		os.Setenv("SAKURA_RATE_LIMIT", testDefaultAPIRateLimit) //nolint:errcheck,gosec
	}
}
