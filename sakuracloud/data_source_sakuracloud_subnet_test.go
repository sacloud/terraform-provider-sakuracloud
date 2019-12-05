package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudSubnetDataSource_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudSubnetDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceSubnetBase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSubnetDataSourceID("sakuracloud_subnet.foobar"),
					testAccCheckSakuraCloudSubnetDataSourceID("sakuracloud_subnet.foobar2"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSubnetDataSourceID("data.sakuracloud_subnet.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_subnet.foobar", "ipaddresses.#", "16"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSubnetConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSubnetDataSourceNotExists("data.sakuracloud_subnet.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudSubnetDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Subnet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Subnet data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudSubnetDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found Subnet data source: %s", n)
		}
		return nil
	}
}

var testAccCheckSakuraCloudDataSourceSubnetBase = `
resource sakuracloud_internet "foobar" {
  name = "subnet_test"
}
resource "sakuracloud_subnet" "foobar" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  next_hop = "${sakuracloud_internet.foobar.ipaddresses[0]}"
}
resource "sakuracloud_subnet" "foobar2" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  next_hop = "${sakuracloud_internet.foobar.ipaddresses[1]}"
}
`

var testAccCheckSakuraCloudDataSourceSubnetConfig = `
resource sakuracloud_internet "foobar" {
  name = "subnet_test"
}
resource "sakuracloud_subnet" "foobar" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  next_hop = "${sakuracloud_internet.foobar.ipaddresses[0]}"
}
resource "sakuracloud_subnet" "foobar2" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  next_hop = "${sakuracloud_internet.foobar.ipaddresses[1]}"
}

data sakuracloud_subnet "foobar" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  index = 1
}
`

var testAccCheckSakuraCloudDataSourceSubnetConfig_NotExists = `
resource sakuracloud_internet "foobar" {
  name = "subnet_test"
}
resource "sakuracloud_subnet" "foobar" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  next_hop = "${sakuracloud_internet.foobar.ipaddresses[0]}"
}
resource "sakuracloud_subnet" "foobar2" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  next_hop = "${sakuracloud_internet.foobar.ipaddresses[1]}"
}
data sakuracloud_subnet "foobar" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  index = 2
}
`
