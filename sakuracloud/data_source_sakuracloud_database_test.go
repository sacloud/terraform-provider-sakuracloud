package sakuracloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceDatabase_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
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
		},
	})
}

func testAccCheckSakuraCloudDatabaseDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("could not find Database data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("ID is not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudDatabaseDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("found Database data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudDatabaseDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_database" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		dbOp := sacloud.NewDatabaseOp(client)
		zone := rs.Primary.Attributes["zone"]
		_, err := dbOp.Read(context.Background(), zone, types.StringID(rs.Primary.ID))

		if err == nil {
			return errors.New("database still exists")
		}
	}

	return nil
}

func testAccCheckSakuraCloudDataSourceDatabaseBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "sw" {
  name = "%s"
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
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceDatabaseConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "sw" {
  name = "%s"
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
}
data "sakuracloud_database" "foobar" {
  filters {
	names = ["%s"]
  }
}`, name, name, name)
}

func testAccCheckSakuraCloudDataSourceDatabaseConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "sw" {
  name = "%s"
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
}
data "sakuracloud_database" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceDatabaseConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "sw" {
  name = "%s"
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
}
data "sakuracloud_database" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceDatabaseConfig_NotExists(name string) string {

	return fmt.Sprintf(`
resource "sakuracloud_switch" "sw" {
  name = "%s"
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
}
data "sakuracloud_database" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, name, name)
}
