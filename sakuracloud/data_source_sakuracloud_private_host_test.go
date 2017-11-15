package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"regexp"
	"testing"
)

func TestAccSakuraCloudDataSourcePrivateHost_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudPrivateHostDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostBase,
				Check:  testAccCheckSakuraCloudPrivateHostDataSourceID("sakuracloud_private_host.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceID("data.sakuracloud_private_host.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "tags.2", "tag3"),
					resource.TestMatchResourceAttr("data.sakuracloud_private_host.foobar",
						"hostname",
						regexp.MustCompile(".+")), // should be not empty
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceID("data.sakuracloud_private_host.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHost_NameSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceID("data.sakuracloud_private_host.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHost_TagSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceID("data.sakuracloud_private_host.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceNotExists("data.sakuracloud_private_host.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceNotExists("data.sakuracloud_private_host.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHost_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceNotExists("data.sakuracloud_private_host.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHost_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceNotExists("data.sakuracloud_private_host.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudPrivateHostDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find PrivateHost data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("PrivateHost data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudPrivateHostDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found PrivateHost data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudPrivateHostDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	originalZone := client.Zone
	client.Zone = "tk1a"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_private_host" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.PrivateHost.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("PrivateHost still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourcePrivateHostBase = `
resource "sakuracloud_private_host" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
`

var testAccCheckSakuraCloudDataSourcePrivateHostConfig = `
resource "sakuracloud_private_host" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
data "sakuracloud_private_host" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
    zone = "tk1a"
}`

var testAccCheckSakuraCloudDataSourcePrivateHostConfig_With_Tag = `
resource "sakuracloud_private_host" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
data "sakuracloud_private_host" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
    zone = "tk1a"
}`

var testAccCheckSakuraCloudDataSourcePrivateHostConfig_With_NotExists_Tag = `
resource "sakuracloud_private_host" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
data "sakuracloud_private_host" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
    zone = "tk1a"
}`

var testAccCheckSakuraCloudDataSourcePrivateHostConfig_NotExists = `
resource "sakuracloud_private_host" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
data "sakuracloud_private_host" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
    zone = "tk1a"
}`

var testAccCheckSakuraCloudDataSourcePrivateHost_NameSelector_Exists = `
resource "sakuracloud_private_host" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
data "sakuracloud_private_host" "foobar" {
    name_selectors = ["name", "test"]
    zone = "tk1a"
}
`
var testAccCheckSakuraCloudDataSourcePrivateHost_NameSelector_NotExists = `
data "sakuracloud_private_host" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
    zone = "tk1a"
}
`

var testAccCheckSakuraCloudDataSourcePrivateHost_TagSelector_Exists = `
resource "sakuracloud_private_host" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
data "sakuracloud_private_host" "foobar" {
	tag_selectors = ["tag1","tag2","tag3"]
    zone = "tk1a"
}`

var testAccCheckSakuraCloudDataSourcePrivateHost_TagSelector_NotExists = `
data "sakuracloud_private_host" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
    zone = "tk1a"
}`
