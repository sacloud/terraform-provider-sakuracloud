package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccSakuraCloudVPCRouterInterface_Basic(t *testing.T) {
	var vpcRouter sacloud.VPCRouter
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudVPCRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudVPCRouterInterfaceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "index", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "vip", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "ipaddress.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "ipaddress.0", "192.168.100.1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "nw_mask_len", "24"),
				),
			},
		},
	})
}

func TestAccSakuraCloudVPCRouterInterface_Update(t *testing.T) {
	var vpcRouter sacloud.VPCRouter
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudVPCRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudVPCRouterInterfaceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "index", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "vip", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "ipaddress.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "ipaddress.0", "192.168.100.1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "nw_mask_len", "24"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudVPCRouterInterfaceConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "index", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "vip", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "ipaddress.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "ipaddress.0", "192.168.100.1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "nw_mask_len", "24"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth2", "index", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth2", "vip", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth2", "ipaddress.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth2", "ipaddress.0", "192.168.200.1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth2", "nw_mask_len", "24"),
				),
			},
		},
	})
}

func TestAccSakuraCloudVPCRouterInterface_WithRouter(t *testing.T) {
	var vpcRouter sacloud.VPCRouter
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudVPCRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudVPCRouterInterfaceConfig_WithRouter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "index", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router_interface.eth1", "ipaddress.#", "2"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudVPCRouterInterfaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_vpc_router" {
			continue
		}

		_, err := client.VPCRouter.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("VPCRouter still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudVPCRouterInterfaceConfig_basic = `
resource sakuracloud_switch "sw01"{
    name = "sw01"
}
resource "sakuracloud_vpc_router" "foobar" {
    name = "name"
}
resource "sakuracloud_vpc_router_interface" "eth1"{
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    index = 1
    switch_id = "${sakuracloud_switch.sw01.id}"
    ipaddress = ["192.168.100.1"]
    nw_mask_len = 24
}
`

var testAccCheckSakuraCloudVPCRouterInterfaceConfig_update = `
resource sakuracloud_switch "sw01"{
    name = "sw01"
}
resource sakuracloud_switch "sw02"{
    name = "sw02"
}
resource "sakuracloud_vpc_router" "foobar" {
    name = "name"
}
resource "sakuracloud_vpc_router_interface" "eth1"{
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    index = 1
    switch_id = "${sakuracloud_switch.sw01.id}"
    ipaddress = ["192.168.100.1"]
    nw_mask_len = 24
}
resource "sakuracloud_vpc_router_interface" "eth2"{
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    index = 2
    switch_id = "${sakuracloud_switch.sw02.id}"
    ipaddress = ["192.168.200.1"]
    nw_mask_len = 24
}
`

var testAccCheckSakuraCloudVPCRouterInterfaceConfig_WithRouter = `
resource "sakuracloud_internet" "router1" {
    name = "myinternet1"
}
resource sakuracloud_switch "sw01"{
    name = "sw01"
}
resource "sakuracloud_vpc_router" "foobar" {
    name = "name"
    plan = "premium"
    switch_id = "${sakuracloud_internet.router1.switch_id}"
    vip = "${sakuracloud_internet.router1.nw_ipaddresses.0}"
    ipaddress1 = "${sakuracloud_internet.router1.nw_ipaddresses.1}"
    ipaddress2 = "${sakuracloud_internet.router1.nw_ipaddresses.2}"
    VRID = 1
}
resource "sakuracloud_vpc_router_interface" "eth1"{
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    index = 1
    switch_id = "${sakuracloud_switch.sw01.id}"
    vip = "192.168.100.1"
    ipaddress = ["192.168.100.2","192.168.100.3"]
    nw_mask_len = "24"
}
`
