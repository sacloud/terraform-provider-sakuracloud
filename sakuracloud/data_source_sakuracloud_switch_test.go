package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccSakuraCloudDataSourceSwitch_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudSwitchDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchBase,
				Check:  testAccCheckSakuraCloudSwitchDataSourceID("sakuracloud_switch.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchDataSourceID("data.sakuracloud_switch.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchDataSourceID("data.sakuracloud_switch.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitch_NameSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchDataSourceID("data.sakuracloud_switch.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitch_TagSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchDataSourceID("data.sakuracloud_switch.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchDataSourceNotExists("data.sakuracloud_switch.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchDataSourceNotExists("data.sakuracloud_switch.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitch_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchDataSourceNotExists("data.sakuracloud_switch.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitch_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchDataSourceNotExists("data.sakuracloud_switch.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudSwitchDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Switch data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Switch data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudSwitchDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found Switch data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudSwitchDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_switch" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Switch.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Switch still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceSwitchBase = `
resource "sakuracloud_switch" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
`

var testAccCheckSakuraCloudDataSourceSwitchConfig = `
resource "sakuracloud_switch" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`

var testAccCheckSakuraCloudDataSourceSwitchConfig_With_Tag = `
resource "sakuracloud_switch" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`

var testAccCheckSakuraCloudDataSourceSwitchConfig_With_NotExists_Tag = `
resource "sakuracloud_switch" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceSwitchConfig_NotExists = `
resource "sakuracloud_switch" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceSwitch_NameSelector_Exists = `
resource "sakuracloud_switch" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
    name_selectors = ["name", "test"]
}`

var testAccCheckSakuraCloudDataSourceSwitch_NameSelector_NotExists = `
data "sakuracloud_switch" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}`

var testAccCheckSakuraCloudDataSourceSwitch_TagSelector_Exists = `
resource "sakuracloud_switch" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
	tag_selectors = ["tag1","tag2","tag3"]
}`

var testAccCheckSakuraCloudDataSourceSwitch_TagSelector_NotExists = `
data "sakuracloud_switch" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
}`
