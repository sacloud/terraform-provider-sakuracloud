package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceNote_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudNoteDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceNoteBase,
				Check:  testAccCheckSakuraCloudNoteDataSourceID("sakuracloud_note.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceID("data.sakuracloud_note.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "content", "content_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "class", "shell"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceID("data.sakuracloud_note.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNote_NameSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceID("data.sakuracloud_note.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNote_TagSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceID("data.sakuracloud_note.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceNotExists("data.sakuracloud_note.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceNotExists("data.sakuracloud_note.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNote_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceNotExists("data.sakuracloud_note.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNote_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceNotExists("data.sakuracloud_note.foobar"),
				),
				Destroy: true,
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
			return errors.New("Note data source ID not set")
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
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_note" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Note.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Note still exists")
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

var testAccCheckSakuraCloudDataSourceNote_NameSelector_Exists = `
resource "sakuracloud_note" "foobar" {
    name = "name_test"
    description = "description_test"
    content = "content_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_note" "foobar" {
    name_selectors = ["name", "test"]
}
`
var testAccCheckSakuraCloudDataSourceNote_NameSelector_NotExists = `
data "sakuracloud_note" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}
`

var testAccCheckSakuraCloudDataSourceNote_TagSelector_Exists = `
resource "sakuracloud_note" "foobar" {
    name = "name_test"
    description = "description_test"
    content = "content_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_note" "foobar" {
	tag_selectors = ["tag1","tag2","tag3"]
}`

var testAccCheckSakuraCloudDataSourceNote_TagSelector_NotExists = `
data "sakuracloud_note" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
}`
