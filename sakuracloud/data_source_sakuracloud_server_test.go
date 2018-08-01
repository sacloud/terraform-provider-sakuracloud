package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceServer_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudServerDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceServerBase(name),
				Check:  testAccCheckSakuraCloudServerDataSourceID("sakuracloud_server.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServerConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceID("data.sakuracloud_server.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "interface_driver", "virtio"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "tags.2", "tag3"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "core", "1"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "memory", "1"),
					//resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "disks.#", "1"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "nic", "shared"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "additional_nics.#", "0"),
					resource.TestCheckResourceAttr("data.sakuracloud_server.foobar", "macaddresses.#", "1"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServerConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceID("data.sakuracloud_server.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServer_NameSelector_Exists(name, randString1, randString2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceID("data.sakuracloud_server.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServer_TagSelector_Exists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceID("data.sakuracloud_server.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServerConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceNotExists("data.sakuracloud_server.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServerConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceNotExists("data.sakuracloud_server.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServer_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceNotExists("data.sakuracloud_server.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceServer_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerDataSourceNotExists("data.sakuracloud_server.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudServerDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Server data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Server data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudServerDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found Server data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudServerDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_server" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Server.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Server still exists")
		}
	}

	return nil
}

func testAccCheckSakuraCloudDataSourceServerBase(name string) string {
	return fmt.Sprintf(`
data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu"
}
resource "sakuracloud_disk" "foobar" {
  name = "%s"
  source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}
resource "sakuracloud_server" "foobar" {
  name = "%s"
  disks = ["${sakuracloud_disk.foobar.id}"]
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceServerConfig(name string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    filter = {
	name = "Name"
	values = ["%s"]
    }
}`, testAccCheckSakuraCloudDataSourceServerBase(name), name)
}

func testAccCheckSakuraCloudDataSourceServerConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`, testAccCheckSakuraCloudDataSourceServerBase(name))
}

func testAccCheckSakuraCloudDataSourceServerConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`, testAccCheckSakuraCloudDataSourceServerBase(name))
}

func testAccCheckSakuraCloudDataSourceServerConfig_NotExists(name string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`, testAccCheckSakuraCloudDataSourceServerBase(name))
}

func testAccCheckSakuraCloudDataSourceServer_NameSelector_Exists(name, p1, p2 string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
    name_selectors = ["%s", "%s"]
}`, testAccCheckSakuraCloudDataSourceServerBase(name), p1, p2)
}

var testAccCheckSakuraCloudDataSourceServer_NameSelector_NotExists = `
data "sakuracloud_server" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}`

func testAccCheckSakuraCloudDataSourceServer_TagSelector_Exists(name string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_server" "foobar" {
	tag_selectors = ["tag1","tag2","tag3"]
}`, testAccCheckSakuraCloudDataSourceServerBase(name))
}

var testAccCheckSakuraCloudDataSourceServer_TagSelector_NotExists = `
data "sakuracloud_server" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
}`
