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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/sacloud"
)

func TestAccResourceSakuraCloudGSLBServer(t *testing.T) {
	var gslb sacloud.GSLB
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudGSLBServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudGSLBServerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBExists("sakuracloud_gslb.foobar", &gslb),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.0", "ipaddress", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.1", "ipaddress", "8.8.4.4"),
				),
			},
			{
				Config: testAccCheckSakuraCloudGSLBServerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBExists("sakuracloud_gslb.foobar", &gslb),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.0", "ipaddress", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.1", "ipaddress", "8.8.4.4"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.2", "ipaddress", "208.67.222.123"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.3", "ipaddress", "208.67.220.123"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudGSLBServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_gslb" {
			continue
		}

		_, err := client.GSLB.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("GSLB still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudGSLBServerConfig_basic = `
variable "gslb_ip_list" {
    default = ["8.8.8.8", "8.8.4.4"]
}
resource "sakuracloud_gslb" "foobar" {
    name = "terraform.io"
    health_check {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    description = "GSLB from TerraForm for SAKURA CLOUD"
    tags = ["hoge1", "hoge2"]
}
resource "sakuracloud_gslb_server" "foobar" {
    count = 2
    gslb_id = "${sakuracloud_gslb.foobar.id}"
    ipaddress = "${var.gslb_ip_list[count.index]}"
}`

var testAccCheckSakuraCloudGSLBServerConfig_update = `
variable "gslb_ip_list" {
    default = ["8.8.8.8","8.8.4.4", "208.67.222.123", "208.67.220.123"]
}
resource "sakuracloud_gslb" "foobar" {
    name = "terraform.io"
    health_check {
        protocol = "https"
        delay_loop = 20
        host_header = "update.terraform.io"
        path = "/"
        status = "200"
    }
    description = "GSLB from TerraForm for SAKURA CLOUD"
    tags = ["hoge1", "hoge2"]
}
resource "sakuracloud_gslb_server" "foobar" {
    count = 4
    gslb_id = "${sakuracloud_gslb.foobar.id}"
    ipaddress = "${var.gslb_ip_list[count.index]}"
}`
