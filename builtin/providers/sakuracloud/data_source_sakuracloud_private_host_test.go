package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudPrivateHostDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudPrivateHostDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostBase,
				Check:  testAccCheckSakuraCloudPrivateHostDataSourceID("sakuracloud_private_host.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePrivateHostConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceID("data.sakuracloud_private_host.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_private_host.foobar", "description", "description_test"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourcePrivateHostConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostDataSourceNotExists("data.sakuracloud_private_host.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudPrivateHostDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find PrivateHost data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("PrivateHost data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudPrivateHostDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found PrivateHost data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudPrivateHostDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "tk1a"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_private_host" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.PrivateHost.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("PrivateHost still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourcePrivateHostBase = `
resource "sakuracloud_private_host" "foobar" {
    name = "name_test"
    description = "description_test"
}
`

var testAccCheckSakuraCloudDataSourcePrivateHostConfig = `
resource "sakuracloud_private_host" "foobar" {
    name = "name_test"
    description = "description_test"
}
data "sakuracloud_private_host" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`

var testAccCheckSakuraCloudDataSourcePrivateHostConfig_NotExists = `
resource "sakuracloud_private_host" "foobar" {
    name = "name_test"
    description = "description_test"
}
data "sakuracloud_private_host" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`
