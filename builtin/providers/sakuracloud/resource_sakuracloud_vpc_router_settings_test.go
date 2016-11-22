package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"testing"
)

func TestAccSakuraCloudVPCRouterSetting_Basic(t *testing.T) {
	var vpcRouter sacloud.VPCRouter
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudVPCRouterSettingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudVPCRouterSettingConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_pptp.pptp", "range_start", "192.168.11.101"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_pptp.pptp", "range_stop", "192.168.11.150"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_static_nat.staticNAT1", "private_address", "192.168.11.11"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_port_forwarding.forward1", "protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_port_forwarding.forward1", "global_port", "10022"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_port_forwarding.forward1", "private_address", "192.168.11.11"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_port_forwarding.forward1", "private_port", "22"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_dhcp_server.dhcp", "range_start", "192.168.11.151"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_dhcp_server.dhcp", "range_stop", "192.168.11.200"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_dhcp_static_mapping.dhcp_map", "ipaddress", "192.168.11.20"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_dhcp_static_mapping.dhcp_map", "macaddress", "aa:bb:cc:aa:bb:cc"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_l2tp.l2tp", "pre_shared_secret", "hogehoge"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_l2tp.l2tp", "range_start", "192.168.11.51"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_l2tp.l2tp", "range_stop", "192.168.11.100"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_user.user1", "name", "username"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_user.user1", "password", "password"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "peer", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "remote_id", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "pre_shared_secret", "presharedsecret"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "routes.0", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "local_prefix.0", "192.168.21.0/24"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "direction", "send"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "expressions.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "expressions.0.allow", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "expressions.0.source_nw", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "expressions.0.source_port", "80"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "expressions.0.dest_nw", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "expressions.0.dest_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "expressions.1.protocol", "ip"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "expressions.1.allow", "false"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "expressions.1.source_nw", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "expressions.1.source_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "expressions.1.dest_nw", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.send_fw", "expressions.1.dest_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "direction", "receive"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "expressions.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "expressions.0.allow", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "expressions.0.source_nw", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "expressions.0.source_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "expressions.0.dest_nw", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "expressions.0.dest_port", "22"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "expressions.1.protocol", "ip"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "expressions.1.allow", "false"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "expressions.1.source_nw", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "expressions.1.source_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "expressions.1.dest_nw", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_firewall.receive_fw", "expressions.1.dest_port", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_static_route.route", "prefix", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_static_route.route", "next_hop", "192.168.11.99"),
				),
			},
		},
	})
}

