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
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/sacloud/iaas-api-go"
)

func TestAccSakuraCloudAutoScale_basic(t *testing.T) {
	resourceName := "sakuracloud_auto_scale.foobar"
	rand := randomName()
	if !isFakeModeEnabled() {
		skipIfEnvIsNotSet(t, "SAKURACLOUD_API_KEY_ID")
	}
	apiKeyId := os.Getenv("SAKURACLOUD_API_KEY_ID")
	if apiKeyId == "" {
		apiKeyId = "111111111111" // dummy
	}

	var autoScale iaas.AutoScale
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudAutoScaleDestroy,
			testCheckSakuraCloudServerDestroy,
			testCheckSakuraCloudIconDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudAutoScale_basic, rand, apiKeyId),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudAutoScaleExists(resourceName, &autoScale),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttrPair(
						resourceName, "icon_id",
						"sakuracloud_icon.foobar", "id",
					),
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
			{
				Config: buildConfigWithArgs(testAccSakuraCloudAutoScale_update, rand, apiKeyId),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudAutoScaleExists(resourceName, &autoScale),
					resource.TestCheckResourceAttr(resourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "description-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2-upd"),
					resource.TestCheckResourceAttr(resourceName, "icon_id", ""),

					resource.TestCheckResourceAttr(resourceName, "zones.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "zones.0", "is1b"),

					resource.TestCheckResourceAttr(resourceName, "cpu_threshold_scaling.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cpu_threshold_scaling.0.server_prefix", rand),
					resource.TestCheckResourceAttr(resourceName, "cpu_threshold_scaling.0.up", "81"),
					resource.TestCheckResourceAttr(resourceName, "cpu_threshold_scaling.0.down", "21"),

					resource.TestCheckResourceAttr(resourceName, "config", buildConfigWithArgs(testAccSakuraCloudAutoScale_encodedConfig_update, rand)),
					resource.TestCheckResourceAttrSet(resourceName, "api_key_id"),
				),
			},
		},
	})
}

func TestAccSakuraCloudAutoScale_withRouter(t *testing.T) {
	resourceName := "sakuracloud_auto_scale.foobar"
	rand := randomName()
	if !isFakeModeEnabled() {
		skipIfEnvIsNotSet(t, "SAKURACLOUD_API_KEY_ID")
	}
	apiKeyId := os.Getenv("SAKURACLOUD_API_KEY_ID")
	if apiKeyId == "" {
		apiKeyId = "111111111111" // dummy
	}

	var autoScale iaas.AutoScale
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudAutoScaleDestroy,
			testCheckSakuraCloudInternetDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudAutoScale_withRouter, rand, apiKeyId),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudAutoScaleExists(resourceName, &autoScale),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "cpu_threshold_scaling.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "router_threshold_scaling.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "router_threshold_scaling.0.router_prefix", rand),
					resource.TestCheckResourceAttr(resourceName, "router_threshold_scaling.0.direction", "in"),
					resource.TestCheckResourceAttr(resourceName, "router_threshold_scaling.0.mbps", "20"),

					resource.TestCheckResourceAttr(resourceName, "config", buildConfigWithArgs(testAccSakuraCloudAutoScale_encodedConfigWithRouter, rand)),
					resource.TestCheckResourceAttrSet(resourceName, "api_key_id"),
				),
			},
		},
	})
}

func testCheckSakuraCloudAutoScaleExists(n string, auto_scale *iaas.AutoScale) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no AutoScale ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		autoScaleOp := iaas.NewAutoScaleOp(client)

		foundAutoScale, err := autoScaleOp.Read(context.Background(), sakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundAutoScale.ID.String() != rs.Primary.ID {
			return errors.New("resource not found")
		}

		*auto_scale = *foundAutoScale

		return nil
	}
}

func testCheckSakuraCloudAutoScaleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_auto_scale" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		autoScaleOp := iaas.NewAutoScaleOp(client)
		_, err := autoScaleOp.Read(context.Background(), sakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("AutoScale still exists")
		}
	}

	return nil
}

