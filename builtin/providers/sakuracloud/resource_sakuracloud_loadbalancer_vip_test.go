package sakuracloud

import (
	"errors"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"testing"
)

func TestAccSakuraCloudLoadBalancerVIP(t *testing.T) {
	var loadBalancer sacloud.LoadBalancer
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudLoadBalancerVIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudLoadBalancerVIPConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerExists("sakuracloud_load_balancer.foobar", &loadBalancer),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "vip", "192.168.11.201"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "port", "80"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "delay_loop", "100"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "sorry_server", "192.168.11.11"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip2", "vip", "192.168.11.202"),
				),
			},
			{
				Config: testAccCheckSakuraCloudLoadBalancerVIPConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer.foobar", "vip_ids.#", "2"),
				),
			},
			{
				Config: testAccCheckSakuraCloudLoadBalancerVIPConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerExists("sakuracloud_load_balancer.foobar", &loadBalancer),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "vip", "192.168.11.201"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "port", "80"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "delay_loop", "50"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "sorry_server", "192.168.11.22"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip2", "vip", "192.168.11.202"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip3", "vip", "192.168.11.203"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip4", "vip", "192.168.11.204"),
				),
			},
			{
				Config: testAccCheckSakuraCloudLoadBalancerVIPConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer.foobar", "vip_ids.#", "4"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudLoadBalancerVIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

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

var testAccCheckSakuraCloudLoadBalancerVIPConfig_basic = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    VRID = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name"
    description = "description"
    tags = ["hoge1" , "hoge2"]
}
resource "sakuracloud_load_balancer_vip" "vip1" {
    load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
    vip = "192.168.11.201"
    port = 80
    delay_loop = 100
    sorry_server = "192.168.11.11"
}
resource "sakuracloud_load_balancer_vip" "vip2" {
    load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
    vip = "192.168.11.202"
    port = 80
}
`

var testAccCheckSakuraCloudLoadBalancerVIPConfig_update = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    VRID = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name"
    description = "description"
    tags = ["hoge1" , "hoge2"]
}
resource "sakuracloud_load_balancer_vip" "vip1" {
    load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
    vip = "192.168.11.201"
    port = 80
    delay_loop = 50
    sorry_server = "192.168.11.22"
}
resource "sakuracloud_load_balancer_vip" "vip2" {
    load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
    vip = "192.168.11.202"
    port = 80
}
resource "sakuracloud_load_balancer_vip" "vip3" {
    load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
    vip = "192.168.11.203"
    port = 80
}
resource "sakuracloud_load_balancer_vip" "vip4" {
    load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
    vip = "192.168.11.204"
    port = 80
}
`
