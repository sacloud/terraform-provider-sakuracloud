package sakuracloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"

	"github.com/sacloud/libsacloud/sacloud"
)

func TestAccMarkerTags(t *testing.T) {
	var sw sacloud.Switch
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkerTags_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchHasTags("sakuracloud_switch.foobar", &sw),
					resource.TestCheckResourceAttr("sakuracloud_switch.foobar", "name", "myswitch"),
					resource.TestCheckResourceAttr("sakuracloud_switch.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("sakuracloud_switch.foobar", "tags.1", "tag2"),
				),
			},
			{
				Config: tesetAccMarkerTags_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchHasTags("sakuracloud_switch.foobar", &sw),
					resource.TestCheckResourceAttr("sakuracloud_switch.foobar", "name", "myswitch_upd"),
					resource.TestCheckResourceAttr("sakuracloud_switch.foobar", "tags.#", "0"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudSwitchHasTags(n string, sw *sacloud.Switch) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Switch ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)

		foundSwitch, err := client.Switch.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundSwitch.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("Switch not found")
		}

		if !foundSwitch.HasTag("@terraform") {
			return errors.New("Switch should have marker tags")
		}

		*sw = *foundSwitch

		return nil
	}
}

var testAccMarkerTags_basic = `
provider sakuracloud {
    use_marker_tags = true
}

resource "sakuracloud_switch" "foobar" {
    name = "myswitch"
    tags = ["tag1" , "tag2"]
}
`

var tesetAccMarkerTags_update = `
provider sakuracloud {
    use_marker_tags = true
}
resource "sakuracloud_switch" "foobar" {
    name = "myswitch_upd"
    tags = []
}
`
