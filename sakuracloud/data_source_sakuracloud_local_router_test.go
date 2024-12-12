// Copyright 2016-2023 terraform-provider-sakuracloud authors
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

func TestAccSakuraCloudDataSourceLocalRouter_basic(t *testing.T) {
	if !isFakeModeEnabled() {
		t.Skip("This test only run if FAKE_MODE environment variable is set")
	}

	resourceName := "data.sakuracloud_local_router.foobar"
	rand := randomName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceLocalRouter_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "tag3"),
					resource.TestCheckResourceAttrPair(
						resourceName, "switch.0.code",
						"sakuracloud_switch.foobar", "id"),
					resource.TestCheckResourceAttr(resourceName, "switch.0.category", "cloud"),
					resource.TestCheckResourceAttrPair(
						resourceName, "switch.0.zone_id",
						"data.sakuracloud_zone.current", "name"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.vip", "192.168.11.1"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.ip_addresses.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.ip_addresses.0", "192.168.11.11"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.ip_addresses.1", "192.168.11.12"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.vrid", "1"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceLocalRouter_basic = `
resource sakuracloud_switch "foobar" {
  name = "{{ .arg0 }}"
}

data sakuracloud_zone "current" {}

resource "sakuracloud_local_router" "foobar" {
  switch {
    code     = sakuracloud_switch.foobar.id
    category = "cloud"
    zone_id  = data.sakuracloud_zone.current.name
  }
  network_interface {
    vip          = "192.168.11.1"
    ip_addresses = ["192.168.11.11", "192.168.11.12"]
    netmask      = 24
    vrid         = 1
  }

  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2", "tag3"]
}

data "sakuracloud_local_router" "foobar" {
  filter {
	names = [sakuracloud_local_router.foobar.name]
  }
}`
