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

func TestAccSakuraCloudDataSourceNote_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudNoteDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceNoteBase(name),
				Check:  testAccCheckSakuraCloudDataSourceExists("sakuracloud_note.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_note.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "content", "content_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "class", "shell"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_note.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_note.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_note.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudDataSourceNoteBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_note" "foobar" {
  name    = "%s"
  content = "content_test"
  tags    = ["tag1", "tag2", "tag3"]
}`, name)
}

func testAccCheckSakuraCloudDataSourceNoteConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_note" "foobar" {
  name    = "%s"
  content = "content_test"
  tags    = ["tag1", "tag2", "tag3"]
}

data "sakuracloud_note" "foobar" {
  filters {
	names = ["%s"]
  }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceNoteConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_note" "foobar" {
  name    = "%s"
  content = "content_test"
  tags    = ["tag1", "tag2", "tag3"]
}

data "sakuracloud_note" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourceNoteConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_note" "foobar" {
  name    = "%s"
  content = "content_test"
  tags    = ["tag1", "tag2", "tag3"]
}

data "sakuracloud_note" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourceNoteConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_note" "foobar" {
  name    = "%s"
  content = "content_test"
  tags    = ["tag1", "tag2", "tag3"]
}

data "sakuracloud_note" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, name)
}
