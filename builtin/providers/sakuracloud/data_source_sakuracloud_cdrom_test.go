package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudCDROMDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudCDROMDataSourceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceCDROMConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudCDROMDataSourceID("data.sakuracloud_cdrom.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "name", "Ubuntu server 16.04.1 LTS 64bit"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "size", "5"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.#", "5"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.0", "arch-64bit"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.1", "current-stable"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.2", "distro-ubuntu"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.3", "distro-ver-16.04.1"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.4", "os-linux"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceCDROMConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudCDROMDataSourceID("data.sakuracloud_cdrom.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceCDROMConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudCDROMDataSourceNotExists("data.sakuracloud_cdrom.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceCDROMConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudCDROMDataSourceNotExists("data.sakuracloud_cdrom.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudCDROMDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find CDROM data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CDROM data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudCDROMDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found CDROM data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudCDROMDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_cdrom" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.CDROM.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("CDROM still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceCDROMConfig = `
data "sakuracloud_cdrom" "foobar" {
    filter = {
	name = "Name"
	values = ["Ubuntu Server 16"]
    }
}`

var testAccCheckSakuraCloudDataSourceCDROMConfig_With_Tag = `
data "sakuracloud_cdrom" "foobar" {
    filter = {
	name = "Tags"
	values = ["distro-ubuntu","os-linux"]
    }
}`

var testAccCheckSakuraCloudDataSourceCDROMConfig_With_NotExists_Tag = `
data "sakuracloud_cdrom" "foobar" {
    filter = {
	name = "Tags"
	values = ["distro-ubuntu-xxxxxxxxxxx","os-linux-xxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceCDROMConfig_NotExists = `
data "sakuracloud_cdrom" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`
