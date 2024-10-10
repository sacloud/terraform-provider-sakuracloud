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

func TestAccSakuraCloudDataSourceAutoScale_basic(t *testing.T) {
	resourceName := "data.sakuracloud_auto_scale.foobar"
	rand := randomName()
	if !isFakeModeEnabled() {
		skipIfEnvIsNotSet(t, "SAKURACLOUD_API_KEY_ID")
	}
	apiKeyId := os.Getenv("SAKURACLOUD_API_KEY_ID")
	if apiKeyId == "" {
		apiKeyId = "111111111111" // dummy
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceAutoScale_basic, rand, apiKeyId),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),

					resource.TestCheckResourceAttr(resourceName, "zones.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "zones.0", "is1b"),
					resource.TestCheckResourceAttr(resourceName, "cpu_threshold_scaling.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cpu_threshold_scaling.0.server_prefix", rand),
					resource.TestCheckResourceAttr(resourceName, "cpu_threshold_scaling.0.up", "80"),
					resource.TestCheckResourceAttr(resourceName, "cpu_threshold_scaling.0.down", "20"),
					resource.TestCheckResourceAttr(resourceName, "config", buildConfigWithArgs(testAccSakuraCloudAutoScale_encodedConfig, rand)),
					resource.TestCheckResourceAttrSet(resourceName, "api_key_id"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceAutoScale_basic = `
resource "sakuracloud_server" "foobar" {
  name = "{{ .arg0 }}"
  force_shutdown = true
  zone = "is1b"
}

resource "sakuracloud_auto_scale" "foobar" {
  name           = "{{ .arg0 }}"
  description    = "description"
  tags           = ["tag1", "tag2"]

  zones  = ["is1b"]
  config = yamlencode({
    resources: [{
      type: "Server",
      selector: {
        names: [sakuracloud_server.foobar.name],
        zones: ["is1b"],
      },
      shutdown_force: true,
    }],
  })
  api_key_id = "{{ .arg1 }}"

  cpu_threshold_scaling {
    server_prefix = sakuracloud_server.foobar.name

    up   = 80
    down = 20
  }
}


data "sakuracloud_auto_scale" "foobar" {
  filter {
    names = [sakuracloud_auto_scale.foobar.name]
  }
}`
