package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"
)

func TestAccSakuraCloudDatabaseDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudDatabaseDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceDatabaseBase,
				Check:  testAccCheckSakuraCloudDatabaseDataSourceID("sakuracloud_database.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDatabaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseDataSourceID("data.sakuracloud_database.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_database.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_database.foobar", "plan", "10g"),
					resource.TestCheckResourceAttr("data.sakuracloud_database.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_database.foobar", "tags.#", "3"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceDatabaseConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseDataSourceID("data.sakuracloud_database.foobar"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceDatabaseConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseDataSourceNotExists("data.sakuracloud_database.foobar"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceDatabaseConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseDataSourceNotExists("data.sakuracloud_database.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudDatabaseDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Database data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Database data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudDatabaseDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found Database data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudDatabaseDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	originalZone := client.Zone
	client.Zone = "tk1a"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_database" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Database.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Database still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceDatabaseBase = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]

    user_name = "defuser"
    user_password = "DatabasePasswordUser397"

    allow_networks = ["192.168.11.0/24","192.168.12.0/24"]

    switch_id = "${sakuracloud_switch.sw.id}"
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    port = 54321

    backup_time = "00:00"
    zone = "tk1a"
}`

var testAccCheckSakuraCloudDataSourceDatabaseConfig = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]

    user_name = "defuser"
    user_password = "DatabasePasswordUser397"

    allow_networks = ["192.168.11.0/24","192.168.12.0/24"]

    switch_id = "${sakuracloud_switch.sw.id}"
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    port = 54321

    backup_time = "00:00"
    zone = "tk1a"

}
data "sakuracloud_database" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
    zone = "tk1a"
}`

var testAccCheckSakuraCloudDataSourceDatabaseConfig_With_Tag = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]

    user_name = "defuser"
    user_password = "DatabasePasswordUser397"

    allow_networks = ["192.168.11.0/24","192.168.12.0/24"]

    switch_id = "${sakuracloud_switch.sw.id}"
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    port = 54321

    backup_time = "00:00"
    zone = "tk1a"

}
data "sakuracloud_database" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
    zone = "tk1a"
}`

var testAccCheckSakuraCloudDataSourceDatabaseConfig_With_NotExists_Tag = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]

    user_name = "defuser"
    user_password = "DatabasePasswordUser397"

    allow_networks = ["192.168.11.0/24","192.168.12.0/24"]

    switch_id = "${sakuracloud_switch.sw.id}"
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    port = 54321

    backup_time = "00:00"
    zone = "tk1a"

}
data "sakuracloud_database" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
    zone = "tk1a"
}`

var testAccCheckSakuraCloudDataSourceDatabaseConfig_NotExists = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]

    user_name = "defuser"
    user_password = "DatabasePasswordUser397"

    allow_networks = ["192.168.11.0/24","192.168.12.0/24"]

    switch_id = "${sakuracloud_switch.sw.id}"
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    port = 54321

    backup_time = "00:00"
    zone = "tk1a"

}
data "sakuracloud_database" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
    zone = "tk1a"
}`
