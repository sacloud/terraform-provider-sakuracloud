// Copyright 2016-2022 The sacloud/terraform-provider-sakuracloud Authors
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

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sacloud/terraform-provider-sakuracloud/version"
)

func protoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"sakuracloud": providerserver.NewProtocol5WithError(New(version.Version)()),
	}
}

func TestProvider_InvalidProviderConfig(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: `
					provider "sakuracloud" { 
						api_request_rate_limit = 11
					}
					data sakuracloud_zone "example" {}
				`,
				ExpectError: regexp.MustCompile(`value must be between 1 and 10`),
			},
			{
				Config: `
					provider "sakuracloud" { 
						retry_max = -1
					}
					data sakuracloud_zone "example" {}
				`,
				ExpectError: regexp.MustCompile(`value must be between 0 and 100`),
			},
		},
	})
}