func TestAccSakuraCloudVPCRouterSetting_Update(t *testing.T) {
	var vpcRouter sacloud.VPCRouter
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudVPCRouterSettingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudVPCRouterSettingConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_pptp.pptp", "range_start", "192.168.11.101"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_pptp.pptp", "range_stop", "192.168.11.150"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_static_nat.staticNAT1", "private_address", "192.168.11.11"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_port_forwarding.forward1", "protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_port_forwarding.forward1", "global_port", "10022"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_port_forwarding.forward1", "private_address", "192.168.11.11"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_port_forwarding.forward1", "private_port", "22"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_dhcp_server.dhcp", "range_start", "192.168.11.151"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_dhcp_server.dhcp", "range_stop", "192.168.11.200"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_dhcp_static_mapping.dhcp_map", "ipaddress", "192.168.11.20"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_dhcp_static_mapping.dhcp_map", "macaddress", "aa:bb:cc:aa:bb:cc"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_l2tp.l2tp", "pre_shared_secret", "hogehoge"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_l2tp.l2tp", "range_start", "192.168.11.51"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_l2tp.l2tp", "range_stop", "192.168.11.100"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_user.user1", "name", "username"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_user.user1", "password", "password"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "peer", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "remote_id", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "pre_shared_secret", "presharedsecret"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "routes.0", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "local_prefix.0", "192.168.21.0/24"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudVPCRouterSettingConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_pptp.pptp", "range_start", "192.168.11.201"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_pptp.pptp", "range_stop", "192.168.11.250"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_static_nat.staticNAT1", "private_address", "192.168.11.12"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_port_forwarding.forward1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_port_forwarding.forward1", "global_port", "10022"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_port_forwarding.forward1", "private_address", "192.168.11.11"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_port_forwarding.forward1", "private_port", "22"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_dhcp_server.dhcp", "range_start", "192.168.11.151"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_dhcp_server.dhcp", "range_stop", "192.168.11.200"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_dhcp_static_mapping.dhcp_map", "ipaddress", "192.168.11.21"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_dhcp_static_mapping.dhcp_map", "macaddress", "aa:bb:cc:aa:bb:cc"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_l2tp.l2tp", "pre_shared_secret", "hogehoge"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_l2tp.l2tp", "range_start", "192.168.11.51"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_l2tp.l2tp", "range_stop", "192.168.11.100"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_user.user1", "name", "username"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_user.user1", "password", "password"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "peer", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "remote_id", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "pre_shared_secret", "presharedsecret"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "routes.0", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_site_to_site_vpn.s2s", "local_prefix.0", "192.168.21.0/24"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudVPCRouterSettingDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_vpc_router" {
			continue
		}

		_, err := client.VPCRouter.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("VPCRouter still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudVPCRouterSettingConfig_basic = `
resource "sakuracloud_internet" "router1" {
    name = "myinternet1"
}
resource "sakuracloud_switch" "sw01"{
    name = "sw01"
}
resource "sakuracloud_vpc_router" "foobar" {
    name = "vpc_router_setting_test"
    plan = "premium"
    switch_id = "${sakuracloud_internet.router1.switch_id}"
    vip = "${sakuracloud_internet.router1.nw_ipaddresses.0}"
    ipaddress1 = "${sakuracloud_internet.router1.nw_ipaddresses.1}"
    ipaddress2 = "${sakuracloud_internet.router1.nw_ipaddresses.2}"
    aliases = ["${sakuracloud_internet.router1.nw_ipaddresses.3}"]
    VRID = 1

}
resource "sakuracloud_vpc_router_interface" "eth1"{
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    index = 1
    switch_id = "${sakuracloud_switch.sw01.id}"
    vip = "192.168.11.1"
    ipaddress = ["192.168.11.2" , "192.168.11.3"]
    nw_mask_len = 24
}
resource "sakuracloud_vpc_router_pptp" "pptp"{
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"

    range_start = "192.168.11.101"
    range_stop = "192.168.11.150"
}
resource "sakuracloud_vpc_router_static_nat" "staticNAT1" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"

    global_address = "${sakuracloud_internet.router1.nw_ipaddresses.3}"
    private_address = "192.168.11.11"
}
resource "sakuracloud_vpc_router_port_forwarding" "forward1" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"

    protocol = "tcp"
    global_port = 10022
    private_address = "192.168.11.11"
    private_port = 22
}
resource "sakuracloud_vpc_router_dhcp_server" "dhcp" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_index = "${sakuracloud_vpc_router_interface.eth1.index}"

    range_start = "192.168.11.151"
    range_stop = "192.168.11.200"
}
resource "sakuracloud_vpc_router_dhcp_static_mapping" "dhcp_map" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_dhcp_server_id = "${sakuracloud_vpc_router_dhcp_server.dhcp.id}"

    ipaddress = "192.168.11.20"
    macaddress = "aa:bb:cc:aa:bb:cc"
}
resource "sakuracloud_vpc_router_l2tp" "l2tp" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"

    pre_shared_secret = "hogehoge"
    range_start = "192.168.11.51"
    range_stop = "192.168.11.100"

}
resource "sakuracloud_vpc_router_user" "user1" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    name = "username"
    password = "password"
}
resource "sakuracloud_vpc_router_site_to_site_vpn" "s2s" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    peer = "8.8.8.8"
    remote_id = "8.8.8.8"
    pre_shared_secret = "presharedsecret"
    routes = ["10.0.0.0/8"]
    local_prefix = ["192.168.21.0/24"]
}

resource "sakuracloud_vpc_router_firewall" "send_fw" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    direction = "send"
    expressions = {
        protocol = "tcp"
        source_nw = ""
        source_port = "80"
        dest_nw = ""
        dest_port = ""
        allow = true
    }

    expressions = {
        protocol = "ip"
        source_nw = ""
        source_port = ""
        dest_nw = ""
        dest_port = ""
        allow = false
    }
}

