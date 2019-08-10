package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceNote_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudNoteDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceNoteBase(name),
				Check:  testAccCheckSakuraCloudNoteDataSourceID("sakuracloud_note.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceID("data.sakuracloud_note.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "content", "content_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "class", "shell"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_note.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceID("data.sakuracloud_note.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudNoteDataSourceNotExists("data.sakuracloud_note.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceNoteConfig_With_NotExists_Tag(name),
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
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
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

func testAccCheckSakuraCloudDataSourceNoteBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_note" "foobar" {
  name = "%s"
  content = "content_test"
  tags = ["tag1","tag2","tag3"]
}`, name)
}

func testAccCheckSakuraCloudDataSourceNoteConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_note" "foobar" {
  name = "%s"
  content = "content_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_note" "foobar" {
  filters {
	names = ["%s"]
  }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceNoteConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_note" "foobar" {
  name = "%s"
  content = "content_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_note" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourceNoteConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_note" "foobar" {
  name = "%s"
  content = "content_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_note" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourceNoteConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_note" "foobar" {
  name = "%s"
  content = "content_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_note" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, name)
}
