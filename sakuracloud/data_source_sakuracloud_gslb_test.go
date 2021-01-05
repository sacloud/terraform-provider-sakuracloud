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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudDataSourceGSLB_basic(t *testing.T) {
	resourceName := "data.sakuracloud_gslb.foobar"
	rand := randomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceGSLB_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.delay_loop", "10"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.host_header", "usacloud.jp"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.path", "/"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.status", "200"),
					resource.TestCheckResourceAttr(resourceName, "sorry_server", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.4151227546", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1852302624", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "tags.425776566", "tag3"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceGSLB_basic = `
resource "sakuracloud_gslb" "foobar" {
  name = "{{ .arg0 }}"
  health_check {
    protocol    = "http"
    delay_loop  = 10
    host_header = "usacloud.jp"
    path        = "/"
    status      = "200"
  }
  sorry_server = "8.8.8.8"
  description  = "description"
  tags         = ["tag1", "tag2", "tag3"]
}

data "sakuracloud_gslb" "foobar" {
  filter {
	names = [sakuracloud_gslb.foobar.name]
  }
}`
