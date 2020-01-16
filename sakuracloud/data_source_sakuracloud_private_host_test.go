// Copyright 2016-2020 terraform-provider-sakuracloud authors
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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSakuraCloudDataSourcePrivateHost_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudPrivateHostDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostBase(name),
				Check:  testAccCheckSakuraCloudPrivateHostDataSourceID("sakuracloud_private_host.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceID("data.sakuracloud_private_host.foobar"),
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
					testAccCheckSakuraCloudPrivateHostDataSourceID("data.sakuracloud_private_host.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHost_NameSelector_Exists(name, randString1, randString2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceID("data.sakuracloud_private_host.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHost_TagSelector_Exists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceID("data.sakuracloud_private_host.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceNotExists("data.sakuracloud_private_host.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceNotExists("data.sakuracloud_private_host.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHost_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceNotExists("data.sakuracloud_private_host.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHost_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceNotExists("data.sakuracloud_private_host.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudPrivateHostDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find PrivateHost data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("PrivateHost data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudPrivateHostDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found PrivateHost data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudPrivateHostDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	originalZone := client.Zone
	client.Zone = "tk1a"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_private_host" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.PrivateHost.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("PrivateHost still exists")
		}
	}

	return nil
}

func testAccCheckSakuraCloudDataSourcePrivateHostBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_private_host" "foobar" {
    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}`, name)
}

func testAccCheckSakuraCloudDataSourcePrivateHostConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_private_host" "foobar" {
    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
data "sakuracloud_private_host" "foobar" {
    filter {
	name = "Name"
	values = ["%s"]
    }
    zone = "tk1a"
}`, name, name)
}

func testAccCheckSakuraCloudDataSourcePrivateHostConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_private_host" "foobar" {
    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
data "sakuracloud_private_host" "foobar" {
    filter {
	name = "Tags"
	values = ["tag1","tag3"]
    }
    zone = "tk1a"
}`, name)
}

func testAccCheckSakuraCloudDataSourcePrivateHostConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_private_host" "foobar" {
    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
data "sakuracloud_private_host" "foobar" {
    filter {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
    zone = "tk1a"
}`, name)
}

func testAccCheckSakuraCloudDataSourcePrivateHostConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_private_host" "foobar" {
    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
data "sakuracloud_private_host" "foobar" {
    filter {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
    zone = "tk1a"
}`, name)
}

func testAccCheckSakuraCloudDataSourcePrivateHost_NameSelector_Exists(name, p1, p2 string) string {
	return fmt.Sprintf(`
resource "sakuracloud_private_host" "foobar" {
    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
data "sakuracloud_private_host" "foobar" {
    name_selectors = ["%s", "%s"]
    zone = "tk1a"
}`, name, p1, p2)
}

var testAccCheckSakuraCloudDataSourcePrivateHost_NameSelector_NotExists = `
data "sakuracloud_private_host" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
    zone = "tk1a"
}
`

func testAccCheckSakuraCloudDataSourcePrivateHost_TagSelector_Exists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_private_host" "foobar" {
    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    zone = "tk1a"
}
data "sakuracloud_private_host" "foobar" {
	tag_selectors = ["tag1","tag2","tag3"]
    zone = "tk1a"
}`, name)
}

var testAccCheckSakuraCloudDataSourcePrivateHost_TagSelector_NotExists = `
data "sakuracloud_private_host" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
    zone = "tk1a"
}`
