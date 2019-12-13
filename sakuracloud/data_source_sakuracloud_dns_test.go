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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudDataSourceDNS_basic(t *testing.T) {
	resourceName := "data.sakuracloud_dns.foobar"
	zone := fmt.Sprintf("%s.com", randomName())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceDNS_basic, zone),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "zone", zone),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "tag3"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceDNS_basic = `
resource "sakuracloud_dns" "foobar" {
  zone = "{{ .arg0 }}"
  description = "description"
  tags = ["tag1", "tag2", "tag3"]
}

data "sakuracloud_dns" "foobar" {
  filters {
    names = [sakuracloud_dns.foobar.zone]
  }
}`
