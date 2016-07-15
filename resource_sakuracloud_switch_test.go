package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccResourceSakuraCloudSwitch_Basic(t *testing.T) {
	var sw sacloud.Switch
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudSwitchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudSwitchConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchExists("sakuracloud_switch.foobar", &sw),
					resource.TestCheckResourceAttr(
						"sakuracloud_switch.foobar", "name", "myswitch"),
					resource.TestCheckResourceAttr(
						"sakuracloud_switch.foobar", "server_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceSakuraCloudSwitch_Update(t *testing.T) {
	var sw sacloud.Switch
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudSwitchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudSwitchConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchExists("sakuracloud_switch.foobar", &sw),
					resource.TestCheckResourceAttr(
						"sakuracloud_switch.foobar", "name", "myswitch"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudSwitchConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchExists("sakuracloud_switch.foobar", &sw),
					resource.TestCheckResourceAttr(
						"sakuracloud_switch.foobar", "name", "myswitch_upd"),
				),
			},
		},
	})
}

func TestAccResourceSakuraCloudSwitch_WithServer(t *testing.T) {
	var sw sacloud.Switch
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudSwitchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudSwitchConfig_with_server,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchExists("sakuracloud_switch.foobar", &sw),
					resource.TestCheckResourceAttr(
						"sakuracloud_switch.foobar", "name", "myswitch"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudSwitchConfig_with_server,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_switch.foobar", "server_ids.#", "1"),
				),
			},
		},
	})
}

//func TestAccResourceSakuraCloudSwitch_Import(t *testing.T) {
//	resourceName := "sakuracloud_switch.foobar"
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testAccCheckSakuraCloudSwitchDestroy,
//		Steps: []resource.TestStep{
//			resource.TestStep{
//				Config: testAccCheckSakuraCloudSwitchConfig_basic,
//			},
//			resource.TestStep{
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//		},
//	})
//}

func testAccCheckSakuraCloudSwitchExists(n string, sw *sacloud.Switch) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Switch ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)

		foundSwitch, err := client.Switch.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundSwitch.ID != rs.Primary.ID {
			return fmt.Errorf("Switch not found")
		}

		*sw = *foundSwitch

		return nil
	}
}

func testAccCheckSakuraCloudSwitchDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_switch" {
			continue
		}

		_, err := client.Switch.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Switch still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudSwitchConfig_basic = `
resource "sakuracloud_switch" "foobar" {
    name = "myswitch"
    description = "Switch from TerraForm for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
}`

var testAccCheckSakuraCloudSwitchConfig_update = `
resource "sakuracloud_switch" "foobar" {
    name = "myswitch_upd"
    description = "Switch from TerraForm for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
}`

var testAccCheckSakuraCloudSwitchConfig_with_server = `
resource "sakuracloud_server" "foobar" {
    name = "myserver"
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
    additional_interfaces = ["${sakuracloud_switch.foobar.id}"]
}
resource "sakuracloud_switch" "foobar" {
    name = "myswitch"
    description = "Switch from TerraForm for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
}
`
