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

func TestAccSakuraCloudDataSourceCDROM_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudCDROMDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceCDROMConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_cdrom.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "name", "Ubuntu Server 18.04.3 LTS 64bit"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "size", "5"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.#", "5"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.0", "arch-64bit"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.1", "current-stable"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.2", "distro-ubuntu"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.3", "distro-ver-18.04.3"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.4", "os-unix"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceCDROM_NameSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_cdrom.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceCDROM_TagSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_cdrom.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceCDROMConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_cdrom.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceCDROM_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_cdrom.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceCDROM_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_cdrom.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

var testAccCheckSakuraCloudDataSourceCDROMConfig = `
data "sakuracloud_cdrom" "foobar" {
  filters {
    conditions {
	  name = "Name"
	  values = ["Ubuntu Server 18.04.3 LTS 64bit"]
    }
  }
}`

var testAccCheckSakuraCloudDataSourceCDROMConfig_NotExists = `
data "sakuracloud_cdrom" "foobar" {
  filters {
    conditions {
	  name = "Name"
	  values = ["xxxxxxxxxxxxxxxxxx"]
    }
  }
}`

var testAccCheckSakuraCloudDataSourceCDROM_NameSelector_Exists = `
data "sakuracloud_cdrom" "foobar" {
  filters {
    names = ["Ubuntu","server","18"]
  }
}
`
var testAccCheckSakuraCloudDataSourceCDROM_NameSelector_NotExists = `
data "sakuracloud_cdrom" "foobar" {
  filters {
    names = ["xxxxxxxxxx"]
  }
}
`

var testAccCheckSakuraCloudDataSourceCDROM_TagSelector_Exists = `
data "sakuracloud_cdrom" "foobar" {
  filters {
	tags = ["distro-ubuntu","os-unix"]
  }
}`

var testAccCheckSakuraCloudDataSourceCDROM_TagSelector_NotExists = `
data "sakuracloud_cdrom" "foobar" {
  filters {
	tags = ["xxxxxxxxxx"]
  }
}`
