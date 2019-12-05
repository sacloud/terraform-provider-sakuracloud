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

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSakuraCloudDataSourceInternet_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudInternetDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceInternetBase(name),
				Check:  testAccCheckSakuraCloudInternetDataSourceID("sakuracloud_internet.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternetConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceID("data.sakuracloud_internet.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "tags.2", "tag3"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "nw_mask_len", "28"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "band_width", "100"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "server_ids.#", "0"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "ipaddresses.#", "11"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternetConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceID("data.sakuracloud_internet.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternetConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceNotExists("data.sakuracloud_internet.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternetConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceNotExists("data.sakuracloud_internet.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudInternetDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Internet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Internet data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudInternetDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found Internet data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudDataSourceInternetBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_internet" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}`, name)
}

func testAccCheckSakuraCloudDataSourceInternetConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_internet" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}

data "sakuracloud_internet" "foobar" {
  filters {
	names = ["%s"]
  }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceInternetConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_internet" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourceInternetConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_internet" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourceInternetConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_internet" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, name)
}
