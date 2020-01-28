// Copyright 2016-2020 terraform-provider-sakuracloud authors
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/sacloud"
)

func TestAccResourceSakuraCloudSubnet_basic(t *testing.T) {
	var subnet sacloud.Subnet
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudSubnetConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSubnetExists("sakuracloud_subnet.foobar", &subnet),
					resource.TestCheckResourceAttr(
						"sakuracloud_subnet.foobar", "ipaddresses.#", "16"),
				),
			},
			{
				Config: testAccCheckSakuraCloudSubnetConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSubnetExists("sakuracloud_subnet.foobar", &subnet),
					resource.TestCheckResourceAttr(
						"sakuracloud_subnet.foobar", "ipaddresses.#", "16"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudSubnetExists(n string, subnet *sacloud.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Subnet ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)

		foundSubnet, err := client.Subnet.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundSubnet.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("Subnet not found")
		}

		*subnet = *foundSubnet

		return nil
	}
}

func testAccCheckSakuraCloudSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_subnet" {
			continue
		}

		_, err := client.Subnet.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Subnet still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudSubnetConfig_basic = `
resource sakuracloud_internet "foobar" {
    name = "myinternet"
}
resource "sakuracloud_subnet" "foobar" {
    internet_id = "${sakuracloud_internet.foobar.id}"
    next_hop = "${sakuracloud_internet.foobar.min_ipaddress}"
}`

var testAccCheckSakuraCloudSubnetConfig_update = `
resource sakuracloud_internet "foobar" {
    name = "myinternet"
}
resource "sakuracloud_subnet" "foobar" {
    internet_id = "${sakuracloud_internet.foobar.id}"
    next_hop = "${sakuracloud_internet.foobar.max_ipaddress}"
}`
