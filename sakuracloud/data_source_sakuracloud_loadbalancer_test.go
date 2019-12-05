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

func TestAccSakuraCloudDataSourceLoadBalancer_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudLoadBalancerDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceLoadBalancerBase(name),
				Check:  testAccCheckSakuraCloudLoadBalancerDataSourceID("sakuracloud_load_balancer.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceLoadBalancerConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerDataSourceID("data.sakuracloud_load_balancer.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_load_balancer.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceLoadBalancerConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerDataSourceID("data.sakuracloud_load_balancer.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceLoadBalancerConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerDataSourceNotExists("data.sakuracloud_load_balancer.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceLoadBalancerConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerDataSourceNotExists("data.sakuracloud_load_balancer.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudLoadBalancerDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find LoadBalancer data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("LoadBalancer data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudLoadBalancerDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found LoadBalancer data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudDataSourceLoadBalancerBase(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_switch "sw"{
  name = "%s"
}
resource "sakuracloud_load_balancer" "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"
  vrid = 1
  ipaddress1 = "192.168.11.101"
  nw_mask_len = 24
  default_route = "192.168.11.1"

  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceLoadBalancerConfig(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_switch "sw"{
  name = "%s"
}
resource "sakuracloud_load_balancer" "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"
  vrid = 1
  ipaddress1 = "192.168.11.101"
  nw_mask_len = 24
  default_route = "192.168.11.1"

  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_load_balancer" "foobar" {
  filters {
	names = ["%s"]
  }
}`, name, name, name)
}

func testAccCheckSakuraCloudDataSourceLoadBalancerConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_switch "sw"{
  name = "%s"
}
resource "sakuracloud_load_balancer" "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"
  vrid = 1
  ipaddress1 = "192.168.11.101"
  nw_mask_len = 24
  default_route = "192.168.11.1"

  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_load_balancer" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceLoadBalancerConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_switch "sw"{
  name = "%s"
}
resource "sakuracloud_load_balancer" "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"
  vrid = 1
  ipaddress1 = "192.168.11.101"
  nw_mask_len = 24
  default_route = "192.168.11.1"

  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_load_balancer" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceLoadBalancerConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_switch "sw"{
  name = "%s"
}
resource "sakuracloud_load_balancer" "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"
  vrid = 1
  ipaddress1 = "192.168.11.101"
  nw_mask_len = 24
  default_route = "192.168.11.1"

  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_load_balancer" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, name, name)
}
