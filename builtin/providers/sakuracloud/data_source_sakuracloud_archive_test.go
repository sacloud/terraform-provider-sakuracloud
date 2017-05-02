package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudArchiveDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudArchiveDataSourceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceArchiveConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceID("data.sakuracloud_archive.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "name", "Ubuntu Server 16.04.2 LTS 64bit"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "size", "20"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "zone", "tk1v"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.#", "6"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.0", "@size-extendable"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.1", "arch-64bit"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.2", "current-stable"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.3", "distro-ubuntu"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.4", "distro-ver-16.04.2"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.5", "os-linux"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceArchive_OSType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceID("data.sakuracloud_archive.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceArchiveConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceID("data.sakuracloud_archive.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceArchiveConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceNotExists("data.sakuracloud_archive.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceArchiveConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceNotExists("data.sakuracloud_archive.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudArchiveDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Archive data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Archive data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudArchiveDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found Archive data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudArchiveDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "tk1v"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_archive" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Archive.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("Archive still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceArchiveConfig = `
data "sakuracloud_archive" "foobar" {
    filter = {
	name = "Name"
	values = ["Ubuntu Server 16"]
    }
    zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchiveConfig_With_Tag = `
data "sakuracloud_archive" "foobar" {
    filter = {
	name = "Tags"
	values = ["distro-ubuntu","os-linux"]
    }
    zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchiveConfig_With_NotExists_Tag = `
data "sakuracloud_archive" "foobar" {
    filter = {
	name = "Tags"
	values = ["distro-ubuntu-xxxxxxxxxxx","os-linux-xxxxxxxx"]
    }
    zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchiveConfig_NotExists = `
data "sakuracloud_archive" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
    zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchive_OSType = `
data "sakuracloud_archive" "foobar" {
    os_type = "centos"
    zone = "tk1v"
}
`
