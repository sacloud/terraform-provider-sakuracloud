package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceVPCRouter_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudVPCRouterDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouterBase(name),
				Check:  testAccCheckSakuraCloudVPCRouterDataSourceID("sakuracloud_vpc_router.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouterConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceID("data.sakuracloud_vpc_router.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_vpc_router.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_vpc_router.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_vpc_router.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_vpc_router.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_vpc_router.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_vpc_router.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouterConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceID("data.sakuracloud_vpc_router.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouter_NameSelector_Exists(name, randString1, randString2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceID("data.sakuracloud_vpc_router.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouter_TagSelector_Exists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceID("data.sakuracloud_vpc_router.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouterConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceNotExists("data.sakuracloud_vpc_router.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouterConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceNotExists("data.sakuracloud_vpc_router.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouter_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceNotExists("data.sakuracloud_vpc_router.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouter_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceNotExists("data.sakuracloud_vpc_router.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudVPCRouterDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find VPCRouter data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("VPCRouter data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudVPCRouterDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found VPCRouter data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudVPCRouterDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_vpc_router" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.VPCRouter.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("VPCRouter still exists")
		}
	}

	return nil
}

func testAccCheckSakuraCloudDataSourceVPCRouterBase(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
`, name)
}

func testAccCheckSakuraCloudDataSourceVPCRouterConfig(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_vpc_router" "foobar" {
    filter = {
	name = "Name"
	values = ["%s"]
    }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceVPCRouterConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_vpc_router" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`, name)
}

func testAccCheckSakuraCloudDataSourceVPCRouterConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_vpc_router" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`, name)
}

func testAccCheckSakuraCloudDataSourceVPCRouterConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_vpc_router" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`, name)
}

func testAccCheckSakuraCloudDataSourceVPCRouter_NameSelector_Exists(name, p1, p2 string) string {
	return fmt.Sprintf(`
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_vpc_router" "foobar" {
    name_selectors = ["%s", "%s"]
}`, name, p1, p2)
}

var testAccCheckSakuraCloudDataSourceVPCRouter_NameSelector_NotExists = `
data "sakuracloud_vpc_router" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}
`

func testAccCheckSakuraCloudDataSourceVPCRouter_TagSelector_Exists(name string) string {
	return fmt.Sprintf(`
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "%s"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_vpc_router" "foobar" {
	tag_selectors = ["tag1","tag2","tag3"]
}`, name)
}

var testAccCheckSakuraCloudDataSourceVPCRouter_TagSelector_NotExists = `
data "sakuracloud_vpc_router" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
}`
