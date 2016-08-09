package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccResourceSakuraCloudLoadBalancer_Basic(t *testing.T) {
	var loadBalancer sacloud.LoadBalancer
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerExists("sakuracloud_load_balancer.foobar", &loadBalancer),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "name", "name_before"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "description", "description_before"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "tags.0", "hoge1"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "tags.1", "hoge2"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "VRID", "1"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "ipaddress1", "192.168.11.101"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "ipaddress2", ""),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "nw_mask_len", "24"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "default_route", "192.168.11.1"),
				),
			},
		},
	})
}

func TestAccResourceSakuraCloudLoadBalancer_Update(t *testing.T) {
	var loadBalancer sacloud.LoadBalancer
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerExists("sakuracloud_load_balancer.foobar", &loadBalancer),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "name", "name_before"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "description", "description_before"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "tags.0", "hoge1"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "tags.1", "hoge2"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "VRID", "1"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "ipaddress1", "192.168.11.101"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "ipaddress2", ""),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "nw_mask_len", "24"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "default_route", "192.168.11.1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudLoadBalancerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerExists("sakuracloud_load_balancer.foobar", &loadBalancer),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "name", "name_after"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "description", "description_after"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "tags.0", "hoge1_after"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "tags.1", "hoge2_after"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "VRID", "1"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "ipaddress1", "192.168.11.101"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "ipaddress2", ""),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "nw_mask_len", "24"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "default_route", "192.168.11.1"),
				),
			},
		},
	})
}

func TestAccResourceSakuraCloudLoadBalancer_WithRouter(t *testing.T) {
	var loadBalancer sacloud.LoadBalancer
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudLoadBalancerConfig_WithRouter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerExists("sakuracloud_load_balancer.foobar", &loadBalancer),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "name", "name_before"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "description", "description_before"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "tags.0", "hoge1"),
					resource.TestCheckResourceAttr("sakuracloud_load_balancer.foobar", "tags.1", "hoge2"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudLoadBalancerExists(n string, loadBalancer *sacloud.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LoadBalancer ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)

		foundLoadBalancer, err := client.LoadBalancer.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundLoadBalancer.ID != rs.Primary.ID {
			return fmt.Errorf("LoadBalancer not found")
		}

		*loadBalancer = *foundLoadBalancer

		return nil
	}
}

func testAccCheckSakuraCloudLoadBalancerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_load_balancer" {
			continue
		}

		_, err := client.LoadBalancer.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("LoadBalancer still exists")
		}
	}

	return nil
}

const testAccCheckSakuraCloudLoadBalancerConfig_basic = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    VRID = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_before"
    description = "description_before"
    tags = ["hoge1" , "hoge2"]
}`

const testAccCheckSakuraCloudLoadBalancerConfig_update = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    VRID = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_after"
    description = "description_after"
    tags = ["hoge1_after" , "hoge2_after"]
}`

const testAccCheckSakuraCloudLoadBalancerConfig_WithRouter = `
resource "sakuracloud_internet" "router" {
    name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_internet.router.switch_id}"
    is_double = true
    plan = "highspec"
    VRID = 1
    ipaddress1 = "${sakuracloud_internet.router.nw_ipaddresses.0}"
    ipaddress2 = "${sakuracloud_internet.router.nw_ipaddresses.1}"
    nw_mask_len = "${sakuracloud_internet.router.nw_mask_len}"
    default_route = "${sakuracloud_internet.router.nw_gateway}"

    name = "name_before"
    description = "description_before"
    tags = ["hoge1" , "hoge2"]
}`
