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

func TestAccSakuraCloudDataSourceServerVNCInfo_basic(t *testing.T) {
	resourceName := "data.sakuracloud_server_vnc_info.foobar"
	rand := randomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceServerVNCInfo_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "server_id"),
					resource.TestCheckResourceAttrSet(resourceName, "host"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
					resource.TestCheckResourceAttrSet(resourceName, "password"),
					resource.TestCheckResourceAttrSet(resourceName, "zone"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceServerVNCInfo_basic = `
resource "sakuracloud_server" "foobar" {
  name        = "{{ .arg0 }}"
  network_interface {
    upstream = "shared"
  }
  force_shutdown = true
}

data "sakuracloud_server_vnc_info" "foobar" {
  server_id = sakuracloud_server.foobar.id
}`
