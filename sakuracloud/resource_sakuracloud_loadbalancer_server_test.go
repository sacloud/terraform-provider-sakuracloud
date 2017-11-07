package sakuracloud

import (
	"errors"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/sacloud/libsacloud/sacloud"
	"testing"
)

func TestAccResourceSakuraCloudLoadBalancerServer(t *testing.T) {
	var loadBalancer sacloud.LoadBalancer
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudLoadBalancerServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudLoadBalancerServerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerExists("sakuracloud_load_balancer.foobar", &loadBalancer),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_server.server01", "ipaddress", "192.168.11.51"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_server.server01", "check_protocol", "http"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_server.server01", "check_path", "/"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_server.server01", "check_status", "200"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_server.server01", "enabled", "true"),
				),
			},
			{
				Config: testAccCheckSakuraCloudLoadBalancerServerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerExists("sakuracloud_load_balancer.foobar", &loadBalancer),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_server.server01", "ipaddress", "192.168.11.51"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_server.server01", "check_protocol", "ping"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_server.server01", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_server.server02", "ipaddress", "192.168.11.52"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_server.server02", "check_protocol", "ping"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_server.server02", "enabled", "true"),
				),
			},
			{
				Config: testAccCheckSakuraCloudLoadBalancerServerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "servers.#", "2"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudLoadBalancerServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_load_balancer" {
			continue
		}

		_, err := client.LoadBalancer.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("LoadBalancer still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudLoadBalancerServerConfig_basic = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    vrid = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    name = "name"
}
resource "sakuracloud_load_balancer_vip" "vip1" {
    load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
    vip = "192.168.11.201"
    port = 80
}
resource "sakuracloud_load_balancer_server" "server01"{
    load_balancer_vip_id = "${sakuracloud_load_balancer_vip.vip1.id}"
    ipaddress = "192.168.11.51"
    check_protocol = "http"
    check_path = "/"
    check_status = "200"
}
`

var testAccCheckSakuraCloudLoadBalancerServerConfig_update = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    vrid = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    name = "name"
}
resource "sakuracloud_load_balancer_vip" "vip1" {
    load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
    vip = "192.168.11.201"
    port = 80
}
resource "sakuracloud_load_balancer_server" "server01"{
    load_balancer_vip_id = "${sakuracloud_load_balancer_vip.vip1.id}"
    ipaddress = "192.168.11.51"
    check_protocol = "ping"
}
resource "sakuracloud_load_balancer_server" "server02"{
    load_balancer_vip_id = "${sakuracloud_load_balancer_vip.vip1.id}"
    ipaddress = "192.168.11.52"
    check_protocol = "ping"
}
`
