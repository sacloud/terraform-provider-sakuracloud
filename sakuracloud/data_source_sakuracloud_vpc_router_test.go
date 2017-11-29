package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccSakuraCloudDataSourceVPCRouter_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudVPCRouterDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouterBase,
				Check:  testAccCheckSakuraCloudVPCRouterDataSourceID("sakuracloud_vpc_router.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceID("data.sakuracloud_vpc_router.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_vpc_router.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_vpc_router.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_vpc_router.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_vpc_router.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_vpc_router.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_vpc_router.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouterConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceID("data.sakuracloud_vpc_router.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouter_NameSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceID("data.sakuracloud_vpc_router.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouter_TagSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceID("data.sakuracloud_vpc_router.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouterConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterDataSourceNotExists("data.sakuracloud_vpc_router.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceVPCRouterConfig_With_NotExists_Tag,
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

var testAccCheckSakuraCloudDataSourceVPCRouterBase = `
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
`

var testAccCheckSakuraCloudDataSourceVPCRouterConfig = `
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_vpc_router" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`

var testAccCheckSakuraCloudDataSourceVPCRouterConfig_With_Tag = `
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_vpc_router" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`

var testAccCheckSakuraCloudDataSourceVPCRouterConfig_With_NotExists_Tag = `
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_vpc_router" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceVPCRouterConfig_NotExists = `
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_vpc_router" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceVPCRouter_NameSelector_Exists = `
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_vpc_router" "foobar" {
    name_selectors = ["name", "test"]
}
`
var testAccCheckSakuraCloudDataSourceVPCRouter_NameSelector_NotExists = `
data "sakuracloud_vpc_router" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}
`

var testAccCheckSakuraCloudDataSourceVPCRouter_TagSelector_Exists = `
resource sakuracloud_vpc_router "foobar" {
    plan = "standard"

    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_vpc_router" "foobar" {
	tag_selectors = ["tag1","tag2","tag3"]
}`

var testAccCheckSakuraCloudDataSourceVPCRouter_TagSelector_NotExists = `
data "sakuracloud_vpc_router" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
}`
