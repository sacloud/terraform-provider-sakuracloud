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

func TestAccSakuraCloudDataSourceServer_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudServerDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceServerBase(name),
				Check:  testAccCheckSakuraCloudDataSourceExists("sakuracloud_server.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServerConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_server.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "interface_driver", "virtio"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.2", "tag3"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "core", "1"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "memory", "1"),
					//resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "disks.#", "1"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "interfaces.0.upstream", "shared"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "interfaces.#", "1"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServerConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_server.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServerConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_server.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServerConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_server.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudDataSourceServerBase(name string) string {
	return fmt.Sprintf(`
data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu"
}
resource "sakuracloud_disk" "foobar" {
  name              = "%s"
  source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}
resource "sakuracloud_server" "foobar" {
  name        = "%s"
  disks       = ["${sakuracloud_disk.foobar.id}"]
  description = "description_test"
  tags        = ["tag1", "tag2", "tag3"]
  interfaces {
    upstream = "shared"
  }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceServerConfig(name string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
  filters {
	names = ["%s"]
  }
}`, testAccCheckSakuraCloudDataSourceServerBase(name), name)
}

func testAccCheckSakuraCloudDataSourceServerConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
}`, testAccCheckSakuraCloudDataSourceServerBase(name))
}

func testAccCheckSakuraCloudDataSourceServerConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
}`, testAccCheckSakuraCloudDataSourceServerBase(name))
}

func testAccCheckSakuraCloudDataSourceServerConfig_NotExists(name string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, testAccCheckSakuraCloudDataSourceServerBase(name))
}
