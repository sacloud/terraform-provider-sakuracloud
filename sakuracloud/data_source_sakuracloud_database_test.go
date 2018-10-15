package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceDatabase_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudDatabaseDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceDatabaseBase(name),
				Check:  testAccCheckSakuraCloudDatabaseDataSourceID("sakuracloud_database.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDatabaseConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseDataSourceID("data.sakuracloud_database.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_database.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_database.foobar", "plan", "10g"),
					resource.TestCheckResourceAttr("data.sakuracloud_database.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_database.foobar", "tags.#", "3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDatabaseConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseDataSourceID("data.sakuracloud_database.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDatabase_NameSelector_Exists(name, randString1, randString2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseDataSourceID("data.sakuracloud_database.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDatabase_TagSelector_Exists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseDataSourceID("data.sakuracloud_database.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDatabaseConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseDataSourceNotExists("data.sakuracloud_database.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDatabaseConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseDataSourceNotExists("data.sakuracloud_database.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDatabase_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseDataSourceNotExists("data.sakuracloud_database.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDatabase_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDatabaseDataSourceNotExists("data.sakuracloud_database.foobar"),
				),
				Destroy: true,
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

func testAccCheckSakuraCloudDataSourceDatabaseBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "sw" {
    name = "%s"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {
    name = "%s"
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

    backup_weekdays = ["mon", "tue"]
    backup_time = "00:00"
    zone = "tk1a"
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceDatabaseConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "sw" {
    name = "%s"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {
    name = "%s"
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

    backup_weekdays = ["mon", "tue"]
    backup_time = "00:00"
    zone = "tk1a"
}
data "sakuracloud_database" "foobar" {
    filter = {
	name = "Name"
	values = ["%s"]
    }
    zone = "tk1a"
}`, name, name, name)
}

func testAccCheckSakuraCloudDataSourceDatabaseConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "sw" {
    name = "%s"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {
    name = "%s"
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

    backup_weekdays = ["mon", "tue"]
    backup_time = "00:00"
    zone = "tk1a"

}
data "sakuracloud_database" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
    zone = "tk1a"
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceDatabaseConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "sw" {
    name = "%s"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {
    name = "%s"
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

    backup_weekdays = ["mon", "tue"]
    backup_time = "00:00"
    zone = "tk1a"

}
data "sakuracloud_database" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
    zone = "tk1a"
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceDatabaseConfig_NotExists(name string) string {

	return fmt.Sprintf(`
resource "sakuracloud_switch" "sw" {
    name = "%s"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {
    name = "%s"
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

    backup_weekdays = ["mon", "tue"]
    backup_time = "00:00"
    zone = "tk1a"

}
data "sakuracloud_database" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
    zone = "tk1a"
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceDatabase_NameSelector_Exists(name, p1, p2 string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "sw" {
    name = "%s"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {
    name = "%s"
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

    backup_weekdays = ["mon", "tue"]
    backup_time = "00:00"
    zone = "tk1a"

}
data "sakuracloud_database" "foobar" {
    name_selectors = ["%s", "%s"]
    zone = "tk1a"
}`, name, name, p1, p2)
}

var testAccCheckSakuraCloudDataSourceDatabase_NameSelector_NotExists = `
data "sakuracloud_database" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
    zone = "tk1a"
}
`

func testAccCheckSakuraCloudDataSourceDatabase_TagSelector_Exists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "sw" {
    name = "%s"
    zone = "tk1a"
}
resource "sakuracloud_database" "foobar" {
    name = "%s"
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

    backup_weekdays = ["mon", "tue"]
    backup_time = "00:00"
    zone = "tk1a"

}
data "sakuracloud_database" "foobar" {
	tag_selectors = ["tag1","tag2","tag3"]
    zone = "tk1a"
}`, name, name)
}

var testAccCheckSakuraCloudDataSourceDatabase_TagSelector_NotExists = `
data "sakuracloud_database" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
    zone = "tk1a"
}`
