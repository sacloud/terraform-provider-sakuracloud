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

func TestAccSakuraCloudSubnetDataSource_basic(t *testing.T) {
	resourceName := "data.sakuracloud_subnet.foobar"
	rand := randomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSubnetDataSource_pre, rand),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSubnetDataSource_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "16"),
				),
				Destroy: true,
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSubnetDataSource_notExists, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceNotExists(resourceName),
				),
				Destroy: true,
			},
		},
	})
}

var testAccSakuraCloudSubnetDataSource_pre = `
resource sakuracloud_internet "foobar" {
  name = "{{ .arg0 }}"
}
resource "sakuracloud_subnet" "foobar" {
  internet_id = sakuracloud_internet.foobar.id
  next_hop    = sakuracloud_internet.foobar.ip_addresses[0]
}
resource "sakuracloud_subnet" "foobar2" {
  internet_id = sakuracloud_internet.foobar.id
  next_hop    = sakuracloud_internet.foobar.ip_addresses[1]
}
`

var testAccSakuraCloudSubnetDataSource_basic = `
resource sakuracloud_internet "foobar" {
  name = "{{ .arg0 }}"
}
resource "sakuracloud_subnet" "foobar" {
  internet_id = sakuracloud_internet.foobar.id
  next_hop    = sakuracloud_internet.foobar.ip_addresses[0]
}
resource "sakuracloud_subnet" "foobar2" {
  internet_id = sakuracloud_internet.foobar.id
  next_hop    = sakuracloud_internet.foobar.ip_addresses[1]
}

data sakuracloud_subnet "foobar" {
  internet_id = sakuracloud_internet.foobar.id
  index       = 1
}
`

var testAccSakuraCloudSubnetDataSource_notExists = `
resource sakuracloud_internet "foobar" {
  name = "{{ .arg0 }}"
}
resource "sakuracloud_subnet" "foobar" {
  internet_id = sakuracloud_internet.foobar.id
  next_hop    = sakuracloud_internet.foobar.ip_addresses[0]
}
resource "sakuracloud_subnet" "foobar2" {
  internet_id = sakuracloud_internet.foobar.id
  next_hop    = sakuracloud_internet.foobar.ip_addresses[1]
}

data sakuracloud_subnet "foobar" {
  internet_id = sakuracloud_internet.foobar.id
  index       = 2
}
`
