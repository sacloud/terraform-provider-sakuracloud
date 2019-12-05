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
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func TestAccResourceSakuraCloudPacketFilter(t *testing.T) {
	var filter sacloud.PacketFilter
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudPacketFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudPacketFilterConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPacketFilterExists("sakuracloud_packet_filter.foobar", &filter),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "name", "mypacket_filter"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.source_network", "0.0.0.0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.source_port", "0-65535"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.destination_port", "80"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.allow", "true"),
				),
			},
			{
				Config: testAccCheckSakuraCloudPacketFilterConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "name", "mypacket_filter_upd"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.#", "5"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.source_network", "192.168.2.0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.source_port", "8080"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.destination_port", "8080"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.allow", "false"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.4.protocol", "ip"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.4.allow", "true"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudPacketFilterExists(n string, filter *sacloud.PacketFilter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no PacketFilter ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		zone := rs.Primary.Attributes["zone"]
		pfOp := sacloud.NewPacketFilterOp(client)

		foundPacketFilter, err := pfOp.Read(context.Background(), zone, types.StringID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if foundPacketFilter.ID.String() != rs.Primary.ID {
			return fmt.Errorf("not found PacketFilter: %s", rs.Primary.ID)
		}

		*filter = *foundPacketFilter
		return nil
	}
}

func testAccCheckSakuraCloudPacketFilterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	pfOp := sacloud.NewPacketFilterOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_packet_filter" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		zone := rs.Primary.Attributes["zone"]
		_, err := pfOp.Read(context.Background(), zone, types.StringID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("still exists PacketFilter: %s", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCheckSakuraCloudPacketFilterConfig_basic = `
resource "sakuracloud_packet_filter" "foobar" {
  name        = "mypacket_filter"
  description = "PacketFilter from TerraForm for SAKURA CLOUD"
  expressions {
    protocol         = "tcp"
    source_network   = "0.0.0.0"
    source_port      = "0-65535"
    destination_port = "80"
    allow            = true
  }
  expressions {
    protocol         = "udp"
    source_network   = "0.0.0.0"
    source_port      = "0-65535"
    destination_port = "80"
    allow            = true
  }
}`

var testAccCheckSakuraCloudPacketFilterConfig_update = `
resource "sakuracloud_packet_filter" "foobar" {
  name        = "mypacket_filter_upd"
  description = "PacketFilter from TerraForm for SAKURA CLOUD"
  expressions {
    protocol         = "tcp"
    source_network   = "192.168.2.0"
    source_port      = "8080"
    destination_port = "8080"
    allow            = false
  }
  expressions {
    protocol         = "udp"
    source_network   = "0.0.0.0"
    source_port      = "0-65535"
    destination_port = "80"
    allow            = true
  }
  expressions {
    protocol       = "icmp"
    source_network = "0.0.0.0"
    allow          = true
  }
  expressions {
    protocol       = "fragment"
    source_network = "0.0.0.0"
    allow          = true
  }
  expressions {
    protocol       = "ip"
    source_network = "0.0.0.0"
    allow          = true
  }
}`
