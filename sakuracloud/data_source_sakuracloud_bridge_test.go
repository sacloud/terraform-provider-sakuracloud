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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudDataSourceBridge_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudBridgeDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeBase(name),
				Check:  testAccCheckSakuraCloudDataSourceExists("sakuracloud_bridge.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_bridge.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "description", "description_test"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeConfig_NameSelector_Exists(name, randString1, randString2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_bridge.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_bridge.foobar", "description", "description_test"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_bridge.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceBridgeConfig_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_bridge.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudDataSourceBridgeBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_bridge" "foobar" {
  name        = "%s"
  description = "description_test"
} 
`, name)
}

func testAccCheckSakuraCloudDataSourceBridgeConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_bridge" "foobar" {
  name        = "%s"
  description = "description_test"
}
data "sakuracloud_bridge" "foobar" {
  filters {
    conditions {
      name = "Name"
      values = ["%s"]
    }
  }
}`, name, name)
}

var testAccCheckSakuraCloudDataSourceBridgeConfig_NotExists = `
data "sakuracloud_bridge" "foobar" {
  filters {
    conditions {
      name = "Name"
      values = ["xxxxxxxxxxxxxxxxxx"]
    }
  }
}`

func testAccCheckSakuraCloudDataSourceBridgeConfig_NameSelector_Exists(name, p1, p2 string) string {
	return fmt.Sprintf(`
resource "sakuracloud_bridge" "foobar" {
  name = "%s"
  description = "description_test"
}
data "sakuracloud_bridge" "foobar" {
  filters {
    names = ["%s", "%s"]
  }
}`, name, p1, p2)
}

var testAccCheckSakuraCloudDataSourceBridgeConfig_NameSelector_NotExists = `
data "sakuracloud_bridge" "foobar" {
  filters {
    names = ["xxxxxxxx"]
  }
}`
