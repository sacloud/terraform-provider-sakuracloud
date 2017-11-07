package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"
)

func TestAccSakuraCloudLoadBalancerDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudLoadBalancerDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceLoadBalancerBase,
				Check:  testAccCheckSakuraCloudLoadBalancerDataSourceID("sakuracloud_load_balancer.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceLoadBalancerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerDataSourceID("data.sakuracloud_load_balancer.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "tags.2", "tag3"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceLoadBalancerConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerDataSourceID("data.sakuracloud_load_balancer.foobar"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceLoadBalancerConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerDataSourceNotExists("data.sakuracloud_load_balancer.foobar"),
				),
			},
			{
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
			return errors.New("LoadBalancer data source ID not set")
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
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_load_balancer" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.LoadBalancer.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("LoadBalancer still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceLoadBalancerBase = `
resource sakuracloud_switch "sw"{
    name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    vrid = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}`

var testAccCheckSakuraCloudDataSourceLoadBalancerConfig = `
resource sakuracloud_switch "sw"{
    name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    vrid = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_load_balancer" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`

var testAccCheckSakuraCloudDataSourceLoadBalancerConfig_With_Tag = `
resource sakuracloud_switch "sw"{
    name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    vrid = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_load_balancer" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`

var testAccCheckSakuraCloudDataSourceLoadBalancerConfig_With_NotExists_Tag = `
resource sakuracloud_switch "sw"{
    name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    vrid = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_load_balancer" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceLoadBalancerConfig_NotExists = `
resource sakuracloud_switch "sw"{
    name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    vrid = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_load_balancer" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`
