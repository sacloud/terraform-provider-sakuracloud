package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccSakuraCloudDatabase_Basic(t *testing.T) {
	var database sacloud.Database
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudDatabaseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDatabaseConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseExists("sakuracloud_database.foobar", &database),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "name", "name_before"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "description", "description_before"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "tags.0", "hoge1"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "tags.1", "hoge2"),
					//resource.TestCheckResourceAttr("sakuracloud_database.foobar", "plan", "mini"),
					//resource.TestCheckResourceAttr("sakuracloud_database.foobar", "is_double", "false"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "admin_password", "DatabasePasswordAdmin397"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "user_name", "defuser"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "user_password", "DatabasePasswordUser397"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "allow_networks.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "allow_networks.0", "192.168.11.0/24"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "allow_networks.1", "192.168.12.0/24"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "port", "54321"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "backup_rotate", "8"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "backup_time", "00:00"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "switch_id", "shared"),
					//resource.TestCheckResourceAttr("sakuracloud_database.foobar", "ipaddress1", ""),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "nw_mask_len", "0"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "default_route", ""),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudDatabaseConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseExists("sakuracloud_database.foobar", &database),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "name", "name_after"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "description", "description_after"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "tags.0", "hoge1_after"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "tags.1", "hoge2_after"),
					//resource.TestCheckResourceAttr("sakuracloud_database.foobar", "plan", "mini"),
					//resource.TestCheckResourceAttr("sakuracloud_database.foobar", "is_double", "false"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "admin_password", "DatabasePasswordAdmin397"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "user_name", "defuser"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "user_password", "DatabasePasswordUser397_upd"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "allow_networks.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "allow_networks.0", "192.168.110.0/24"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "allow_networks.1", "192.168.120.0/24"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "port", "54322"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "backup_rotate", "7"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "backup_time", "00:30"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "switch_id", "shared"),
					//resource.TestCheckResourceAttr("sakuracloud_database.foobar", "ipaddress1", ""),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "nw_mask_len", "0"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "default_route", ""),
				),
			},
		},
	})
}

func TestAccResourceSakuraCloudDatabase_WithSwitch(t *testing.T) {
	var database sacloud.Database
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudDatabaseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDatabaseConfig_WithSwitch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseExists("sakuracloud_database.foobar", &database),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "name", "name_before"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "description", "description_before"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "tags.0", "hoge1"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "tags.1", "hoge2"),
					//resource.TestCheckResourceAttr("sakuracloud_database.foobar", "plan", "mini"),
					//resource.TestCheckResourceAttr("sakuracloud_database.foobar", "is_double", "false"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "admin_password", "DatabasePasswordAdmin397"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "user_name", "defuser"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "user_password", "DatabasePasswordUser397"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "allow_networks.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "allow_networks.0", "192.168.11.0/24"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "allow_networks.1", "192.168.12.0/24"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "port", "54321"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "backup_rotate", "8"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "backup_time", "00:00"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "ipaddress1", "192.168.11.101"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "nw_mask_len", "24"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "default_route", "192.168.11.1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudDatabaseConfig_WithSwitchUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseExists("sakuracloud_database.foobar", &database),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "name", "name_after"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "description", "description_after"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "tags.0", "hoge1_after"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "tags.1", "hoge2_after"),
					//resource.TestCheckResourceAttr("sakuracloud_database.foobar", "plan", "mini"),
					//resource.TestCheckResourceAttr("sakuracloud_database.foobar", "is_double", "false"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "admin_password", "DatabasePasswordAdmin397"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "user_name", "defuser"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "user_password", "DatabasePasswordUser397_upd"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "allow_networks.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "allow_networks.0", "192.168.110.0/24"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "allow_networks.1", "192.168.120.0/24"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "port", "54322"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "backup_rotate", "7"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "backup_time", "00:30"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "ipaddress1", "192.168.11.101"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "nw_mask_len", "24"),
					resource.TestCheckResourceAttr("sakuracloud_database.foobar", "default_route", "192.168.11.1"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudDatabaseExists(n string, database *sacloud.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Database ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)
		originalZone := client.Zone
		client.Zone = "tk1a"
		defer func() { client.Zone = originalZone }()

		foundDatabase, err := client.Database.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundDatabase.ID != rs.Primary.ID {
			return fmt.Errorf("Database not found")
		}

		*database = *foundDatabase

		return nil
	}
}

func testAccCheckSakuraCloudDatabaseDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "tk1a"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_database" {
			continue
		}

		_, err := client.Database.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Database still exists")
		}
	}

	return nil
}

const testAccCheckSakuraCloudDatabaseConfig_basic = `
resource "sakuracloud_database" "foobar" {

    admin_password = "DatabasePasswordAdmin397"
    user_name = "defuser"
    user_password = "DatabasePasswordUser397"

    allow_networks = ["192.168.11.0/24","192.168.12.0/24"]

    port = 54321

    backup_rotate = 8
    backup_time = "00:00"

    name = "name_before"
    description = "description_before"
    tags = ["hoge1" , "hoge2"]
    zone = "tk1a"
}`

const testAccCheckSakuraCloudDatabaseConfig_update = `
resource "sakuracloud_database" "foobar" {

    admin_password = "DatabasePasswordAdmin397"
    user_name = "defuser"
    user_password = "DatabasePasswordUser397_upd"

    allow_networks = ["192.168.110.0/24","192.168.120.0/24"]

    port = 54322

    backup_rotate = 7
    backup_time = "00:30"
    name = "name_after"
    description = "description_after"
    tags = ["hoge1_after" , "hoge2_after"]
    zone = "tk1a"
}`

const testAccCheckSakuraCloudDatabaseConfig_WithSwitch = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {

    admin_password = "DatabasePasswordAdmin397"
    user_name = "defuser"
    user_password = "DatabasePasswordUser397"

    allow_networks = ["192.168.11.0/24","192.168.12.0/24"]

    port = 54321

    backup_rotate = 8
    backup_time = "00:00"


    switch_id = "${sakuracloud_switch.sw.id}"
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    name = "name_before"
    description = "description_before"
    tags = ["hoge1" , "hoge2"]
    zone = "tk1a"
}`

const testAccCheckSakuraCloudDatabaseConfig_WithSwitchUpdate = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {

    admin_password = "DatabasePasswordAdmin397"
    user_name = "defuser"
    user_password = "DatabasePasswordUser397_upd"

    allow_networks = ["192.168.110.0/24","192.168.120.0/24"]

    port = 54322

    backup_rotate = 7
    backup_time = "00:30"

    name = "name_after"
    description = "description_after"
    tags = ["hoge1_after" , "hoge2_after"]

    switch_id = "${sakuracloud_switch.sw.id}"
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    default_route = "192.168.11.1"

    zone = "tk1a"
}`
