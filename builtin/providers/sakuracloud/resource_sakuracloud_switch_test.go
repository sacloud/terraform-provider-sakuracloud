package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"testing"
)

func TestAccResourceSakuraCloudSwitch(t *testing.T) {
	var sw sacloud.Switch
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudSwitchConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchExists("sakuracloud_switch.foobar", &sw),
					resource.TestCheckResourceAttr(
						"sakuracloud_switch.foobar", "name", "myswitch"),
				),
			},
			{
				Config: testAccCheckSakuraCloudSwitchConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchExists("sakuracloud_switch.foobar", &sw),
					resource.TestCheckResourceAttr(
						"sakuracloud_switch.foobar", "name", "myswitch_upd"),
				),
			},
			{
				Config: testAccCheckSakuraCloudSwitchConfig_update,
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
//			{
//				Config: testAccCheckSakuraCloudSwitchConfig_basic,
//			},
//			{
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
			return errors.New("No Switch ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)

		foundSwitch, err := client.Switch.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundSwitch.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("Switch not found")
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

		_, err := client.Switch.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Switch still exists")
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
resource "sakuracloud_server" "foobar" {
    name = "myserver"
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
    additional_nics = ["${sakuracloud_switch.foobar.id}"]
}
resource "sakuracloud_switch" "foobar" {
    name = "myswitch_upd"
    description = "Switch from TerraForm for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
}
`
