package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudDiskDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudDiskDataSourceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceDiskConfigBase,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_disk.disk01", "name", "hoge_Ubuntu_fuga"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceDiskConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskDataSourceID("data.sakuracloud_disk.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "name", "hoge_Ubuntu_fuga"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "plan", "4"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "connection", "virtio"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "size", "20"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "description", "source_disk_description"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "tags.2", "tag3"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceDiskConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskDataSourceID("data.sakuracloud_disk.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceDiskConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskDataSourceNotExists("data.sakuracloud_disk.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceDiskConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskDataSourceNotExists("data.sakuracloud_disk.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudDiskDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Disk data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Disk data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudDiskDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found Disk data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudDiskDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_disk" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Disk.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Disk still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceDiskConfigBase = `
resource "sakuracloud_disk" "disk01"{
    name = "hoge_Ubuntu_fuga"
    tags = ["tag1","tag2","tag3"]
    description = "source_disk_description"
}
`

var testAccCheckSakuraCloudDataSourceDiskConfig = `
resource "sakuracloud_disk" "disk01"{
    name = "hoge_Ubuntu_fuga"
    tags = ["tag1","tag2","tag3"]
    description = "source_disk_description"
}

data "sakuracloud_disk" "foobar" {
    filter = {
	name = "Name"
	values = ["Ubuntu"]
    }
}`

var testAccCheckSakuraCloudDataSourceDiskConfig_With_Tag = `
resource "sakuracloud_disk" "disk01"{
    name = "hoge_Ubuntu_fuga"
    tags = ["tag1","tag2","tag3"]
    description = "source_disk_description"
}

data "sakuracloud_disk" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag2","tag3"]
    }
}`

var testAccCheckSakuraCloudDataSourceDiskConfig_With_NotExists_Tag = `
resource "sakuracloud_disk" "disk01"{
    name = "hoge_Ubuntu_fuga"
    tags = ["tag1","tag2","tag3"]
    description = "source_disk_description"
}

data "sakuracloud_disk" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag2","tag4"]
    }
}`

var testAccCheckSakuraCloudDataSourceDiskConfig_NotExists = `
resource "sakuracloud_disk" "disk01"{
    name = "hoge_Ubuntu_fuga"
    tags = ["tag1","tag2","tag3"]
    description = "source_disk_description"
}

data "sakuracloud_disk" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`
