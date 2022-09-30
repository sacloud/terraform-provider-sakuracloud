// Copyright 2016-2022 terraform-provider-sakuracloud authors
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

func TestAccSakuraCloudDataSourceArchive_osType(t *testing.T) {
	resourceName := "data.sakuracloud_archive.foobar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceArchive_osType,
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
				),
			},
		},
	})
}

func TestAccSakuraCloudDataSourceArchive_withTag(t *testing.T) {
	resourceName := "data.sakuracloud_archive.foobar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceArchive_withTag,
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
				),
			},
		},
	})
}

var testAccCheckSakuraCloudDataSourceArchive_withTag = `
data "sakuracloud_archive" "foobar" {
  filter {
    tags = ["distro-ubuntu","os-linux"]
  }
}`

var testAccCheckSakuraCloudDataSourceArchive_osType = `
data "sakuracloud_archive" "foobar" {
    os_type = "ubuntu"
}
`
