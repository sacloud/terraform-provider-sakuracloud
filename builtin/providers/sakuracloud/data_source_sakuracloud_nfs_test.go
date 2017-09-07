package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudNFSDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudNFSDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceNFSBase,
				Check:  testAccCheckSakuraCloudNFSDataSourceID("sakuracloud_nfs.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNFSConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNFSDataSourceID("data.sakuracloud_nfs.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_nfs.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_nfs.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_nfs.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_nfs.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_nfs.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_nfs.foobar", "tags.2", "tag3"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceNFSConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNFSDataSourceID("data.sakuracloud_nfs.foobar"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceNFSConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNFSDataSourceNotExists("data.sakuracloud_nfs.foobar"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceNFSConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNFSDataSourceNotExists("data.sakuracloud_nfs.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudNFSDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find NFS data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("NFS data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudNFSDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found NFS data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudNFSDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_nfs" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.NFS.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("NFS still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceNFSBase = `
resource sakuracloud_switch "sw"{
    name = "sw"
}
resource "sakuracloud_nfs" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    ipaddress = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}`

var testAccCheckSakuraCloudDataSourceNFSConfig = `
resource sakuracloud_switch "sw"{
    name = "sw"
}
resource "sakuracloud_nfs" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    ipaddress = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_nfs" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`

var testAccCheckSakuraCloudDataSourceNFSConfig_With_Tag = `
resource sakuracloud_switch "sw"{
    name = "sw"
}
resource "sakuracloud_nfs" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    ipaddress = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_nfs" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`

var testAccCheckSakuraCloudDataSourceNFSConfig_With_NotExists_Tag = `
resource sakuracloud_switch "sw"{
    name = "sw"
}
resource "sakuracloud_nfs" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    ipaddress = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_nfs" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceNFSConfig_NotExists = `
resource sakuracloud_switch "sw"{
    name = "sw"
}
resource "sakuracloud_nfs" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    ipaddress = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_nfs" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`
