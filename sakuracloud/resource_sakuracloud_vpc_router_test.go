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
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/sacloud/iaas-api-go"
)

func TestAccSakuraCloudVPCRouter_basic(t *testing.T) {
	resourceName := "sakuracloud_vpc_router.foobar"
	rand := randomName()

	var vpcRouter iaas.VPCRouter
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudIconDestroy,
			testCheckSakuraCloudVPCRouterDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudVPCRouter_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudVPCRouterExists(resourceName, &vpcRouter),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "plan", "standard"),
					resource.TestCheckResourceAttr(resourceName, "version", "2"),
					resource.TestCheckResourceAttr(resourceName, "internet_connection", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "public_network_interface.#"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "public_netmask"),
					resource.TestCheckResourceAttrPair(
						resourceName, "icon_id",
						"sakuracloud_icon.foobar", "id",
					),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudVPCRouter_update, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudVPCRouterExists(resourceName, &vpcRouter),
					resource.TestCheckResourceAttr(resourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "description-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2-upd"),
					resource.TestCheckResourceAttr(resourceName, "plan", "standard"),
					resource.TestCheckResourceAttr(resourceName, "internet_connection", "false"),
					resource.TestCheckNoResourceAttr(resourceName, "public_network_interface.#"),
					resource.TestCheckResourceAttr(resourceName, "syslog_host", "192.168.0.2"),
					resource.TestCheckResourceAttr(resourceName, "icon_id", ""),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "public_netmask"),
				),
			},
		},
	})
}

