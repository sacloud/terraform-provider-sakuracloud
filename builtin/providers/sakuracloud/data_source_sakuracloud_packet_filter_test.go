package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudPacketFilterDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudPacketFilterDataSourceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourcePacketFilterBase,
				Check:  testAccCheckSakuraCloudPacketFilterDataSourceID("sakuracloud_packet_filter.foobar"),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourcePacketFilterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPacketFilterDataSourceID("data.sakuracloud_packet_filter.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.#", "2"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.source_nw", "0.0.0.0"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.source_port", "0-65535"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.dest_port", "80"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.allow", "true"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourcePacketFilterConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPacketFilterDataSourceNotExists("data.sakuracloud_packet_filter.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudPacketFilterDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find PacketFilter data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("PacketFilter data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudPacketFilterDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found PacketFilter data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudPacketFilterDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_packet_filter" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.PacketFilter.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("PacketFilter still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourcePacketFilterBase = `
resource "sakuracloud_packet_filter" "foobar" {
    name = "name_test"
    description = "description_test"
    expressions = {
    	protocol = "tcp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
    expressions = {
    	protocol = "udp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
}
`

var testAccCheckSakuraCloudDataSourcePacketFilterConfig = `
resource "sakuracloud_packet_filter" "foobar" {
    name = "name_test"
    description = "description_test"
    expressions = {
    	protocol = "tcp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
    expressions = {
    	protocol = "udp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
}
data "sakuracloud_packet_filter" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`

var testAccCheckSakuraCloudDataSourcePacketFilterConfig_NotExists = `
resource "sakuracloud_packet_filter" "foobar" {
    name = "name_test"
    description = "description_test"
    expressions = {
    	protocol = "tcp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
    expressions = {
    	protocol = "udp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
}
data "sakuracloud_packet_filter" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`
