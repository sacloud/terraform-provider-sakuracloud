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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSakuraCloudDataSourceZone_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceZoneBase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudZoneDataSourceID("data.sakuracloud_zone.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "name", "is1a"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "zone_id", "31001"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "description", "石狩第1ゾーン"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "region_id", "310"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "region_name", "石狩"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "dns_servers.0", "133.242.0.3"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "dns_servers.1", "133.242.0.4"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudZoneDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Zone data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Zone data source ID not set")
		}
		return nil
	}
}

var testAccCheckSakuraCloudDataSourceZoneBase = `
data "sakuracloud_zone" "foobar" { 
  name = "is1a"
}`
