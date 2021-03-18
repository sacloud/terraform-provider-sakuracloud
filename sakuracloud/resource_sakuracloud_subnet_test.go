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
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func TestAccSakuraCloudSubnet_basic(t *testing.T) {
	resourceName := "sakuracloud_subnet.foobar"
	rand := randomName()

	var subnet sacloud.Subnet
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudInternetDestroy,
			testCheckSakuraCloudSubnetDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSubnet_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudSubnetExists(resourceName, &subnet),
					resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "16"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSubnet_update, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudSubnetExists(resourceName, &subnet),
					resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "16"),
				),
			},
		},
	})
}

func testCheckSakuraCloudSubnetExists(n string, subnet *sacloud.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no Subnet ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		subnetOp := sacloud.NewSubnetOp(client)
		zone := rs.Primary.Attributes["zone"]

		foundSubnet, err := subnetOp.Read(context.Background(), zone, sakuraCloudID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if foundSubnet.ID.String() != rs.Primary.ID {
			return fmt.Errorf("not found Subnet: %s", rs.Primary.ID)
		}

		*subnet = *foundSubnet
		return nil
	}
}

func testCheckSakuraCloudSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	subnetOp := sacloud.NewSubnetOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_subnet" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		zone := rs.Primary.Attributes["zone"]
		_, err := subnetOp.Read(context.Background(), zone, sakuraCloudID(rs.Primary.ID))
		if err == nil {
			return fmt.Errorf("still exists Subnet: %s", rs.Primary.ID)
		}
	}

	return nil
}

var testAccSakuraCloudSubnet_basic = `
resource sakuracloud_internet "foobar" {
  name = "{{ .arg0 }}"
}
resource "sakuracloud_subnet" "foobar" {
  internet_id = sakuracloud_internet.foobar.id
  next_hop    = sakuracloud_internet.foobar.min_ip_address
}`

var testAccSakuraCloudSubnet_update = `
resource sakuracloud_internet "foobar" {
  name = "{{ .arg0 }}"
}
resource "sakuracloud_subnet" "foobar" {
  internet_id = sakuracloud_internet.foobar.id
  next_hop    = sakuracloud_internet.foobar.max_ip_address
}`
