package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccResourceSakuraCloudNote_Basic(t *testing.T) {
	var note sacloud.Note
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudNoteDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudNoteConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteExists("sakuracloud_note.foobar", &note),
					resource.TestCheckResourceAttr(
						"sakuracloud_note.foobar", "name", "mynote"),
					resource.TestCheckResourceAttr(
						"sakuracloud_note.foobar", "content", "content"),
					resource.TestCheckResourceAttr(
						"sakuracloud_note.foobar", "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceSakuraCloudNote_Update(t *testing.T) {
	var note sacloud.Note
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudNoteDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudNoteConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteExists("sakuracloud_note.foobar", &note),
					resource.TestCheckResourceAttr(
						"sakuracloud_note.foobar", "name", "mynote"),
					resource.TestCheckResourceAttr(
						"sakuracloud_note.foobar", "content", "content"),
					resource.TestCheckResourceAttr(
						"sakuracloud_note.foobar", "tags.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudNoteConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteExists("sakuracloud_note.foobar", &note),
					resource.TestCheckResourceAttr(
						"sakuracloud_note.foobar", "name", "mynote_upd"),
					resource.TestCheckResourceAttr(
						"sakuracloud_note.foobar", "content", "content_upd"),
					resource.TestCheckResourceAttr(
						"sakuracloud_note.foobar", "tags.#", "0"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudNoteExists(n string, note *sacloud.Note) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Note ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)
		foundNote, err := client.Note.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundNote.ID != rs.Primary.ID {
			return fmt.Errorf("Note not found")
		}

		*note = *foundNote

		return nil
	}
}

func testAccCheckSakuraCloudNoteDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_note" {
			continue
		}

		_, err := client.Note.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Note still exists")
		}
	}

	return nil
}

const testAccCheckSakuraCloudNoteConfig_basic = `
resource "sakuracloud_note" "foobar" {
    name = "mynote"
    content = "content"
    description = "Note from TerraForm for SAKURA CLOUD"
    tags = ["hoge" , "hoge2"]
}`

const testAccCheckSakuraCloudNoteConfig_update = `
resource "sakuracloud_note" "foobar" {
    name = "mynote_upd"
    content = "content_upd"
}`
