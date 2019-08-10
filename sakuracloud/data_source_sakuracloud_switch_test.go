package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceSwitch_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudSwitchDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchBase(name),
				Check:  testAccCheckSakuraCloudSwitchDataSourceID("sakuracloud_switch.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchDataSourceID("data.sakuracloud_switch.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchDataSourceID("data.sakuracloud_switch.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchDataSourceNotExists("data.sakuracloud_switch.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSwitchDataSourceNotExists("data.sakuracloud_switch.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudSwitchDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Switch data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Switch data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudSwitchDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found Switch data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudSwitchDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_switch" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Switch.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Switch still exists")
		}
	}

	return nil
}

func testAccCheckSakuraCloudDataSourceSwitchBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
`, name)
}

func testAccCheckSakuraCloudDataSourceSwitchConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
  filters {
	names = ["%s"]
  }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceSwitchConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourceSwitchConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourceSwitchConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, name)
}
