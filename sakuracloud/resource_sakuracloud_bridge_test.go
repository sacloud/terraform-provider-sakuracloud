// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func TestAccSakuraCloudBridge_basic(t *testing.T) {
	resourceName := "sakuracloud_bridge.foobar"
	rand := randomName()

	var bridge sacloud.Bridge
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudBridgeDestroy,
			testCheckSakuraCloudSwitchDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudBridge_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudBridgeExists(resourceName, &bridge),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudBridge_disconnectSwitch, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudBridgeExists(resourceName, &bridge),
					resource.TestCheckResourceAttr(resourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "description-upd"),
				),
			},
		},
	})
}

func testCheckSakuraCloudBridgeExists(n string, bridge *sacloud.Bridge) resource.TestCheckFunc {
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
		foundBridge, err := bridgeOp.Read(context.Background(), zone, sakuraCloudID(rs.Primary.ID))

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

func testCheckSakuraCloudBridgeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_bridge" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		bridgeOp := sacloud.NewBridgeOp(client)
		zone := rs.Primary.Attributes["zone"]
		_, err := bridgeOp.Read(context.Background(), zone, sakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("bridge still exists")
		}
	}

	return nil
}

func TestAccImportSakuraCloudBridge_basic(t *testing.T) {
	rand := randomName()
	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %#v", s)
		}
		expects := map[string]string{
			"name":        rand,
			"description": "description",
			"zone":        os.Getenv("SAKURACLOUD_ZONE"),
		}

		return compareStateMulti(s[0], expects)
	}

	resourceName := "sakuracloud_bridge.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudBridgeDestroy,
			testCheckSakuraCloudSwitchDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudBridge_basic, rand),
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

var testAccSakuraCloudBridge_basic = `
resource "sakuracloud_switch" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  bridge_id   = sakuracloud_bridge.foobar.id
}
resource "sakuracloud_bridge" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
}`

var testAccSakuraCloudBridge_disconnectSwitch = `
resource "sakuracloud_bridge" "foobar" {
  name        = "{{ .arg0 }}-upd"
  description = "description-upd"
}`
