package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudBridgeDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudBridgeDataSourceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceBridgeBase,
				Check:  testAccCheckSakuraCloudBridgeDataSourceID("sakuracloud_bridge.foobar"),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceBridgeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeDataSourceID("data.sakuracloud_bridge.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "description", "description_test"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceBridgeConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeDataSourceNotExists("data.sakuracloud_bridge.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudBridgeDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Bridge data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Bridge data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudBridgeDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found Bridge data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudBridgeDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "tk1v"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_bridge" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Bridge.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("Bridge still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceBridgeBase = `
resource "sakuracloud_bridge" "foobar" {
    name = "name_test"
    description = "description_test"
}
`

var testAccCheckSakuraCloudDataSourceBridgeConfig = `
resource "sakuracloud_bridge" "foobar" {
    name = "name_test"
    description = "description_test"
}
data "sakuracloud_bridge" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`

var testAccCheckSakuraCloudDataSourceBridgeConfig_NotExists = `
resource "sakuracloud_bridge" "foobar" {
    name = "name_test"
    description = "description_test"
}
data "sakuracloud_bridge" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`
