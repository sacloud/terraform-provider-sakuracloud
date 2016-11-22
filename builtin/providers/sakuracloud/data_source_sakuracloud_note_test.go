package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudNoteDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudNoteDataSourceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceNoteBase,
				Check:  testAccCheckSakuraCloudNoteDataSourceID("sakuracloud_note.foobar"),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceID("data.sakuracloud_note.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "content", "content_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.2", "tag3"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceNoteConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceID("data.sakuracloud_note.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceNoteConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceNotExists("data.sakuracloud_note.foobar"),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceNoteConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceNotExists("data.sakuracloud_note.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudNoteDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Note data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Note data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudNoteDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found Note data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudNoteDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_note" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Note.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("Note still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceNoteBase = `
resource "sakuracloud_note" "foobar" {
    name = "name_test"
    description = "description_test"
    content = "content_test"
    tags = ["tag1","tag2","tag3"]
}
`

var testAccCheckSakuraCloudDataSourceNoteConfig = `
resource "sakuracloud_note" "foobar" {
    name = "name_test"
    description = "description_test"
    content = "content_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_note" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`

var testAccCheckSakuraCloudDataSourceNoteConfig_With_Tag = `
resource "sakuracloud_note" "foobar" {
    name = "name_test"
    description = "description_test"
    content = "content_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_note" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`

var testAccCheckSakuraCloudDataSourceNoteConfig_With_NotExists_Tag = `
resource "sakuracloud_note" "foobar" {
    name = "name_test"
    description = "description_test"
    content = "content_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_note" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceNoteConfig_NotExists = `
resource "sakuracloud_note" "foobar" {
    name = "name_test"
    description = "description_test"
    content = "content_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_note" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`
