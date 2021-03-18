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

func TestAccSakuraCloudDataSourceArchive_basic(t *testing.T) {
	resourceName := "data.sakuracloud_archive.foobar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSakuraCloudDataSourceArchive_basic,
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "Ubuntu Server 16.04.6 LTS 64bit"),
					resource.TestCheckResourceAttr(resourceName, "size", "20"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "tags.1695116635", "@size-extendable"),
					resource.TestCheckResourceAttr(resourceName, "tags.2816018188", "arch-64bit"),
					resource.TestCheckResourceAttr(resourceName, "tags.1490716481", "distro-ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "tags.4143804920", "distro-ver-16.04.5"),
					resource.TestCheckResourceAttr(resourceName, "tags.1583874418", "os-linux"),
					resource.TestCheckResourceAttr(resourceName, "tags.1550871325", "ubuntu-16.04-latest"),
				),
			},
		},
	})
}

func TestAccSakuraCloudDataSourceArchive_osType(t *testing.T) {
	resourceName := "data.sakuracloud_archive.foobar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

func TestAccSakuraCloudDataSourceArchive_notExists(t *testing.T) {
	name := "data.sakuracloud_archive.foobar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceArchive_notExists,
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceNotExists(name),
				),
				Destroy: true,
			},
		},
	})
}

func TestAccSakuraCloudDataSourceArchive_tagNotExists(t *testing.T) {
	name := "data.sakuracloud_archive.foobar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceArchive_tagNotExists,
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceNotExists(name),
				),
				Destroy: true,
			},
		},
	})
}

var testAccSakuraCloudDataSourceArchive_basic = `
data "sakuracloud_archive" "foobar" {
  filter {
    names = ["Ubuntu Server 16"]
  }
}`

var testAccCheckSakuraCloudDataSourceArchive_withTag = `
data "sakuracloud_archive" "foobar" {
  filter {
    tags = ["distro-ubuntu","os-linux"]
  }
}`

var testAccCheckSakuraCloudDataSourceArchive_tagNotExists = `
data "sakuracloud_archive" "foobar" {
  filter {
    tags = ["distro-ubuntu-xxxxxxxxxxx","os-linux-xxxxxxxx"]
  }
}`

var testAccCheckSakuraCloudDataSourceArchive_notExists = `
data "sakuracloud_archive" "foobar" {
  filter {
    names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`

var testAccCheckSakuraCloudDataSourceArchive_osType = `
data "sakuracloud_archive" "foobar" {
    os_type = "rancheros"
}
`
