// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var (
	testDefaultTargetZone   = "is1b"
	testDefaultAPIRetryMax  = "30"
	testDefaultAPIRateLimit = "5"
)

func init() {
	if v := os.Getenv("SAKURACLOUD_TEST_ZONE"); v != "" {
		os.Setenv("SAKURACLOUD_ZONE", v)
	}
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"sakuracloud": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
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
				t.Fatal(fmt.Sprintf("%s must be set for acceptance tests", env))
			}
		}
	}

	if v := os.Getenv("SAKURACLOUD_ZONE"); v == "" {
		os.Setenv("SAKURACLOUD_ZONE", testDefaultTargetZone)
	}

	if v := os.Getenv("SAKURACLOUD_RETRY_MAX"); v == "" {
		os.Setenv("SAKURACLOUD_RETRY_MAX", testDefaultAPIRetryMax)
	}

	if v := os.Getenv("SAKURACLOUD_RATE_LIMIT"); v == "" {
		os.Setenv("SAKURACLOUD_RATE_LIMIT", testDefaultAPIRateLimit)
	}
}
