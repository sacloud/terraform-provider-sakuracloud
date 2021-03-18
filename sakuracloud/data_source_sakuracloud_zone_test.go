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

func TestAccSakuraCloudDataSourceZone_basic(t *testing.T) {
	resourceName := "data.sakuracloud_zone.foobar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSakuraCloudDataSourceZone_basic,
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "is1a"),
					resource.TestCheckResourceAttr(resourceName, "zone_id", "31001"),
					resource.TestCheckResourceAttr(resourceName, "description", "石狩第1ゾーン"),
					resource.TestCheckResourceAttr(resourceName, "region_id", "310"),
					resource.TestCheckResourceAttr(resourceName, "region_name", "石狩"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "133.242.0.3"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "133.242.0.4"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceZone_basic = `
data "sakuracloud_zone" "foobar" { 
  name = "is1a"
}`
