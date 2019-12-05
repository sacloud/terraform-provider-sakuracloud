package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceArchive_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudArchiveDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceArchiveConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceID("data.sakuracloud_archive.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "name", "Ubuntu Server 16.04.6 LTS 64bit"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "size", "20"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "zone", "tk1v"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.#", "5"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.0", "@size-extendable"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.1", "arch-64bit"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.2", "distro-ubuntu"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.3", "distro-ver-16.04.5"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.4", "os-linux"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceArchive_OSType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceID("data.sakuracloud_archive.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceArchiveConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceID("data.sakuracloud_archive.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceArchiveConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceNotExists("data.sakuracloud_archive.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceArchiveConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceNotExists("data.sakuracloud_archive.foobar"),
				),
				Destroy: true,
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
			return errors.New("Archive data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudArchiveDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found Archive data source: %s", n)
		}
		return nil
	}
}

var testAccCheckSakuraCloudDataSourceArchiveConfig = `
data "sakuracloud_archive" "foobar" {
  filters {
    names = ["Ubuntu Server 16"]
  }
  zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchiveConfig_With_Tag = `
data "sakuracloud_archive" "foobar" {
  filters {
    tags = ["distro-ubuntu","os-linux"]
  }
  zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchiveConfig_With_NotExists_Tag = `
data "sakuracloud_archive" "foobar" {
  filters {
    tags = ["distro-ubuntu-xxxxxxxxxxx","os-linux-xxxxxxxx"]
  }
  zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchiveConfig_NotExists = `
data "sakuracloud_archive" "foobar" {
  filters {
    names = ["xxxxxxxxxxxxxxxxxx"]
  }
  zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchive_OSType = `
data "sakuracloud_archive" "foobar" {
    os_type = "rancheros"
    zone    = "tk1v"
}
`