func TestAccSakuraCloudVPCRouter_Full(t *testing.T) {
	resourceName := "sakuracloud_vpc_router.foobar"
	rand := randomName()

	var vpcRouter iaas.VPCRouter
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudInternetDestroy,
			testCheckSakuraCloudSwitchDestroy,
			testCheckSakuraCloudVPCRouterDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudVPCRouter_complete, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudVPCRouterExists(resourceName, &vpcRouter),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "version", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "public_netmask"),
					resource.TestCheckResourceAttr(resourceName, "private_network_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "private_network_interface.0.vip", "192.168.11.1"),
					resource.TestCheckResourceAttr(resourceName, "private_network_interface.0.ip_addresses.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "private_network_interface.0.ip_addresses.0", "192.168.11.2"),
					resource.TestCheckResourceAttr(resourceName, "private_network_interface.0.ip_addresses.1", "192.168.11.3"),
					resource.TestCheckResourceAttr(resourceName, "private_network_interface.0.netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_server.0.interface_index", "1"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_server.0.range_start", "192.168.11.11"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_server.0.range_stop", "192.168.11.20"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_server.0.dns_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_server.0.dns_servers.0", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_server.0.dns_servers.1", "8.8.4.4"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_static_mapping.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_static_mapping.0.ip_address", "192.168.11.10"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_static_mapping.0.mac_address", "aa:bb:cc:aa:bb:cc"),
					resource.TestCheckResourceAttr(resourceName, "dns_forwarding.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_forwarding.0.interface_index", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_forwarding.0.dns_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_forwarding.0.dns_servers.0", "133.242.0.3"),
					resource.TestCheckResourceAttr(resourceName, "dns_forwarding.0.dns_servers.1", "133.242.0.4"),
					resource.TestCheckResourceAttr(resourceName, "firewall.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.interface_index", "1"),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.direction", "send"),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.0.allow", "true"),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.0.source_network", ""),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.0.source_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.0.destination_network", ""),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.0.destination_port", ""),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.1.protocol", "ip"),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.1.allow", "false"),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.1.source_network", ""),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.1.source_port", ""),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.1.destination_network", ""),
					resource.TestCheckResourceAttr(resourceName, "firewall.0.expression.1.destination_port", ""),
					resource.TestCheckResourceAttr(resourceName, "l2tp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "l2tp.0.pre_shared_secret", "example"),
					resource.TestCheckResourceAttr(resourceName, "l2tp.0.range_start", "192.168.11.21"),
					resource.TestCheckResourceAttr(resourceName, "l2tp.0.range_stop", "192.168.11.30"),
					resource.TestCheckResourceAttr(resourceName, "port_forwarding.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "port_forwarding.0.protocol", "udp"),
					resource.TestCheckResourceAttr(resourceName, "port_forwarding.0.public_port", "10022"),
					resource.TestCheckResourceAttr(resourceName, "port_forwarding.0.private_ip", "192.168.11.11"),
					resource.TestCheckResourceAttr(resourceName, "port_forwarding.0.private_port", "22"),
					resource.TestCheckResourceAttr(resourceName, "port_forwarding.0.description", "desc"),
					resource.TestCheckResourceAttr(resourceName, "pptp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "pptp.0.range_start", "192.168.11.31"),
					resource.TestCheckResourceAttr(resourceName, "pptp.0.range_stop", "192.168.11.40"),
					resource.TestCheckResourceAttr(resourceName, "wire_guard.0.ip_address", "192.168.31.1/24"),
					resource.TestCheckResourceAttrSet(resourceName, "wire_guard.0.public_key"),
					resource.TestCheckResourceAttr(resourceName, "wire_guard.0.peer.0.name", "example"),
					resource.TestCheckResourceAttr(resourceName, "wire_guard.0.peer.0.ip_address", "192.168.31.11"),
					resource.TestCheckResourceAttr(resourceName, "wire_guard.0.peer.0.public_key", "fqxOlS2X0Jtg4P9zVf8D3BAUtJmrp+z2mjzUmgxxxxx="),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn.0.peer", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn.0.remote_id", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn.0.pre_shared_secret", "example"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn.0.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn.0.routes.0", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn.0.local_prefix.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn.0.local_prefix.0", "192.168.21.0/24"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn_parameter.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn_parameter.0.ike.0.lifetime", "28801"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn_parameter.0.ike.0.dpd.0.interval", "16"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn_parameter.0.ike.0.dpd.0.timeout", "31"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn_parameter.0.esp.0.lifetime", "1801"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn_parameter.0.encryption_algo", "aes256"),
					resource.TestCheckResourceAttr(resourceName, "site_to_site_vpn_parameter.0.hash_algo", "sha256"),
					resource.TestCheckResourceAttr(resourceName, "static_nat.#", "1"),
					resource.TestCheckResourceAttrPair(
						resourceName, "static_nat.0.public_ip",
						"sakuracloud_internet.foobar", "ip_addresses.3"),
					resource.TestCheckResourceAttr(resourceName, "static_nat.0.private_ip", "192.168.11.12"),
					resource.TestCheckResourceAttr(resourceName, "static_nat.0.description", "desc"),
					resource.TestCheckResourceAttr(resourceName, "static_route.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "static_route.0.prefix", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "static_route.0.next_hop", "192.168.11.99"),
					resource.TestCheckResourceAttr(resourceName, "user.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user.0.name", "username"),
					resource.TestCheckResourceAttr(resourceName, "user.0.password", "password"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudVPCRouter_completeUpdate, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudVPCRouterExists(resourceName, &vpcRouter),
					resource.TestCheckResourceAttr(resourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "public_netmask"),
					resource.TestCheckNoResourceAttr(resourceName, "private_network_interface.#"),
					resource.TestCheckNoResourceAttr(resourceName, "dhcp_server.#"),
					resource.TestCheckNoResourceAttr(resourceName, "dhcp_static_mapping.#"),
					resource.TestCheckNoResourceAttr(resourceName, "firewall.#"),
					resource.TestCheckNoResourceAttr(resourceName, "l2tp.#"),
					resource.TestCheckNoResourceAttr(resourceName, "port_forwarding.#"),
					resource.TestCheckNoResourceAttr(resourceName, "pptp.#"),
					resource.TestCheckNoResourceAttr(resourceName, "site_to_site_vpn.#"),
					resource.TestCheckNoResourceAttr(resourceName, "static_nat.#"),
					resource.TestCheckNoResourceAttr(resourceName, "static_route.#"),
					resource.TestCheckNoResourceAttr(resourceName, "user.#"),
				),
			},
		},
	})
}

func testCheckSakuraCloudVPCRouterExists(n string, vpcRouter *iaas.VPCRouter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no VPCRouter ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		vrOp := iaas.NewVPCRouterOp(client)
		zone := rs.Primary.Attributes["zone"]

		foundVPCRouter, err := vrOp.Read(context.Background(), zone, sakuraCloudID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if foundVPCRouter.ID.String() != rs.Primary.ID {
			return fmt.Errorf("not found VPCRouter: %s", rs.Primary.ID)
		}

		*vpcRouter = *foundVPCRouter

		return nil
	}
}

func testCheckSakuraCloudVPCRouterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	vrOp := iaas.NewVPCRouterOp(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_vpc_router" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		zone := rs.Primary.Attributes["zone"]
		_, err := vrOp.Read(context.Background(), zone, sakuraCloudID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("still exists VPCRouter: %s", rs.Primary.ID)
		}
	}

	return nil
}

var testAccSakuraCloudVPCRouter_basic = `
resource "sakuracloud_vpc_router" "foobar" {
  name                = "{{ .arg0 }}"
  description         = "description"
  tags                = ["tag1", "tag2"]
  icon_id             = sakuracloud_icon.foobar.id
  internet_connection = true
}

resource "sakuracloud_icon" "foobar" {
  name          = "{{ .arg0 }}"
  base64content = "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"
}
`

