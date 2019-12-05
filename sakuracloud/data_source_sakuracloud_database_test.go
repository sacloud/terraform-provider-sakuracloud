// Copyright 2016-2019 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSakuraCloudDataSourceDatabase_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudDatabaseDestroy,

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
