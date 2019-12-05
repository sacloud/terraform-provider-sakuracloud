package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSakuraCloudDataSourceNFS_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudNFSDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceNFSBase(name),
				Check:  testAccCheckSakuraCloudNFSDataSourceID("sakuracloud_nfs.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNFSConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNFSDataSourceID("data.sakuracloud_nfs.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_nfs.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_nfs.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_nfs.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_nfs.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_nfs.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_nfs.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNFSConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNFSDataSourceID("data.sakuracloud_nfs.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNFSConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNFSDataSourceNotExists("data.sakuracloud_nfs.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNFSConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNFSDataSourceNotExists("data.sakuracloud_nfs.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudNFSDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find NFS data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("NFS data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudNFSDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found NFS data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudDataSourceNFSBase(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_switch "sw"{
  name = "%s"
}
resource "sakuracloud_nfs" "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"
  ipaddress = "192.168.11.101"
  nw_mask_len = 24
  default_route = "192.168.11.1"

  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceNFSConfig(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_switch "sw"{
  name = "%s"
}
resource "sakuracloud_nfs" "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"
  ipaddress = "192.168.11.101"
  nw_mask_len = 24
  default_route = "192.168.11.1"

  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_nfs" "foobar" {
  filters {
	names = ["%s"]
  }
}`, name, name, name)
}

func testAccCheckSakuraCloudDataSourceNFSConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_switch "sw"{
  name = "%s"
}
resource "sakuracloud_nfs" "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"
  ipaddress = "192.168.11.101"
  nw_mask_len = 24
  default_route = "192.168.11.1"

  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_nfs" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceNFSConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_switch "sw"{
  name = "%s"
}
resource "sakuracloud_nfs" "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"
  ipaddress = "192.168.11.101"
  nw_mask_len = 24
  default_route = "192.168.11.1"

  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_nfs" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceNFSConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_switch "sw"{
  name = "%s"
}
resource "sakuracloud_nfs" "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"
  ipaddress = "192.168.11.101"
  nw_mask_len = 24
  default_route = "192.168.11.1"

  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_nfs" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, name, name)
}
