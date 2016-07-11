package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudLoadBalancerDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudLoadBalancerDataSourceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceLoadBalancerBase,
				Check:  testAccCheckSakuraCloudLoadBalancerDataSourceID("sakuracloud_load_balancer.foobar"),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceLoadBalancerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerDataSourceID("data.sakuracloud_load_balancer.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "zone", "tk1v"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "tags.2", "tag3"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceLoadBalancerConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerDataSourceID("data.sakuracloud_load_balancer.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceLoadBalancerConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerDataSourceNotExists("data.sakuracloud_load_balancer.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceLoadBalancerConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerDataSourceNotExists("data.sakuracloud_load_balancer.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudLoadBalancerDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find LoadBalancer data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("LoadBalancer data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudLoadBalancerDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found LoadBalancer data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudLoadBalancerDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "tk1v"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_load_balancer" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.LoadBalancer.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("LoadBalancer still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceLoadBalancerBase = `
resource sakuracloud_switch "sw"{
    name = "sw"
    zone = "tk1v"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    VRID = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    zone = "tk1v"
    tags = ["tag1","tag2","tag3"]
}`

var testAccCheckSakuraCloudDataSourceLoadBalancerConfig = `
resource sakuracloud_switch "sw"{
    name = "sw"
    zone = "tk1v"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    VRID = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    zone = "tk1v"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_load_balancer" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
    zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceLoadBalancerConfig_With_Tag = `
resource sakuracloud_switch "sw"{
    name = "sw"
    zone = "tk1v"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    VRID = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    zone = "tk1v"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_load_balancer" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
    zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceLoadBalancerConfig_With_NotExists_Tag = `
resource sakuracloud_switch "sw"{
    name = "sw"
    zone = "tk1v"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    VRID = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    zone = "tk1v"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_load_balancer" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
    zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceLoadBalancerConfig_NotExists = `
resource sakuracloud_switch "sw"{
    name = "sw"
    zone = "tk1v"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    VRID = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    zone = "tk1v"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_load_balancer" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
    zone = "tk1v"
}`
