package sakuracloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"

	"github.com/sacloud/libsacloud/sacloud"
)

func TestAccResourceSakuraCloudPacketFilter(t *testing.T) {
	var filter sacloud.PacketFilter
	resource.Test(t, resource.TestCase{
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
						"sakuracloud_packet_filter.foobar", "expressions.0.source_nw", "0.0.0.0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.source_port", "0-65535"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.dest_port", "80"),
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
						"sakuracloud_packet_filter.foobar", "expressions.0.source_nw", "192.168.2.0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.source_port", "8080"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter.foobar", "expressions.0.dest_port", "8080"),
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
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No PacketFilter ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)

		foundPacketFilter, err := client.PacketFilter.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundPacketFilter.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("PacketFilter not found")
		}

		*filter = *foundPacketFilter

		return nil
	}
}

func testAccCheckSakuraCloudPacketFilterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_packet_filter" {
			continue
		}

		_, err := client.PacketFilter.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("PacketFilter still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudPacketFilterConfig_basic = `
resource "sakuracloud_packet_filter" "foobar" {
    name = "mypacket_filter"
    description = "PacketFilter from TerraForm for SAKURA CLOUD"
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
}`

var testAccCheckSakuraCloudPacketFilterConfig_update = `
resource "sakuracloud_packet_filter" "foobar" {
    name = "mypacket_filter_upd"
    description = "PacketFilter from TerraForm for SAKURA CLOUD"
    expressions = {
    	protocol = "tcp"
    	source_nw = "192.168.2.0"
    	source_port = "8080"
    	dest_port = "8080"
    	allow = false
    }
    expressions = {
    	protocol = "udp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
    expressions = {
    	protocol = "icmp"
    	source_nw = "0.0.0.0"
    	allow = true
    }
    expressions = {
    	protocol = "fragment"
    	source_nw = "0.0.0.0"
    	allow = true
    }
    expressions = {
    	protocol = "ip"
    	source_nw = "0.0.0.0"
    	allow = true
    }
}`