func TestAccImportSakuraCloudAutoScale_basic(t *testing.T) {
	if !isFakeModeEnabled() {
		skipIfEnvIsNotSet(t, "SAKURACLOUD_API_KEY_ID")
	}
	apiKeyId := os.Getenv("SAKURACLOUD_API_KEY_ID")
	rand := randomName()

	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %#v", s)
		}
		expects := map[string]string{
			"name":                                  rand,
			"description":                           "description",
			"tags.0":                                "tag1",
			"tags.1":                                "tag2",
			"zones.0":                               "is1b",
			"api_key_id":                            apiKeyId,
			"cpu_threshold_scaling.0.server_prefix": rand,
			"cpu_threshold_scaling.0.up":            "80",
			"cpu_threshold_scaling.0.down":          "20",
			"config":                                buildConfigWithArgs(testAccSakuraCloudAutoScale_encodedConfig, rand),
		}

		if err := compareStateMulti(s[0], expects); err != nil {
			return err
		}
		return stateNotEmptyMulti(s[0], "icon_id")
	}

	resourceName := "sakuracloud_auto_scale.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudAutoScaleDestroy,
			testCheckSakuraCloudServerDestroy,
			testCheckSakuraCloudIconDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudAutoScale_basic, rand, os.Getenv("SAKURACLOUD_API_KEY_ID")),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateCheck:  checkFn,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccSakuraCloudAutoScale_basic = `
resource "sakuracloud_server" "foobar" {
  name = "{{ .arg0 }}"
  force_shutdown = true
  zone = "is1b"
}

resource "sakuracloud_auto_scale" "foobar" {
  name           = "{{ .arg0 }}"
  description    = "description"
  tags           = ["tag1", "tag2"]
  icon_id        = sakuracloud_icon.foobar.id

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
    server_prefix = "{{ .arg0 }}"

    up   = 80
    down = 20
  }
}

resource "sakuracloud_icon" "foobar" {
  name          = "{{ .arg0 }}"
  base64content = "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"
}
`

var testAccSakuraCloudAutoScale_encodedConfig = `"resources":
- "selector":
    "names":
    - "{{ .arg0 }}"
    "zones":
    - "is1b"
  "shutdown_force": true
  "type": "Server"
`

var testAccSakuraCloudAutoScale_update = `
resource "sakuracloud_server" "foobar" {
  name = "{{ .arg0 }}"
  force_shutdown = true
  zone = "is1b"
}

resource "sakuracloud_auto_scale" "foobar" {
  name           = "{{ .arg0 }}-upd"
  description    = "description-upd"
  tags           = ["tag1-upd", "tag2-upd"]

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
    autoscaler: {
      cooldown: 300,
    },
  })

  api_key_id = "{{ .arg1 }}"

  cpu_threshold_scaling {
    server_prefix = "{{ .arg0 }}"

    up   = 81
    down = 21
  }
}
`

var testAccSakuraCloudAutoScale_encodedConfig_update = `"autoscaler":
  "cooldown": 300
"resources":
- "selector":
    "names":
    - "{{ .arg0 }}"
    "zones":
    - "is1b"
  "shutdown_force": true
  "type": "Server"
`

var testAccSakuraCloudAutoScale_withRouter = `
resource "sakuracloud_internet" "foobar" {
  name = "{{ .arg0 }}"
  zone = "is1b"
}

resource "sakuracloud_auto_scale" "foobar" {
  name           = "{{ .arg0 }}"

  zones  = ["is1b"]
  config = yamlencode({
    resources: [{
      type: "Router",
      selector: {
        names: [sakuracloud_internet.foobar.name],
        zones: ["is1b"],
      },
    }],
  })
  api_key_id = "{{ .arg1 }}"

  trigger_type = "router"
  router_threshold_scaling {
    router_prefix = "{{ .arg0 }}"
    direction     = "in"
    mbps          = 20
  }
}
`

var testAccSakuraCloudAutoScale_encodedConfigWithRouter = `"resources":
- "selector":
    "names":
    - "{{ .arg0 }}"
    "zones":
    - "is1b"
  "type": "Router"
`
