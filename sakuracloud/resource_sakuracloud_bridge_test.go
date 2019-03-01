package sakuracloud

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/sacloud"
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
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Bridge ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)

		foundBridge, err := client.Bridge.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundBridge.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("Bridge not found")
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

		_, err := client.Bridge.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Bridge still exists")
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
    bridge_id = sakuracloud_bridge.foobar.id
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
