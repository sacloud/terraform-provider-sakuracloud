package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"
)

func TestAccSakuraCloudServerDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudServerDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceServerBase,
				Check:  testAccCheckSakuraCloudServerDataSourceID("sakuracloud_server.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceID("data.sakuracloud_server.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "interface_driver", "virtio"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.2", "tag3"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "core", "1"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "memory", "1"),
					//resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "disks.#", "1"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "nic", "shared"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "additional_nics.#", "0"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "macaddresses.#", "1"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceServerConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceID("data.sakuracloud_server.foobar"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceServerConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceNotExists("data.sakuracloud_server.foobar"),
				),
			},
			{
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
			return errors.New("Server data source ID not set")
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
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_server" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Server.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Server still exists")
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
}
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}
resource "sakuracloud_server" "foobar" {
    name = "name_test"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}`

var testAccCheckSakuraCloudDataSourceServerConfig = fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`, testAccCheckSakuraCloudDataSourceServerBase)

var testAccCheckSakuraCloudDataSourceServerConfig_With_Tag = fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`, testAccCheckSakuraCloudDataSourceServerBase)

var testAccCheckSakuraCloudDataSourceServerConfig_With_NotExists_Tag = fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`, testAccCheckSakuraCloudDataSourceServerBase)

var testAccCheckSakuraCloudDataSourceServerConfig_NotExists = fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`, testAccCheckSakuraCloudDataSourceServerBase)
