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

func TestAccSakuraCloudDataSourceCDROM_basic(t *testing.T) {
	resourceName := "data.sakuracloud_cdrom.foobar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSakuraCloudDataSourceCDROM_basic,
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "Parted Magic 2013_08_01"),
					resource.TestCheckResourceAttr(resourceName, "size", "5"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "tags.2816018188", "arch-64bit"),
					resource.TestCheckResourceAttr(resourceName, "tags.1883356667", "current-stable"),
					resource.TestCheckResourceAttr(resourceName, "tags.3295801924", "distro-parted_magic"),
					resource.TestCheckResourceAttr(resourceName, "tags.2388413417", "distro-ver-2013.08.01"),
					resource.TestCheckResourceAttr(resourceName, "tags.1583874418", "os-linux"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceCDROM_basic = `
data "sakuracloud_cdrom" "foobar" {
  filter {
    condition {
	  name    = "Name"
	  values = ["Parted Magic 2013_08_01"]
    }
  }
}`
