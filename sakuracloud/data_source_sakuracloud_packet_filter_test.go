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

func TestAccSakuraCloudDataSourcePacketFilter_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudPacketFilterDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourcePacketFilterBase(name),
				Check:  testAccCheckSakuraCloudDataSourceExists("sakuracloud_packet_filter.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePacketFilterConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_packet_filter.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.#", "2"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.source_network", "0.0.0.0"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.source_port", "0-65535"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.destination_port", "80"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.allow", "true"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePacketFilterConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_packet_filter.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudDataSourcePacketFilterBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_packet_filter" "foobar" {
  name = "%s"
  description = "description_test"
  expressions {
  	protocol = "tcp"
  	source_network = "0.0.0.0"
  	source_port = "0-65535"
  	destination_port = "80"
  	allow = true
  }
  expressions {
  	protocol = "udp"
  	source_network = "0.0.0.0"
  	source_port = "0-65535"
  	destination_port = "80"
  	allow = true
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourcePacketFilterConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_packet_filter" "foobar" {
  name = "%s"
  description = "description_test"
  expressions {
  	protocol = "tcp"
  	source_network = "0.0.0.0"
  	source_port = "0-65535"
  	destination_port = "80"
  	allow = true
  }
  expressions {
  	protocol = "udp"
  	source_network = "0.0.0.0"
  	source_port = "0-65535"
  	destination_port = "80"
  	allow = true
  }
}
data "sakuracloud_packet_filter" "foobar" {
  filters {
	names = ["%s"]
  }
}`, name, name)
}

var testAccCheckSakuraCloudDataSourcePacketFilterConfig_NotExists = `
data "sakuracloud_packet_filter" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`
