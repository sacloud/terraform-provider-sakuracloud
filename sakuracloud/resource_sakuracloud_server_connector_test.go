package sakuracloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "cdrom_id", "",
					),
				),
			},
		},
	})
}

const testAccCheckSakuraCloudServerConnectorConfig_basic = `
resource sakuracloud_server "foobar" {
  name = "foobar"
  additional_nics = [""]
}
resource sakuracloud_disk "foobar" {
  name = "foobar"
  count = 1
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
resource sakuracloud_server "foobar" {
  name = "foobar"
  additional_nics = [""]
}
resource sakuracloud_disk "foobar" {
  name = "foobar"
  count = 2
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
resource sakuracloud_server "foobar" {
  name = "foobar"
  additional_nics = [""]
}
resource sakuracloud_disk "foobar" {
  name = "foobar"
  count = 2
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
