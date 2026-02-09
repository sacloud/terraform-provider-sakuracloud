// Copyright 2016-2025 terraform-provider-sakuracloud authors
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

func TestAccSakuraCloudDataSourceNFS_basic(t *testing.T) {
	resourceName := "data.sakuracloud_nfs.foobar"
	rand := randomName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceNFS_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "tag3"),
					resource.TestCheckResourceAttrPair(
						resourceName, "network_interface.0.switch_id",
						"sakuracloud_switch.foobar", "id",
					),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.ip_address", "192.168.11.101"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.gateway", "192.168.11.1"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceNFS_basic = `
resource sakuracloud_switch "foobar" {
  name = "{{ .arg0 }}"
}

resource "sakuracloud_nfs" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2", "tag3"]

  plan = "ssd"
  size = "500"

  network_interface {
    switch_id   = sakuracloud_switch.foobar.id
    ip_address  = "192.168.11.101"
    netmask     = 24
    gateway     = "192.168.11.1"
  }
}

data "sakuracloud_nfs" "foobar" {
  filter {
	names = [sakuracloud_nfs.foobar.name]
  }
}`
