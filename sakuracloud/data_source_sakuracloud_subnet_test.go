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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudSubnetDataSource_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudSubnetDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceSubnetBase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("sakuracloud_subnet.foobar"),
					testAccCheckSakuraCloudDataSourceExists("sakuracloud_subnet.foobar2"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_subnet.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_subnet.foobar", "ipaddresses.#", "16"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSubnetConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_subnet.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

var testAccCheckSakuraCloudDataSourceSubnetBase = `
resource sakuracloud_internet "foobar" {
  name = "subnet_test"
}
resource "sakuracloud_subnet" "foobar" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  next_hop = "${sakuracloud_internet.foobar.ipaddresses[0]}"
}
resource "sakuracloud_subnet" "foobar2" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  next_hop = "${sakuracloud_internet.foobar.ipaddresses[1]}"
}
`

var testAccCheckSakuraCloudDataSourceSubnetConfig = `
resource sakuracloud_internet "foobar" {
  name = "subnet_test"
}
resource "sakuracloud_subnet" "foobar" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  next_hop = "${sakuracloud_internet.foobar.ipaddresses[0]}"
}
resource "sakuracloud_subnet" "foobar2" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  next_hop = "${sakuracloud_internet.foobar.ipaddresses[1]}"
}

data sakuracloud_subnet "foobar" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  index = 1
}
`

var testAccCheckSakuraCloudDataSourceSubnetConfig_NotExists = `
resource sakuracloud_internet "foobar" {
  name = "subnet_test"
}
resource "sakuracloud_subnet" "foobar" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  next_hop = "${sakuracloud_internet.foobar.ipaddresses[0]}"
}
resource "sakuracloud_subnet" "foobar2" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  next_hop = "${sakuracloud_internet.foobar.ipaddresses[1]}"
}
data sakuracloud_subnet "foobar" {
  internet_id = "${sakuracloud_internet.foobar.id}"
  index = 2
}
`
