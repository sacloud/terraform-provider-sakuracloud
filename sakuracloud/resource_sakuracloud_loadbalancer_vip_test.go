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
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func TestAccSakuraCloudLoadBalancerVIP(t *testing.T) {
	var loadBalancer sacloud.LoadBalancer
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudLoadBalancerVIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudLoadBalancerVIPConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerExists("sakuracloud_load_balancer.foobar", &loadBalancer),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "vip", "192.168.11.201"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "port", "80"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "delay_loop", "100"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "sorry_server", "192.168.11.11"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "description", "description"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip2", "vip", "192.168.11.202"),
				),
			},
			{
				Config: testAccCheckSakuraCloudLoadBalancerVIPConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudLoadBalancerExists("sakuracloud_load_balancer.foobar", &loadBalancer),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "vip", "192.168.11.201"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "port", "80"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "delay_loop", "50"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip1", "sorry_server", "192.168.11.22"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip2", "vip", "192.168.11.202"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip3", "vip", "192.168.11.203"),
					resource.TestCheckResourceAttr(
						"sakuracloud_load_balancer_vip.vip4", "vip", "192.168.11.204"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudLoadBalancerVIPDestroy(s *terraform.State) error {
	// TODO IDをパースしてLBのIDを取得すべき

	client := testAccProvider.Meta().(*APIClient)
	lbOp := sacloud.NewLoadBalancerOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_load_balancer" {
			continue
		}

		zone := rs.Primary.Attributes["zone"]
		_, err := lbOp.Read(context.Background(), zone, types.StringID(rs.Primary.ID))
		if err == nil {
			return errors.New("LoadBalancer still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudLoadBalancerVIPConfig_basic = `
resource "sakuracloud_switch" "sw" {
  name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"
  vrid = 1
  ipaddress1 = "192.168.11.101"
  nw_mask_len = 24
  default_route = "192.168.11.1"

  name = "name"
  description = "description"
  tags = ["hoge1" , "hoge2"]
}
resource "sakuracloud_load_balancer_vip" "vip1" {
  load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
  vip = "192.168.11.201"
  port = 80
  delay_loop = 100
  sorry_server = "192.168.11.11"
  description  = "description"
}
resource "sakuracloud_load_balancer_vip" "vip2" {
  load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
  vip = "192.168.11.202"
  port = 80
}
`

var testAccCheckSakuraCloudLoadBalancerVIPConfig_update = `
resource "sakuracloud_switch" "sw" {
  name = "sw"
}
resource "sakuracloud_load_balancer" "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"
  vrid = 1
  ipaddress1 = "192.168.11.101"
  nw_mask_len = 24
  default_route = "192.168.11.1"

  name = "name"
  description = "description"
  tags = ["hoge1" , "hoge2"]
}
resource "sakuracloud_load_balancer_vip" "vip1" {
  load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
  vip = "192.168.11.201"
  port = 80
  delay_loop = 50
  sorry_server = "192.168.11.22"
}
resource "sakuracloud_load_balancer_vip" "vip2" {
  load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
  vip = "192.168.11.202"
  port = 80
}
resource "sakuracloud_load_balancer_vip" "vip3" {
  load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
  vip = "192.168.11.203"
  port = 80
}
resource "sakuracloud_load_balancer_vip" "vip4" {
  load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
  vip = "192.168.11.204"
  port = 80
}
`
