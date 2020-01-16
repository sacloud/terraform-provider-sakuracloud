// Copyright 2016-2020 terraform-provider-sakuracloud authors
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func TestAccSakuraCloudPacketFilterRules_basic(t *testing.T) {
	resourceName := "sakuracloud_packet_filter_rules.rules"
	rand := randomName()

	var filter sacloud.PacketFilter
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckSakuraCloudPacketFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudPacketFilterRules_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudPacketFilterExists("sakuracloud_packet_filter.foobar", &filter),
					resource.TestCheckResourceAttr(resourceName, "expression.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.source_network", "192.168.2.0"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.source_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.destination_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.allow", "true"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.description", "description"),

					resource.TestCheckResourceAttr(resourceName, "expression.1.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "expression.1.source_network", "192.168.2.0"),
					resource.TestCheckResourceAttr(resourceName, "expression.1.source_port", "443"),
					resource.TestCheckResourceAttr(resourceName, "expression.1.destination_port", "443"),
					resource.TestCheckResourceAttr(resourceName, "expression.1.allow", "true"),

					resource.TestCheckResourceAttr(resourceName, "expression.2.protocol", "ip"),
					resource.TestCheckResourceAttr(resourceName, "expression.2.source_network", ""),
					resource.TestCheckResourceAttr(resourceName, "expression.2.source_port", ""),
					resource.TestCheckResourceAttr(resourceName, "expression.2.destination_port", ""),
					resource.TestCheckResourceAttr(resourceName, "expression.2.allow", "false"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudPacketFilterRules_update, rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expression.0.protocol", "udp"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.source_network", "192.168.2.2"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.source_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.destination_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.allow", "true"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.description", ""),

					resource.TestCheckResourceAttr(resourceName, "expression.1.protocol", "udp"),
					resource.TestCheckResourceAttr(resourceName, "expression.1.source_network", "192.168.2.2"),
					resource.TestCheckResourceAttr(resourceName, "expression.1.source_port", "443"),
					resource.TestCheckResourceAttr(resourceName, "expression.1.destination_port", "443"),
					resource.TestCheckResourceAttr(resourceName, "expression.1.allow", "true"),

					resource.TestCheckResourceAttr(resourceName, "expression.2.protocol", "ip"),
					resource.TestCheckResourceAttr(resourceName, "expression.2.source_network", ""),
					resource.TestCheckResourceAttr(resourceName, "expression.2.source_port", ""),
					resource.TestCheckResourceAttr(resourceName, "expression.2.destination_port", ""),
					resource.TestCheckResourceAttr(resourceName, "expression.2.allow", "false"),
				),
			},
		},
	})
}

var testAccSakuraCloudPacketFilterRules_basic = `
resource "sakuracloud_packet_filter" "foobar" {
  name        = "{{ .arg0 }}"
}

resource sakuracloud_packet_filter_rules "rules" {
  packet_filter_id = sakuracloud_packet_filter.foobar.id
  expression {
 	protocol         = "tcp"
	source_network   = "192.168.2.0"
	source_port      = "80"
	destination_port = "80"
	allow            = true
    description      = "description"
  }
  expression {
	protocol         = "tcp"
	source_network   = "192.168.2.0"
	source_port      = "443"
	destination_port = "443"
	allow            = true
  }
  expression {
 	protocol = "ip"
	allow    = false
  }
}
`

var testAccSakuraCloudPacketFilterRules_update = `
resource "sakuracloud_packet_filter" "foobar" {
  name = "{{ .arg0 }}"
}

resource sakuracloud_packet_filter_rules "rules" {
  packet_filter_id = sakuracloud_packet_filter.foobar.id
  expression {
   	protocol         = "udp"
  	source_network   = "192.168.2.2"
  	source_port      = "80"
  	destination_port = "80"
   	allow            = true
  }
  expression {
   	protocol         = "udp"
  	source_network   = "192.168.2.2"
  	source_port      = "443"
  	destination_port = "443"
  	allow            = true
  }
  expression {
  	protocol = "ip"
	allow    = false
  }
}
`
