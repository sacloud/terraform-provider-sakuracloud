package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceDisk_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudDiskDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceDiskConfigBase(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_disk.disk01", "name", name),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDiskConfig(name, randString1, randString2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskDataSourceID("data.sakuracloud_disk.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "plan", "ssd"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "connector", "virtio"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "size", "20"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "description", "source_disk_description"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDiskConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskDataSourceID("data.sakuracloud_disk.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDisk_NameSelector_Exists(name, randString1, randString2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskDataSourceID("data.sakuracloud_disk.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDisk_TagSelector_Exists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskDataSourceID("data.sakuracloud_disk.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDiskConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskDataSourceNotExists("data.sakuracloud_disk.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDiskConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskDataSourceNotExists("data.sakuracloud_disk.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDisk_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskDataSourceNotExists("data.sakuracloud_disk.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDisk_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskDataSourceNotExists("data.sakuracloud_disk.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudDiskDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Disk data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Disk data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudDiskDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found Disk data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudDiskDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_disk" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Disk.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Disk still exists")
		}
	}

	return nil
}

func testAccCheckSakuraCloudDataSourceDiskConfigBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_disk" "disk01"{
    name = "%s"
    tags = ["tag1","tag2","tag3"]
    description = "source_disk_description"
}`, name)
}

func testAccCheckSakuraCloudDataSourceDiskConfig(name, p1, p2 string) string {
	return fmt.Sprintf(`
resource "sakuracloud_disk" "disk01"{
    name = "%s"
    tags = ["tag1","tag2","tag3"]
    description = "source_disk_description"
}

data "sakuracloud_disk" "foobar" {
    filter {
	name = "Name"
	values = ["%s", "%s"]
    }
}`, name, p1, p2)
}

func testAccCheckSakuraCloudDataSourceDiskConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_disk" "disk01"{
    name = "%s"
    tags = ["tag1","tag2","tag3"]
    description = "source_disk_description"
}

data "sakuracloud_disk" "foobar" {
    filter {
	name = "Tags"
	values = ["tag2","tag3"]
    }
}`, name)
}

func testAccCheckSakuraCloudDataSourceDiskConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_disk" "disk01"{
    name = "%s"
    tags = ["tag1","tag2","tag3"]
    description = "source_disk_description"
}

data "sakuracloud_disk" "foobar" {
    filter {
	name = "Tags"
	values = ["tag2","tag4"]
    }
}`, name)
}

func testAccCheckSakuraCloudDataSourceDiskConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_disk" "disk01"{
    name = "%s"
    tags = ["tag1","tag2","tag3"]
    description = "source_disk_description"
}

data "sakuracloud_disk" "foobar" {
    filter {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`, name)
}

func testAccCheckSakuraCloudDataSourceDisk_NameSelector_Exists(name, p1, p2 string) string {
	return fmt.Sprintf(`
resource "sakuracloud_disk" "disk01"{
    name = "%s"
    tags = ["tag1","tag2","tag3"]
    description = "source_disk_description"
}
data "sakuracloud_disk" "foobar" {
    name_selectors = ["%s", "%s"]
}
`, name, p1, p2)
}

var testAccCheckSakuraCloudDataSourceDisk_NameSelector_NotExists = `
data "sakuracloud_disk" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}
`

func testAccCheckSakuraCloudDataSourceDisk_TagSelector_Exists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_disk" "disk01"{
    name = "%s"
    tags = ["tag1","tag2","tag3"]
    description = "source_disk_description"
}
data "sakuracloud_disk" "foobar" {
	tag_selectors = ["tag1","tag2","tag3"]
}`, name)
}

var testAccCheckSakuraCloudDataSourceDisk_TagSelector_NotExists = `
data "sakuracloud_disk" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
}`
