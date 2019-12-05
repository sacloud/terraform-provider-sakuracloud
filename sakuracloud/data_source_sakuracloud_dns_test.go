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

func TestAccSakuraCloudDataSourceDNS_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	zone := fmt.Sprintf("%s.%s.com", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudDNSDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceDNSBase(zone),
				Check:  testAccCheckSakuraCloudDataSourceExists("sakuracloud_dns.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDNSConfig(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_dns.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "zone", zone),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDNSConfig_With_Tag(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_dns.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDNSConfig_NotExists(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_dns.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDNSConfig_With_NotExists_Tag(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_dns.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudDataSourceDNSBase(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}`, zone)
}

func testAccCheckSakuraCloudDataSourceDNSConfig(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_dns" "foobar" {
  filters {
	names = ["%s"]
  }
}`, zone, zone)
}

func testAccCheckSakuraCloudDataSourceDNSConfig_With_Tag(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_dns" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
}`, zone)
}

func testAccCheckSakuraCloudDataSourceDNSConfig_With_NotExists_Tag(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_dns" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
}`, zone)
}

func testAccCheckSakuraCloudDataSourceDNSConfig_NotExists(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_dns" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, zone)
}
