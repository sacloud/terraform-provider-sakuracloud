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
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func testAccPreCheck(t *testing.T) {
	requiredEnvs := []string{
		"SAKURACLOUD_ACCESS_TOKEN",
		"SAKURACLOUD_ACCESS_TOKEN_SECRET",
	}

	if isFakeModeEnabled() {
		for _, env := range requiredEnvs {
			if err := os.Setenv(env, "dummy"); err != nil {
				t.Fatalf("setting up dummy environment variables is failed: %s", err)
			}
		}
	} else {
		for _, env := range requiredEnvs {
			if v := os.Getenv(env); v == "" {
				t.Fatalf("%s must be set for acceptance tests", env)
			}
		}
	}

	if v := os.Getenv("SAKURACLOUD_ZONE"); v == "" {
		os.Setenv("SAKURACLOUD_ZONE", testDefaultTargetZone) //nolint:errcheck,gosec
	}

	if v := os.Getenv("SAKURACLOUD_RETRY_MAX"); v == "" {
		os.Setenv("SAKURACLOUD_RETRY_MAX", testDefaultAPIRetryMax) //nolint:errcheck,gosec
	}

	if v := os.Getenv("SAKURACLOUD_RATE_LIMIT"); v == "" {
		os.Setenv("SAKURACLOUD_RATE_LIMIT", testDefaultAPIRateLimit) //nolint:errcheck,gosec
	}
}
