package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceBridge_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudBridgeDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeBase,
				Check:  testAccCheckSakuraCloudBridgeDataSourceID("sakuracloud_bridge.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeDataSourceID("data.sakuracloud_bridge.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "description", "description_test"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeConfig_NameSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeDataSourceID("data.sakuracloud_bridge.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "description", "description_test"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeDataSourceNotExists("data.sakuracloud_bridge.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeConfig_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeDataSourceNotExists("data.sakuracloud_bridge.foobar"),
				),
				Destroy: true,
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
			return errors.New("Bridge data source ID not set")
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
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_bridge" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Bridge.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Bridge still exists")
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
data "sakuracloud_bridge" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceBridgeConfig_NameSelector_Exists = `
resource "sakuracloud_bridge" "foobar" {
    name = "name_test"
    description = "description_test"
}
data "sakuracloud_bridge" "foobar" {
    name_selectors = ["name", "test"]
}`

var testAccCheckSakuraCloudDataSourceBridgeConfig_NameSelector_NotExists = `
data "sakuracloud_bridge" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}`