resource "sakuracloud_vpc_router_firewall" "receive_fw" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    direction = "receive"
    expressions = {
        protocol = "tcp"
        source_nw = ""
        source_port = ""
        dest_nw = ""
        dest_port = "22"
        allow = true
    }

    expressions = {
        protocol = "ip"
        source_nw = ""
        source_port = ""
        dest_nw = ""
        dest_port = ""
        allow = false
    }
}
resource "sakuracloud_vpc_router_static_route" "route" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"

    prefix = "172.16.0.0/16"
    next_hop = "192.168.11.99"
}
`

var testAccCheckSakuraCloudVPCRouterSettingConfig_update = `
resource "sakuracloud_internet" "router1" {
    name = "myinternet1"
}
resource "sakuracloud_switch" "sw01"{
    name = "sw01"
}
resource "sakuracloud_vpc_router" "foobar" {
    name = "vpc_router_setting_test"
    plan = "premium"
    switch_id = "${sakuracloud_internet.router1.switch_id}"
    vip = "${sakuracloud_internet.router1.nw_ipaddresses.0}"
    ipaddress1 = "${sakuracloud_internet.router1.nw_ipaddresses.1}"
    ipaddress2 = "${sakuracloud_internet.router1.nw_ipaddresses.2}"
    aliases = ["${sakuracloud_internet.router1.nw_ipaddresses.3}"]
    VRID = 1

}
resource "sakuracloud_vpc_router_interface" "eth1"{
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    index = 1
    switch_id = "${sakuracloud_switch.sw01.id}"
    vip = "192.168.11.1"
    ipaddress = ["192.168.11.2" , "192.168.11.3"]
    nw_mask_len = 24
}
resource "sakuracloud_vpc_router_pptp" "pptp"{
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"

    range_start = "192.168.11.201"
    range_stop = "192.168.11.250"
}
resource "sakuracloud_vpc_router_static_nat" "staticNAT1" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"

    global_address = "${sakuracloud_internet.router1.nw_ipaddresses.3}"
    private_address = "192.168.11.12"
}
resource "sakuracloud_vpc_router_port_forwarding" "forward1" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"

    protocol = "udp"
    global_port = 10022
    private_address = "192.168.11.11"
    private_port = 22
}
resource "sakuracloud_vpc_router_dhcp_server" "dhcp" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_index = "${sakuracloud_vpc_router_interface.eth1.index}"

    range_start = "192.168.11.151"
    range_stop = "192.168.11.200"
}
resource "sakuracloud_vpc_router_dhcp_static_mapping" "dhcp_map" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_dhcp_server_id = "${sakuracloud_vpc_router_dhcp_server.dhcp.id}"

    ipaddress = "192.168.11.21"
    macaddress = "aa:bb:cc:aa:bb:cc"
}
resource "sakuracloud_vpc_router_l2tp" "l2tp" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"

    pre_shared_secret = "hogehoge"
    range_start = "192.168.11.51"
    range_stop = "192.168.11.100"

}
resource "sakuracloud_vpc_router_user" "user1" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    name = "username"
    password = "password"
}
resource "sakuracloud_vpc_router_site_to_site_vpn" "s2s" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    peer = "8.8.8.8"
    remote_id = "8.8.8.8"
    pre_shared_secret = "presharedsecret"
    routes = ["10.0.0.0/8"]
    local_prefix = ["192.168.21.0/24"]
}
resource "sakuracloud_vpc_router_firewall" "send_fw" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    direction = "send"
    expressions = {
        protocol = "tcp"
        source_nw = ""
        source_port = "80"
        dest_nw = ""
        dest_port = ""
        allow = true
    }

    expressions = {
        protocol = "ip"
        source_nw = ""
        source_port = ""
        dest_nw = ""
        dest_port = ""
        allow = false
    }
}

resource "sakuracloud_vpc_router_firewall" "receive_fw" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    direction = "receive"
    expressions = {
        protocol = "tcp"
        source_nw = ""
        source_port = ""
        dest_nw = ""
        dest_port = "22"
        allow = true
    }

    expressions = {
        protocol = "ip"
        source_nw = ""
        source_port = ""
        dest_nw = ""
        dest_port = ""
        allow = false
    }
}
resource "sakuracloud_vpc_router_static_route" "route" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"

    prefix = "172.16.0.0/16"
    next_hop = "192.168.11.99"
}

`
