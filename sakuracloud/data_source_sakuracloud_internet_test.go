package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceInternet_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudInternetDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceInternetBase,
				Check:  testAccCheckSakuraCloudInternetDataSourceID("sakuracloud_internet.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceID("data.sakuracloud_internet.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "tags.2", "tag3"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "nw_mask_len", "28"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "band_width", "100"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "server_ids.#", "0"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "ipaddresses.#", "11"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternetConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceID("data.sakuracloud_internet.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternet_NameSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceID("data.sakuracloud_internet.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternet_TagSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceID("data.sakuracloud_internet.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternetConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceNotExists("data.sakuracloud_internet.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternetConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceNotExists("data.sakuracloud_internet.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternet_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceNotExists("data.sakuracloud_internet.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternet_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceNotExists("data.sakuracloud_internet.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudInternetDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Internet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Internet data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudInternetDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found Internet data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudInternetDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_internet" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Internet.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Internet still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceInternetBase = `
resource "sakuracloud_internet" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
`

var testAccCheckSakuraCloudDataSourceInternetConfig = `
resource "sakuracloud_internet" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`

var testAccCheckSakuraCloudDataSourceInternetConfig_With_Tag = `
resource "sakuracloud_internet" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`

var testAccCheckSakuraCloudDataSourceInternetConfig_With_NotExists_Tag = `
resource "sakuracloud_internet" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceInternetConfig_NotExists = `
resource "sakuracloud_internet" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceInternet_NameSelector_Exists = `
resource "sakuracloud_internet" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
    name_selectors = ["name", "test"]
}
`

var testAccCheckSakuraCloudDataSourceInternet_NameSelector_NotExists = `
data "sakuracloud_internet" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}
`

var testAccCheckSakuraCloudDataSourceInternet_TagSelector_Exists = `
resource "sakuracloud_internet" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
	tag_selectors = ["tag1","tag2","tag3"]
}`

var testAccCheckSakuraCloudDataSourceInternet_TagSelector_NotExists = `
data "sakuracloud_internet" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
}`