var testAccSakuraCloudVPCRouter_update = `
resource "sakuracloud_vpc_router" "foobar" {
  name                = "{{ .arg0 }}-upd"
  description         = "description-upd"
  tags                = ["tag1-upd", "tag2-upd"]
  syslog_host         = "192.168.0.2"
  internet_connection = false
}`

var testAccSakuraCloudVPCRouter_complete = `
resource "sakuracloud_internet" "foobar" {
  name = "{{ .arg0 }}"
}
resource sakuracloud_switch "foobar" {
  name = "{{ .arg0 }}"
}

resource "sakuracloud_vpc_router" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1" , "tag2"]
  plan        = "premium"
  version     = 2

  internet_connection = true

  public_network_interface {
    switch_id    = sakuracloud_internet.foobar.switch_id
    vip          = sakuracloud_internet.foobar.ip_addresses[0]
    ip_addresses = [sakuracloud_internet.foobar.ip_addresses[1], sakuracloud_internet.foobar.ip_addresses[2]]
    aliases      = [sakuracloud_internet.foobar.ip_addresses[3]]
    vrid         = 1
  }

  private_network_interface {
    index        = 1
    switch_id    = sakuracloud_switch.foobar.id
    vip          = "192.168.11.1"
    ip_addresses = ["192.168.11.2", "192.168.11.3"]
    netmask      = 24 
  }

  dhcp_server {
    interface_index = 1

    range_start = "192.168.11.11"
    range_stop  = "192.168.11.20"
    dns_servers = ["8.8.8.8", "8.8.4.4"]
  }

  dhcp_static_mapping {
    ip_address  = "192.168.11.10"
    mac_address = "aa:bb:cc:aa:bb:cc"
  }

  dns_forwarding {
    interface_index = 1
    dns_servers = ["133.242.0.3", "133.242.0.4"]
  }


  firewall {
    interface_index = 1

    direction = "send"
    expression {
        protocol            = "tcp"
        source_network      = ""
        source_port         = "80"
        destination_network = ""
        destination_port    = ""
        allow               = true
        logging             = true
        description         = "desc"
    }

    expression {
        protocol            = "ip"
        source_network      = ""
        source_port         = ""
        destination_network = ""
        destination_port    = ""
        allow               = false
        logging             = true
        description         = "desc"
    }
  }

  l2tp {
    pre_shared_secret = "example"
    range_start       = "192.168.11.21"
    range_stop        = "192.168.11.30"
  }

  port_forwarding {
    protocol     = "udp"
    public_port  = 10022
    private_ip   = "192.168.11.11"
    private_port = 22
    description  = "desc"
  }

  pptp {
    range_start = "192.168.11.31"
    range_stop  = "192.168.11.40"
  }

  wire_guard {
    ip_address = "192.168.31.1/24"
    peer {
      name       = "example"
      ip_address = "192.168.31.11"
      public_key = "fqxOlS2X0Jtg4P9zVf8D3BAUtJmrp+z2mjzUmgxxxxx="
    }
  }

  site_to_site_vpn {
    peer              = "8.8.8.8"
    remote_id         = "8.8.8.8"
    pre_shared_secret = "example"
    routes            = ["10.0.0.0/8"]
    local_prefix      = ["192.168.21.0/24"]
  }

  site_to_site_vpn_parameter {
    ike {
      lifetime = 28801
      dpd {
        interval = 16
        timeout  = 31
      }
    }
    esp {
      lifetime = 1801
    }
    encryption_algo = "aes256"
    hash_algo       = "sha256"
  }

  static_nat {
    public_ip   = sakuracloud_internet.foobar.ip_addresses[3]
    private_ip  = "192.168.11.12"
    description = "desc"
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

var testAccSakuraCloudVPCRouter_completeUpdate = `
resource "sakuracloud_internet" "foobar" {
  name = "{{ .arg0 }}"
}
resource sakuracloud_switch "foobar" {
  name = "{{ .arg0 }}"
}
resource "sakuracloud_vpc_router" "foobar" {
  name        = "{{ .arg0 }}-upd"
  description = "description-upd"
  tags        = ["tag1-upd" , "tag2-upd"]
  plan        = "premium"

  internet_connection = true

  public_network_interface {
    switch_id    = sakuracloud_internet.foobar.switch_id
    vip          = sakuracloud_internet.foobar.ip_addresses[0]
    ip_addresses = [sakuracloud_internet.foobar.ip_addresses[1], sakuracloud_internet.foobar.ip_addresses[2]]
    aliases      = [sakuracloud_internet.foobar.ip_addresses[3]]
    vrid         = 1
  }
}
`
