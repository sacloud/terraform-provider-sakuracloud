// Copyright 2016-2022 terraform-provider-sakuracloud authors
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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSakuraCloudDataSourcePacketFilter_basic(t *testing.T) {
	resourceName := "data.sakuracloud_packet_filter.foobar"
	rand := randomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourcePacketFilter_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.source_network", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.source_port", "0-65535"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.destination_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.allow", "true"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourcePacketFilter_basic = `
resource "sakuracloud_packet_filter" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  expression {
    protocol         = "tcp"
    source_network   = "0.0.0.0/0"
    source_port      = "0-65535"
    destination_port = "80"
    allow            = true
  }
  expression {
    protocol         = "udp"
    source_network   = "0.0.0.0/0"
    source_port      = "0-65535"
    destination_port = "80"
    allow            = true
  }
}

data "sakuracloud_packet_filter" "foobar" {
  filter {
	names = [sakuracloud_packet_filter.foobar.name]
  }
}`
