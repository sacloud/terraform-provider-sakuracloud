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
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSakuraCloudDataSourceArchive_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudArchiveDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceArchiveConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceID("data.sakuracloud_archive.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "name", "Ubuntu Server 16.04.6 LTS 64bit"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "size", "20"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "zone", "tk1v"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.#", "5"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.0", "@size-extendable"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.1", "arch-64bit"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.2", "distro-ubuntu"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.3", "distro-ver-16.04.5"),
					resource.TestCheckResourceAttr("data.sakuracloud_archive.foobar", "tags.4", "os-linux"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceArchive_OSType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceID("data.sakuracloud_archive.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceArchiveConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceID("data.sakuracloud_archive.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceArchive_NameSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceID("data.sakuracloud_archive.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceArchive_TagSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceID("data.sakuracloud_archive.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceArchiveConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceNotExists("data.sakuracloud_archive.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceArchiveConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceNotExists("data.sakuracloud_archive.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceArchive_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceNotExists("data.sakuracloud_archive.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceArchive_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudArchiveDataSourceNotExists("data.sakuracloud_archive.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudArchiveDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Archive data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Archive data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudArchiveDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found Archive data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudArchiveDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	originalZone := client.Zone
	client.Zone = "tk1v"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_archive" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Archive.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Archive still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceArchiveConfig = `
data "sakuracloud_archive" "foobar" {
  filter {
    name = "Name"
    values = ["Ubuntu Server 16"]
  }
  zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchiveConfig_With_Tag = `
data "sakuracloud_archive" "foobar" {
  filter {
    name = "Tags"
    values = ["distro-ubuntu","os-linux"]
  }
  zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchiveConfig_With_NotExists_Tag = `
data "sakuracloud_archive" "foobar" {
  filter {
    name = "Tags"
    values = ["distro-ubuntu-xxxxxxxxxxx","os-linux-xxxxxxxx"]
  }
  zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchiveConfig_NotExists = `
data "sakuracloud_archive" "foobar" {
  filter {
    name = "Name"
    values = ["xxxxxxxxxxxxxxxxxx"]
  }
  zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchive_OSType = `
data "sakuracloud_archive" "foobar" {
    os_type = "rancheros"
    zone = "tk1v"
}
`

var testAccCheckSakuraCloudDataSourceArchive_NameSelector_Exists = `
data "sakuracloud_archive" "foobar" {
    name_selectors = ["Ubuntu", "Server","16"]
    zone = "tk1v"
}
`
var testAccCheckSakuraCloudDataSourceArchive_NameSelector_NotExists = `
data "sakuracloud_archive" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
    zone = "tk1v"
}
`

var testAccCheckSakuraCloudDataSourceArchive_TagSelector_Exists = `
data "sakuracloud_archive" "foobar" {
	tag_selectors = ["distro-ubuntu","os-linux"]
    zone = "tk1v"
}`

var testAccCheckSakuraCloudDataSourceArchive_TagSelector_NotExists = `
data "sakuracloud_archive" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
    zone = "tk1v"
}`
