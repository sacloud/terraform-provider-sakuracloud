package sakuracloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func TestAccResourceSakuraCloudPacketFilterRule(t *testing.T) {
	var filter sacloud.PacketFilter
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudPacketFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudPacketFilterRuleConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPacketFilterExists("sakuracloud_packet_filter.foobar", &filter),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.0.source_network", "192.168.2.0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.0.source_port", "80"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.0.destination_port", "80"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.0.allow", "true"),

					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.1.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.1.source_network", "192.168.2.0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.1.source_port", "443"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.1.destination_port", "443"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.1.allow", "true"),

					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.2.protocol", "ip"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.2.source_network", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.2.source_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.2.destination_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.2.allow", "false"),
				),
			},
			{
				Config: testAccCheckSakuraCloudPacketFilterRuleConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.0.protocol", "udp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.0.source_network", "192.168.2.2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.0.source_port", "80"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.0.destination_port", "80"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.0.allow", "true"),

					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.1.protocol", "udp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.1.source_network", "192.168.2.2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.1.source_port", "443"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.1.destination_port", "443"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.1.allow", "true"),

					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.2.protocol", "ip"),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.2.source_network", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.2.source_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.2.destination_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_packet_filter_rule.rules", "expressions.2.allow", "false"),
				),
			},
		},
	})
}

var testAccCheckSakuraCloudPacketFilterRuleConfig_basic = `
resource "sakuracloud_packet_filter" "foobar" {
  name        = "mypacket_filter"
  description = "PacketFilter from TerraForm for SAKURA CLOUD"
}

resource sakuracloud_packet_filter_rule "rules" {
  packet_filter_id = "${sakuracloud_packet_filter.foobar.id}"
  expressions {
 	protocol         = "tcp"
	source_network   = "192.168.2.0"
	source_port      = "80"
	destination_port = "80"
	allow            = true
  }
  expressions {
	protocol         = "tcp"
	source_network   = "192.168.2.0"
	source_port      = "443"
	destination_port = "443"
	allow            = true
  }
  expressions {
 	protocol = "ip"
	allow    = false
  }
}
`

var testAccCheckSakuraCloudPacketFilterRuleConfig_update = `
resource "sakuracloud_packet_filter" "foobar" {
  name = "mypacket_filter"
  description = "PacketFilter from TerraForm for SAKURA CLOUD"
}

resource sakuracloud_packet_filter_rule "rules" {
  packet_filter_id = "${sakuracloud_packet_filter.foobar.id}"
  expressions {
   	protocol         = "udp"
  	source_network   = "192.168.2.2"
  	source_port      = "80"
  	destination_port = "80"
   	allow            = true
  }
  expressions {
   	protocol         = "udp"
  	source_network   = "192.168.2.2"
  	source_port      = "443"
  	destination_port = "443"
  	allow            = true
  }
  expressions {
  	protocol = "ip"
	allow    = false
  }
}
`
