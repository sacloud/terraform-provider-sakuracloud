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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSakuraCloudDataSourceGSLB_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudGSLBDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceGSLBBase(name),
				Check:  testAccCheckSakuraCloudGSLBDataSourceID("sakuracloud_gslb.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceGSLBConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBDataSourceID("data.sakuracloud_gslb.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "sorry_server", "8.8.8.8"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "health_check.0.protocol", "http"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "health_check.0.delay_loop", "10"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "health_check.0.host_header", "terraform.io"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceGSLBConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBDataSourceID("data.sakuracloud_gslb.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceGSLB_NameSelector_Exists(name, randString1, randString2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBDataSourceID("data.sakuracloud_gslb.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceGSLB_TagSelector_Exists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBDataSourceID("data.sakuracloud_gslb.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceGSLBConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBDataSourceNotExists("data.sakuracloud_gslb.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceGSLBConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBDataSourceNotExists("data.sakuracloud_gslb.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceGSLB_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBDataSourceNotExists("data.sakuracloud_gslb.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceGSLB_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBDataSourceNotExists("data.sakuracloud_gslb.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudGSLBDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find GSLB data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("GSLB data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudGSLBDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found GSLB data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudGSLBDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_gslb" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.GSLB.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("GSLB still exists")
		}
	}

	return nil
}

func testAccCheckSakuraCloudDataSourceGSLBBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_gslb" "foobar" {
    name = "%s"
    health_check {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    sorry_server = "8.8.8.8"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}`, name)
}

func testAccCheckSakuraCloudDataSourceGSLBConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_gslb" "foobar" {
    name = "%s"
    health_check {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    sorry_server = "8.8.8.8"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_gslb" "foobar" {
    filter {
	name = "Name"
	values = ["%s"]
    }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceGSLBConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_gslb" "foobar" {
    name = "%s"
    health_check {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    sorry_server = "8.8.8.8"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_gslb" "foobar" {
    filter {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`, name)
}

func testAccCheckSakuraCloudDataSourceGSLBConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_gslb" "foobar" {
    name = "%s"
    health_check {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    sorry_server = "8.8.8.8"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_gslb" "foobar" {
    filter {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`, name)
}

func testAccCheckSakuraCloudDataSourceGSLBConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_gslb" "foobar" {
    name = "%s"
    health_check {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    sorry_server = "8.8.8.8"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_gslb" "foobar" {
    filter {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`, name)
}

func testAccCheckSakuraCloudDataSourceGSLB_NameSelector_Exists(name, p1, p2 string) string {
	return fmt.Sprintf(`
resource "sakuracloud_gslb" "foobar" {
    name = "%s"
    health_check {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    sorry_server = "8.8.8.8"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_gslb" "foobar" {
    name_selectors = ["%s", "%s"]
}
`, name, p1, p2)
}

var testAccCheckSakuraCloudDataSourceGSLB_NameSelector_NotExists = `
data "sakuracloud_gslb" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}
`

func testAccCheckSakuraCloudDataSourceGSLB_TagSelector_Exists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_gslb" "foobar" {
    name = "%s"
    health_check {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    sorry_server = "8.8.8.8"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_gslb" "foobar" {
	tag_selectors = ["tag1","tag2","tag3"]
}`, name)
}

var testAccCheckSakuraCloudDataSourceGSLB_TagSelector_NotExists = `
data "sakuracloud_gslb" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
}`
