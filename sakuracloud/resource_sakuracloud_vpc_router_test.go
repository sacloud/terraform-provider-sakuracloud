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
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/sacloud"
)

func TestAccResourceSakuraCloudVPCRouter_basic(t *testing.T) {
	var vpcRouter sacloud.VPCRouter
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudVPCRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudVPCRouterConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "name", "name_before"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "description", "description_before"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "tags.0", "hoge1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "tags.1", "hoge2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "plan", "standard"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "internet_connection", "true"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "switch_id"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "vip"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "ipaddress1"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "ipaddress2"),
					resource.TestCheckResourceAttrPair(
						"sakuracloud_vpc_router.foobar", "icon_id",
						"sakuracloud_icon.foobar", "id",
					),
				),
			},
			{
				Config: testAccCheckSakuraCloudVPCRouterConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "name", "name_after"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "description", "description_after"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "tags.0", "hoge1_after"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "tags.1", "hoge2_after"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "plan", "standard"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "internet_connection", "false"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "switch_id"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "vip"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "ipaddress1"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "ipaddress2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "syslog_host", "192.168.0.2"),
					resource.TestCheckNoResourceAttr("sakuracloud_vpc_router.foobar", "icon_id"),
				),
			},
		},
	})
}

func TestAccResourceSakuraCloudVPCRouter_full(t *testing.T) {
	var vpcRouter sacloud.VPCRouter
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudVPCRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudVPCRouterConfig_full_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "name", "name_before"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "interface.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "interface.0.vip", "192.168.11.1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "interface.0.ipaddress.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "interface.0.ipaddress.0", "192.168.11.2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "interface.0.ipaddress.1", "192.168.11.3"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "interface.0.nw_mask_len", "24"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "dhcp_server.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "dhcp_server.0.vpc_router_interface_index", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "dhcp_server.0.range_start", "192.168.11.11"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "dhcp_server.0.range_stop", "192.168.11.20"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "dhcp_server.0.dns_servers.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "dhcp_server.0.dns_servers.0", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "dhcp_server.0.dns_servers.1", "8.8.4.4"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "dhcp_static_mapping.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "dhcp_static_mapping.0.ipaddress", "192.168.11.10"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "dhcp_static_mapping.0.macaddress", "aa:bb:cc:aa:bb:cc"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.vpc_router_interface_index", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.direction", "send"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.0.allow", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.0.source_nw", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.0.source_port", "80"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.0.dest_nw", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.0.dest_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.1.protocol", "ip"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.1.allow", "false"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.1.source_nw", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.1.source_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.1.dest_nw", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall.0.expressions.1.dest_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "l2tp.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "l2tp.0.pre_shared_secret", "example"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "l2tp.0.range_start", "192.168.11.21"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "l2tp.0.range_stop", "192.168.11.30"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "port_forwarding.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "port_forwarding.0.protocol", "udp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "port_forwarding.0.global_port", "10022"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "port_forwarding.0.private_address", "192.168.11.11"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "port_forwarding.0.private_port", "22"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "port_forwarding.0.description", "desc"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "pptp.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "pptp.0.range_start", "192.168.11.31"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "pptp.0.range_stop", "192.168.11.40"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "site_to_site_vpn.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "site_to_site_vpn.0.peer", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "site_to_site_vpn.0.remote_id", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "site_to_site_vpn.0.pre_shared_secret", "example"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "site_to_site_vpn.0.routes.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "site_to_site_vpn.0.routes.0", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "site_to_site_vpn.0.local_prefix.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "site_to_site_vpn.0.local_prefix.0", "192.168.21.0/24"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "static_nat.#", "1"),
					resource.TestCheckResourceAttrPair(
						"sakuracloud_vpc_router.foobar", "static_nat.0.global_address",
						"sakuracloud_internet.router1", "ipaddresses.3"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "static_nat.0.private_address", "192.168.11.12"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "static_nat.0.description", "desc"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "static_route.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "static_route.0.prefix", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "static_route.0.next_hop", "192.168.11.99"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "user.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "user.0.name", "username"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "user.0.password", "password"),
				),
			},
			{
				Config: testAccCheckSakuraCloudVPCRouterConfig_full_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "name", "name_before"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "interface"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "dhcp_server"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "dhcp_static_mapping"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "firewall"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "l2tp"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "port_forwarding"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "pptp"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "site_to_site_vpn"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "static_nat"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "static_route"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "user"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudVPCRouterExists(n string, vpcRouter *sacloud.VPCRouter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No VPCRouter ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)

		foundVPCRouter, err := client.VPCRouter.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundVPCRouter.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("VPCRouter not found")
		}

		*vpcRouter = *foundVPCRouter

		return nil
	}
}

func testAccCheckSakuraCloudVPCRouterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_vpc_router" {
			continue
		}

		_, err := client.VPCRouter.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("VPCRouter still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudVPCRouterConfig_basic = `
resource "sakuracloud_vpc_router" "foobar" {
    name = "name_before"
    description = "description_before"
    tags = ["hoge1" , "hoge2"]
    icon_id = "${sakuracloud_icon.foobar.id}"
    internet_connection = true
}

resource "sakuracloud_icon" "foobar" {
  name = "myicon"
  base64content = "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"
}
`

var testAccCheckSakuraCloudVPCRouterConfig_update = `
resource "sakuracloud_vpc_router" "foobar" {
    name = "name_after"
    description = "description_after"
    tags = ["hoge1_after" , "hoge2_after"]
    syslog_host = "192.168.0.2"
    internet_connection = false
}`

var testAccCheckSakuraCloudVPCRouterConfig_full_basic = `
resource "sakuracloud_internet" "router1" {
    name = "myinternet1"
}
resource sakuracloud_switch "sw" {
  name = "name_before"
}

resource "sakuracloud_vpc_router" "foobar" {
  name        = "name_before"
  description = "description_before"
  tags        = ["hoge1" , "hoge2"]
  plan        = "premium"

  internet_connection = true

  switch_id  = "${sakuracloud_internet.router1.switch_id}"
  vip        = "${sakuracloud_internet.router1.ipaddresses[0]}"
  ipaddress1 = "${sakuracloud_internet.router1.ipaddresses[1]}"
  ipaddress2 = "${sakuracloud_internet.router1.ipaddresses[2]}"
  aliases    = ["${sakuracloud_internet.router1.ipaddresses[3]}"]
  vrid       = 1

  interface {
    switch_id   = "${sakuracloud_switch.sw.id}"
    vip         = "192.168.11.1"
    ipaddress   = ["192.168.11.2" , "192.168.11.3"]
    nw_mask_len = 24 
  }

  dhcp_server {
    vpc_router_interface_index = 1

    range_start = "192.168.11.11"
    range_stop  = "192.168.11.20"
    dns_servers = ["8.8.8.8", "8.8.4.4"]
  }

  dhcp_static_mapping {
    ipaddress  = "192.168.11.10"
    macaddress = "aa:bb:cc:aa:bb:cc"
  }

  firewall {
    vpc_router_interface_index = 1

    direction = "send"
    expressions {
        protocol    = "tcp"
        source_nw   = ""
        source_port = "80"
        dest_nw     = ""
        dest_port   = ""
        allow       = true
        logging     = true
        description = "desc"
    }

    expressions {
        protocol    = "ip"
        source_nw   = ""
        source_port = ""
        dest_nw     = ""
        dest_port   = ""
        allow       = false
        logging     = true
        description = "desc"
    }
  }

  l2tp {
    pre_shared_secret = "example"
    range_start       = "192.168.11.21"
    range_stop        = "192.168.11.30"
  }

  port_forwarding {
    protocol        = "udp"
    global_port     = 10022
    private_address = "192.168.11.11"
    private_port    = 22
    description     = "desc"
  }

  pptp {
    range_start = "192.168.11.31"
    range_stop  = "192.168.11.40"
  }

  site_to_site_vpn {
    peer              = "8.8.8.8"
    remote_id         = "8.8.8.8"
    pre_shared_secret = "example"
    routes            = ["10.0.0.0/8"]
    local_prefix      = ["192.168.21.0/24"]
  }

  static_nat {
    global_address  = "${sakuracloud_internet.router1.ipaddresses[3]}"
    private_address = "192.168.11.12"
    description     = "desc"
  }

  static_route {
    prefix   = "172.16.0.0/16"
    next_hop = "192.168.11.99"
  }

  user {
    name     = "username"
    password = "password"
  }
}
`

var testAccCheckSakuraCloudVPCRouterConfig_full_update = `
resource "sakuracloud_internet" "router1" {
    name = "myinternet1"
}
resource sakuracloud_switch "sw" {
  name = "name_before"
}
resource "sakuracloud_vpc_router" "foobar" {
  depends_on = ["sakuracloud_switch.sw", "sakuracloud_internet.router1"] # TODO for terraform v0.12-alpha4(if without this, deleting switch will fail)

  name        = "name_before"
  description = "description_before"
  tags        = ["hoge1" , "hoge2"]
  plan        = "premium"

  internet_connection = true

  switch_id  = "${sakuracloud_internet.router1.switch_id}"
  vip        = "${sakuracloud_internet.router1.ipaddresses[0]}"
  ipaddress1 = "${sakuracloud_internet.router1.ipaddresses[1]}"
  ipaddress2 = "${sakuracloud_internet.router1.ipaddresses[2]}"
  aliases    = [ "${sakuracloud_internet.router1.ipaddresses[3]}" ]
  vrid       = 1
}
`
