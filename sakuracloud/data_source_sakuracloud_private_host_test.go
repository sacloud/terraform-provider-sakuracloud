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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudDataSourcePrivateHost_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudPrivateHostDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostBase(name),
				Check:  testAccCheckSakuraCloudDataSourceExists("sakuracloud_private_host.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_private_host.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "tags.2", "tag3"),
					resource.TestMatchResourceAttr("data.sakuracloud_private_host.foobar",
						"hostname",
						regexp.MustCompile(".+")), // should be not empty
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_private_host.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_private_host.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_private_host.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudDataSourcePrivateHostBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_private_host" "foobar" {
  name        = "%s"
  description = "description_test"
  tags        = ["tag1", "tag2", "tag3"]
  zone        = "tk1a"
}`, name)
}

func testAccCheckSakuraCloudDataSourcePrivateHostConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_private_host" "foobar" {
  name        = "%s"
  description = "description_test"
  tags        = ["tag1", "tag2", "tag3"]
  zone        = "tk1a"
}

data "sakuracloud_private_host" "foobar" {
  filters {
	names = ["%s"]
  }
  zone = "tk1a"
}`, name, name)
}

func testAccCheckSakuraCloudDataSourcePrivateHostConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_private_host" "foobar" {
  name        = "%s"
  description = "description_test"
  tags        = ["tag1", "tag2", "tag3"]
  zone        = "tk1a"
}

data "sakuracloud_private_host" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
  zone = "tk1a"
}`, name)
}

func testAccCheckSakuraCloudDataSourcePrivateHostConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_private_host" "foobar" {
  name        = "%s"
  description = "description_test"
  tags        = ["tag1", "tag2", "tag3"]
  zone        = "tk1a"
}

data "sakuracloud_private_host" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
  zone = "tk1a"
}`, name)
}

func testAccCheckSakuraCloudDataSourcePrivateHostConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_private_host" "foobar" {
  name        = "%s"
  description = "description_test"
  tags        = ["tag1", "tag2", "tag3"]
  zone        = "tk1a"
}

data "sakuracloud_private_host" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
  zone = "tk1a"
}`, name)
}
