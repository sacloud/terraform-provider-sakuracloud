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
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func TestAccResourceSakuraCloudBridge(t *testing.T) {
	var bridge sacloud.Bridge
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudBridgeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudBridgeConfig_withSwitch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeExists("sakuracloud_bridge.foobar", &bridge),
					resource.TestCheckResourceAttr(
						"sakuracloud_bridge.foobar", "name", "mybridge"),
					resource.TestCheckResourceAttr(
						"sakuracloud_switch.foobar", "name", "myswitch"),
				),
			},
			{
				Config: testAccCheckSakuraCloudBridgeConfig_withSwitch,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_bridge.foobar", "switch_ids.#", "1"),
				),
			},
			{
				Config: testAccCheckSakuraCloudBridgeConfig_withSwitchDisconnect,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeExists("sakuracloud_bridge.foobar", &bridge),
					resource.TestCheckResourceAttr(
						"sakuracloud_bridge.foobar", "name", "mybridge_upd"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudBridgeExists(n string, bridge *sacloud.Bridge) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no bridge ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		bridgeOp := sacloud.NewBridgeOp(client)
		zone := rs.Primary.Attributes["zone"]
		foundBridge, err := bridgeOp.Read(context.Background(), zone, types.StringID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundBridge.ID.String() != rs.Primary.ID {
			return errors.New("bridge not found")
		}

		*bridge = *foundBridge

		return nil
	}
}

func testAccCheckSakuraCloudBridgeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_bridge" {
			continue
		}

		bridgeOp := sacloud.NewBridgeOp(client)
		zone := rs.Primary.Attributes["zone"]
		_, err := bridgeOp.Read(context.Background(), zone, types.StringID(rs.Primary.ID))

		if err == nil {
			return errors.New("bridge still exists")
		}
	}

	return nil
}

func TestAccImportSakuraCloudBridge(t *testing.T) {
	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %#v", s)
		}
		expects := map[string]string{
			"name":        "mybridge",
			"description": "Bridge from TerraForm for SAKURA CLOUD",
			"zone":        os.Getenv("SAKURACLOUD_ZONE"),
		}

		if err := compareStateMulti(s[0], expects); err != nil {
			return err
		}
		return stateNotEmptyMulti(s[0], "switch_ids.0")
	}

	resourceName := "sakuracloud_bridge.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudBridgeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudBridgeConfig_withSwitch,
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

var testAccCheckSakuraCloudBridgeConfig_withSwitch = `
resource "sakuracloud_switch" "foobar" {
  name = "myswitch"
  description = "Switch from TerraForm for SAKURA CLOUD"
  bridge_id = "${sakuracloud_bridge.foobar.id}"
}
resource "sakuracloud_bridge" "foobar" {
  name = "mybridge"
  description = "Bridge from TerraForm for SAKURA CLOUD"
}`

var testAccCheckSakuraCloudBridgeConfig_withSwitchDisconnect = `
resource "sakuracloud_bridge" "foobar" {
  name = "mybridge_upd"
  description = "Bridge from TerraForm for SAKURA CLOUD"
}`
