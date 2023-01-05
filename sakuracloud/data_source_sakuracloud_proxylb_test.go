// Copyright 2016-2023 terraform-provider-sakuracloud authors
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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSakuraCloudDataSourceProxyLB_basic(t *testing.T) {
	skipIfEnvIsNotSet(t, envProxyLBRealServerIP0, envProxyLBRealServerIP1)

	resourceName := "data.sakuracloud_proxylb.foobar"
	rand := randomName()
	ip0 := os.Getenv(envProxyLBRealServerIP0)
	ip1 := os.Getenv(envProxyLBRealServerIP1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceProxyLB_basic, rand, ip0, ip1),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "plan", "100"),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "tag3"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.delay_loop", "20"),
					resource.TestCheckResourceAttr(resourceName, "bind_port.0.proxy_mode", "http"),
					resource.TestCheckResourceAttr(resourceName, "bind_port.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "server.0.ip_address", ip0),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "server.0.group", "group1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.1.ip_address", ip1),
					resource.TestCheckResourceAttr(resourceName, "server.1.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "server.1.group", "group2"),
					resource.TestCheckResourceAttr(resourceName, "server.1.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.host", ""),
					resource.TestCheckResourceAttr(resourceName, "rule.0.path", "/path1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.group", "group1"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.host", ""),
					resource.TestCheckResourceAttr(resourceName, "rule.1.path", "/path2"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.group", "group2"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceProxyLB_basic = `
resource "sakuracloud_proxylb" "foobar" {
  name = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2", "tag3"]

  health_check {
    protocol   = "tcp"
    delay_loop = 20
  }

  bind_port {
    proxy_mode = "http"
    port       = 80
  }

  server {
    ip_address = "{{ .arg1 }}"
    port       = 80
    group      = "group1"
  }
  server {
    ip_address = "{{ .arg2 }}"
    port       = 80
    group      = "group2"
  }

  rule {
    path  = "/path1"
    group = "group1"
  }
  rule {
    path  = "/path2"
    group = "group2"
  }
}

data "sakuracloud_proxylb" "foobar" {
  filter {
	names = [sakuracloud_proxylb.foobar.name]
  }
}`
