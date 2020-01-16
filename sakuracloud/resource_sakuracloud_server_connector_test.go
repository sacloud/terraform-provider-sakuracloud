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
)

func TestAccResourceSakuraCloudServerConnector(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudServerConnectorConfig_basic,
			},
			{
				Config: testAccCheckSakuraCloudServerConnectorConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"sakuracloud_server.foobar", "id",
						"sakuracloud_server_connector.foobar", "server_id",
					),
					resource.TestCheckResourceAttrPair(
						"sakuracloud_server.foobar", "disks.0",
						"sakuracloud_server_connector.foobar", "disks.0",
					),
					resource.TestCheckResourceAttrPair(
						"sakuracloud_server.foobar", "packet_filter_ids.0",
						"sakuracloud_server_connector.foobar", "packet_filter_ids.0",
					),
					resource.TestCheckResourceAttrPair(
						"sakuracloud_server.foobar", "cdrom_id",
						"sakuracloud_server_connector.foobar", "cdrom_id",
					),
				),
			},
			{
				Config: testAccCheckSakuraCloudServerConnectorConfig_update,
			},
			{
				Config: testAccCheckSakuraCloudServerConnectorConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"sakuracloud_server.foobar", "id",
						"sakuracloud_server_connector.foobar", "server_id",
					),
					resource.TestCheckResourceAttrPair(
						"sakuracloud_server.foobar", "disks.0",
						"sakuracloud_server_connector.foobar", "disks.0",
					),
					resource.TestCheckResourceAttrPair(
						"sakuracloud_server.foobar", "disks.1",
						"sakuracloud_server_connector.foobar", "disks.1",
					),
					resource.TestCheckResourceAttrPair(
						"sakuracloud_server.foobar", "packet_filter_ids.0",
						"sakuracloud_server_connector.foobar", "packet_filter_ids.0",
					),
					resource.TestCheckResourceAttrPair(
						"sakuracloud_server.foobar", "packet_filter_ids.1",
						"sakuracloud_server_connector.foobar", "packet_filter_ids.1",
					),
					resource.TestCheckResourceAttrPair(
						"sakuracloud_server.foobar", "cdrom_id",
						"sakuracloud_server_connector.foobar", "cdrom_id",
					),
				),
			},
			{
				Config: testAccCheckSakuraCloudServerConnectorConfig_delete,
			},
			{
				Config: testAccCheckSakuraCloudServerConnectorConfig_delete,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						"sakuracloud_server.foobar", "disks",
					),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_server.foobar", "packet_filter_ids",
					),
				),
			},
		},
	})
}

const testAccCheckSakuraCloudServerConnectorConfig_basic = `
resource sakuracloud_server "foobar" {
  name = "foobar"
  graceful_shutdown_timeout = 5
}
resource sakuracloud_disk "foobar" {
  name = "foobar"
  count = 1
  graceful_shutdown_timeout = 5
}
resource sakuracloud_packet_filter "foobar" {
  name = "foobar"
  count = 1
}
data "sakuracloud_cdrom" "foobar" {
    filter {
	name = "Name"
	values = ["Ubuntu Server 16"]
    }
}

resource sakuracloud_server_connector "foobar" {
  server_id         = "${sakuracloud_server.foobar.id}"
  disks             = ["${sakuracloud_disk.foobar.0.id}"]
  packet_filter_ids = ["${sakuracloud_packet_filter.foobar.0.id}"]
  cdrom_id          = "${data.sakuracloud_cdrom.foobar.id}"
}
`

const testAccCheckSakuraCloudServerConnectorConfig_update = `
resource sakuracloud_switch "foobar" {
  name = "foobar"
}
resource sakuracloud_server "foobar" {
  name = "foobar"
  additional_nics = ["${sakuracloud_switch.foobar.id}"]
  graceful_shutdown_timeout = 5
}
resource sakuracloud_disk "foobar" {
  name = "foobar"
  count = 2
  graceful_shutdown_timeout = 5
}
resource sakuracloud_packet_filter "foobar" {
  name = "foobar"
  count = 2
}
data "sakuracloud_cdrom" "foobar" {
    filter {
	name = "Name"
	values = ["CentOS"]
    }
}

resource sakuracloud_server_connector "foobar" {
  server_id         = "${sakuracloud_server.foobar.id}"
  disks             = ["${sakuracloud_disk.foobar.0.id}", "${sakuracloud_disk.foobar.1.id}"]
  packet_filter_ids = ["${sakuracloud_packet_filter.foobar.0.id}", "${sakuracloud_packet_filter.foobar.1.id}"]
  cdrom_id          = "${data.sakuracloud_cdrom.foobar.id}"
}
`

const testAccCheckSakuraCloudServerConnectorConfig_delete = `
resource sakuracloud_switch "foobar" {
  name = "foobar"
}
resource sakuracloud_server "foobar" {
  name = "foobar"
  graceful_shutdown_timeout = 5
}
resource sakuracloud_disk "foobar" {
  name = "foobar"
  count = 2
  graceful_shutdown_timeout = 5
}
resource sakuracloud_packet_filter "foobar" {
  name = "foobar"
  count = 2
}
data "sakuracloud_cdrom" "foobar" {
    filter {
	name = "Name"
	values = ["CentOS"]
    }
}
`
