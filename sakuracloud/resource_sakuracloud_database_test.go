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
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func TestAccSakuraCloudDatabase_basic(t *testing.T) {
	resourceName := "sakuracloud_database.foobar"
	rand := randomName()
	password := randomPassword()

	var database sacloud.Database
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudDatabaseDestroy,
			testCheckSakuraCloudIconDestroy,
			testCheckSakuraCloudSwitchDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDatabase_basic, rand, password),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDatabaseExists(resourceName, &database),
					testCheckSakuraCloudDatabaseIsMaster(true, &database),
					resource.TestCheckResourceAttr(resourceName, "database_type", "mariadb"),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "plan", "30g"),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "username", "defuser"),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "replica_password", password),
					resource.TestCheckResourceAttr(resourceName, "source_ranges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "source_ranges.0", "192.168.11.0/24"),
					resource.TestCheckResourceAttr(resourceName, "source_ranges.1", "192.168.12.0/24"),
					resource.TestCheckResourceAttr(resourceName, "port", "33061"),
					resource.TestCheckResourceAttr(resourceName, "backup_time", "00:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_weekdays.0", "mon"),
					resource.TestCheckResourceAttr(resourceName, "backup_weekdays.1", "tue"),
					resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_addresses.0", "192.168.11.101"),
					resource.TestCheckResourceAttr(resourceName, "netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.11.1"),
					resource.TestCheckResourceAttrPair(
						resourceName, "icon_id",
						"sakuracloud_icon.foobar", "id",
					),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDatabase_update, rand, password),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDatabaseExists(resourceName, &database),
					testCheckSakuraCloudDatabaseIsMaster(false, &database),
					resource.TestCheckResourceAttr(resourceName, "database_type", "mariadb"),
					resource.TestCheckResourceAttr(resourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "plan", "30g"),
					resource.TestCheckResourceAttr(resourceName, "description", "description-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2-upd"),
					resource.TestCheckResourceAttr(resourceName, "username", "defuser"),
					resource.TestCheckResourceAttr(resourceName, "password", password+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "source_ranges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "source_ranges.0", "192.168.110.0/24"),
					resource.TestCheckResourceAttr(resourceName, "source_ranges.1", "192.168.120.0/24"),
					resource.TestCheckResourceAttr(resourceName, "port", "33062"),
					resource.TestCheckResourceAttr(resourceName, "backup_time", "00:30"),
					resource.TestCheckResourceAttr(resourceName, "backup_weekdays.0", "sun"),
					resource.TestCheckResourceAttr(resourceName, "backup_weekdays.1", "sat"),
					resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_addresses.0", "192.168.11.101"),
					resource.TestCheckResourceAttr(resourceName, "netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "gateway", "192.168.11.1"),
					resource.TestCheckResourceAttr(resourceName, "icon_id", ""),
				),
			},
		},
	})
}

func testCheckSakuraCloudDatabaseExists(n string, database *sacloud.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no Database ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		dbOp := sacloud.NewDatabaseOp(client)
		zone := rs.Primary.Attributes["zone"]

		foundDatabase, err := dbOp.Read(context.Background(), zone, sakuraCloudID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if foundDatabase.ID.String() != rs.Primary.ID {
			return fmt.Errorf("resource Database[%s] not found", rs.Primary.ID)
		}

		*database = *foundDatabase

		return nil
	}
}

func testCheckSakuraCloudDatabaseIsMaster(isMaster bool, database *sacloud.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if database == nil {
			return errors.New("database is nil")
		}

		dbStat := database.ReplicationSetting != nil && database.ReplicationSetting.Model == types.DatabaseReplicationModels.MasterSlave

		if dbStat != isMaster {
			return fmt.Errorf("database replication settings is not match, expect: %t", isMaster)
		}
		return nil
	}
}

func testCheckSakuraCloudDatabaseDestroy(s *terraform.State) error {
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
		_, err := dbOp.Read(context.Background(), zone, sakuraCloudID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("resource Database[%s] still exists", rs.Primary.ID)
		}
	}

	return nil
}

func TestAccImportSakuraCloudDatabase(t *testing.T) {
	name := randomName()
	password := randomPassword()

	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %#v", s)
		}
		expects := map[string]string{
			"name":              name,
			"database_type":     "mariadb",
			"description":       "description",
			"plan":              "30g",
			"username":          "defuser",
			"password":          password,
			"replica_password":  password,
			"source_ranges.0":   "192.168.11.0/24",
			"source_ranges.1":   "192.168.12.0/24",
			"port":              "33061",
			"backup_time":       "00:00",
			"backup_weekdays.0": "mon",
			"backup_weekdays.1": "tue",
			"ip_addresses.0":    "192.168.11.101",
			"netmask":           "24",
			"gateway":           "192.168.11.1",
			"tags.0":            "tag1",
			"tags.1":            "tag2",
		}

		if err := compareStateMulti(s[0], expects); err != nil {
			return err
		}
		return stateNotEmptyMulti(s[0], "icon_id", "switch_id")
	}

	resourceName := "sakuracloud_database.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudDatabaseDestroy,
			testCheckSakuraCloudIconDestroy,
			testCheckSakuraCloudSwitchDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDatabase_basic, name, password),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateCheck:  checkFn,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSakuraCloudDatabase_basic = `
resource "sakuracloud_switch" "foobar" {
  name = "{{ .arg0 }}"
}

resource "sakuracloud_database" "foobar" {
  database_type = "mariadb"
  plan          = "30g"

  username = "defuser"
  password = "{{ .arg1 }}"

  replica_password = "{{ .arg1 }}"

  source_ranges = ["192.168.11.0/24", "192.168.12.0/24"]

  port = 33061

  backup_time     = "00:00"
  backup_weekdays = ["mon", "tue"]

  switch_id    = sakuracloud_switch.foobar.id
  ip_addresses = ["192.168.11.101"]
  netmask      = 24
  gateway      = "192.168.11.1"

  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2"]
  icon_id     = sakuracloud_icon.foobar.id
}

resource "sakuracloud_icon" "foobar" {
  name          = "{{ .arg0 }}"
  base64content = "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"
}
`

const testAccSakuraCloudDatabase_update = `
resource "sakuracloud_switch" "foobar" {
  name = "{{ .arg0 }}"
}
resource "sakuracloud_database" "foobar" {
  database_type = "mariadb"

  plan     = "30g"
  username = "defuser"
  password = "{{ .arg1 }}-upd"

  source_ranges = ["192.168.110.0/24", "192.168.120.0/24"]

  port = 33062

  backup_time     = "00:30"
  backup_weekdays = ["sun", "sat"]

  name        = "{{ .arg0 }}-upd"
  description = "description-upd"
  tags        = ["tag1-upd", "tag2-upd"]

  switch_id    = sakuracloud_switch.foobar.id
  ip_addresses = ["192.168.11.101"]
  netmask      = 24
  gateway      = "192.168.11.1"
}`
