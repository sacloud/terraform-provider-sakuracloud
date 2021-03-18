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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSakuraCloudDataSourceNote_Basic(t *testing.T) {
	resourceName := "data.sakuracloud_note.foobar"
	rand := randomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccCheckSakuraCloudDataSourceNoteConfig, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "content", "content"),
					resource.TestCheckResourceAttr(resourceName, "class", "shell"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.4151227546", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1852302624", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "tags.425776566", "tag3"),
				),
			},
		},
	})
}

var testAccCheckSakuraCloudDataSourceNoteConfig = `
resource "sakuracloud_note" "foobar" {
  name    = "{{ .arg0 }}"
  content = "content"
  tags    = ["tag1", "tag2", "tag3"]
}

data "sakuracloud_note" "foobar" {
  filter {
	names = [sakuracloud_note.foobar.name]
  }
}`
