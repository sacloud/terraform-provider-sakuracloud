package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceBridge_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudBridgeDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeBase(name),
				Check:  testAccCheckSakuraCloudBridgeDataSourceID("sakuracloud_bridge.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeDataSourceID("data.sakuracloud_bridge.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "description", "description_test"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeConfig_NameSelector_Exists(name, randString1, randString2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeDataSourceID("data.sakuracloud_bridge.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "name", name),
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
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
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

func testAccCheckSakuraCloudDataSourceBridgeBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_bridge" "foobar" {
    name = "%s"
    description = "description_test"
} 
`, name)
}

func testAccCheckSakuraCloudDataSourceBridgeConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_bridge" "foobar" {
    name = "%s"
    description = "description_test"
}
data "sakuracloud_bridge" "foobar" {
  filter {
    name = "Name"
    values = ["%s"]
  }
}`, name, name)
}

var testAccCheckSakuraCloudDataSourceBridgeConfig_NotExists = `
data "sakuracloud_bridge" "foobar" {
  filter {
    name = "Name"
    values = ["xxxxxxxxxxxxxxxxxx"]
  }
}`

func testAccCheckSakuraCloudDataSourceBridgeConfig_NameSelector_Exists(name, p1, p2 string) string {
	return fmt.Sprintf(`
resource "sakuracloud_bridge" "foobar" {
  name = "%s"
  description = "description_test"
}
data "sakuracloud_bridge" "foobar" {
  name_selectors = ["%s", "%s"]
}`, name, p1, p2)
}

var testAccCheckSakuraCloudDataSourceBridgeConfig_NameSelector_NotExists = `
data "sakuracloud_bridge" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}`
