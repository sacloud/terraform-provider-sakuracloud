// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSakuraCloudDataSourceDatabase_basic(t *testing.T) {
	resourceName := "data.sakuracloud_database.foobar"
	rand := randomName()
	password := randomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceDatabase_basic, rand, password),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "plan", "10g"),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "tag3"),
					resource.TestCheckResourceAttr(resourceName, "backup.0.weekdays.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "backup.0.weekdays.0", "mon"),
					resource.TestCheckResourceAttr(resourceName, "backup.0.weekdays.1", "tue"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceDatabase_basic = `
resource "sakuracloud_switch" "foobar" {
  name = "{{ .arg0 }}"
}

resource "sakuracloud_database" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2", "tag3"]

  username = "defuser"
  password = "{{ .arg1 }}"

  network_interface {
    switch_id       = "${sakuracloud_switch.foobar.id}"
    ip_address      = "192.168.101.101"
    netmask         = 24
    gateway         = "192.168.101.1"
    port            = 54321
    source_ranges   = ["192.168.101.0/24", "192.168.102.0/24"]
  }
  backup {
    weekdays = ["mon", "tue"]
    time     = "00:00"
  }
}

data "sakuracloud_database" "foobar" {
  filter {
    names = [sakuracloud_database.foobar.name]
  }
}`
