package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudServerDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudServerDataSourceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceServerBase,
				Check:  testAccCheckSakuraCloudServerDataSourceID("sakuracloud_server.foobar"),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceServerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceID("data.sakuracloud_server.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "zone", "tk1v"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.2", "tag3"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "core", "1"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "memory", "1"),
					//resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "disks.#", "1"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "base_interface", "shared"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "additional_interfaces.#", "0"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "mac_addresses.#", "1"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceServerConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceID("data.sakuracloud_server.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceServerConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceNotExists("data.sakuracloud_server.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceServerConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceNotExists("data.sakuracloud_server.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudServerDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Server data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Server data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudServerDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found Server data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudServerDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "tk1v"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_server" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Server.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Server still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceServerBase = `
data "sakuracloud_archive" "ubuntu" {
    filter = {
	name = "Name"
	values = ["Ubuntu Server 16"]
    }
    zone = "tk1v"
}
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
    zone = "tk1v"
}
resource "sakuracloud_server" "foobar" {
    name = "name_test"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceServerConfig = fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
    zone = "tk1v"
}`, testAccCheckSakuraCloudDataSourceServerBase)

var testAccCheckSakuraCloudDataSourceServerConfig_With_Tag = fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
    zone = "tk1v"
}`, testAccCheckSakuraCloudDataSourceServerBase)

var testAccCheckSakuraCloudDataSourceServerConfig_With_NotExists_Tag = fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
    zone = "tk1v"
}`, testAccCheckSakuraCloudDataSourceServerBase)

var testAccCheckSakuraCloudDataSourceServerConfig_NotExists = fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
    zone = "tk1v"
}`, testAccCheckSakuraCloudDataSourceServerBase)
