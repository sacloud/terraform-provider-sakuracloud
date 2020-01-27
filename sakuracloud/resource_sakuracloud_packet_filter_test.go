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
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func TestAccSakuraCloudPacketFilter_basic(t *testing.T) {
	resourceName := "sakuracloud_packet_filter.foobar"
	rand := randomName()

	var filter sacloud.PacketFilter
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckSakuraCloudPacketFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudPacketFilter_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudPacketFilterExists(resourceName, &filter),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.source_network", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.source_port", "0-65535"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.destination_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.allow", "true"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudPacketFilter_update, rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "description-upd"),
					resource.TestCheckResourceAttr(resourceName, "expression.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.source_network", "192.168.2.0"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.source_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.destination_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "expression.0.allow", "false"),
					resource.TestCheckResourceAttr(resourceName, "expression.4.protocol", "ip"),
					resource.TestCheckResourceAttr(resourceName, "expression.4.allow", "true"),
				),
			},
		},
	})
}

func testCheckSakuraCloudPacketFilterExists(n string, filter *sacloud.PacketFilter) resource.TestCheckFunc {
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

		foundPacketFilter, err := pfOp.Read(context.Background(), zone, sakuraCloudID(rs.Primary.ID))
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

func testCheckSakuraCloudPacketFilterDestroy(s *terraform.State) error {
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
		_, err := pfOp.Read(context.Background(), zone, sakuraCloudID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("still exists PacketFilter: %s", rs.Primary.ID)
		}
	}

	return nil
}

var testAccSakuraCloudPacketFilter_basic = `
resource "sakuracloud_packet_filter" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  expression {
    protocol         = "tcp"
    source_network   = "0.0.0.0"
    source_port      = "0-65535"
    destination_port = "80"
    allow            = true
  }
  expression {
    protocol         = "udp"
    source_network   = "0.0.0.0"
    source_port      = "0-65535"
    destination_port = "80"
    allow            = true
  }
}`

var testAccSakuraCloudPacketFilter_update = `
resource "sakuracloud_packet_filter" "foobar" {
  name        = "{{ .arg0 }}-upd"
  description = "description-upd"
  expression {
    protocol         = "tcp"
    source_network   = "192.168.2.0"
    source_port      = "8080"
    destination_port = "8080"
    allow            = false
  }
  expression {
    protocol         = "udp"
    source_network   = "0.0.0.0"
    source_port      = "0-65535"
    destination_port = "80"
    allow            = true
  }
  expression {
    protocol       = "icmp"
    source_network = "0.0.0.0"
    allow          = true
  }
  expression {
    protocol       = "fragment"
    source_network = "0.0.0.0"
    allow          = true
  }
  expression {
    protocol       = "ip"
    source_network = "0.0.0.0"
    allow          = true
  }
}`
