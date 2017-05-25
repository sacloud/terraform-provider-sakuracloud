package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"testing"
)

func TestAccResourceSakuraCloudPrivateHost(t *testing.T) {
	var private_host sacloud.PrivateHost
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudPrivateHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudPrivateHostConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostExists("sakuracloud_private_host.foobar", &private_host),
					resource.TestCheckResourceAttr(
						"sakuracloud_private_host.foobar", "name", "before"),
					resource.TestCheckResourceAttr(
						"sakuracloud_private_host.foobar", "description", "before"),
					resource.TestCheckResourceAttr(
						"sakuracloud_private_host.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttrPair(
						"sakuracloud_private_host.foobar", "id",
						"sakuracloud_server.foobar", "private_host_id",
					),
				),
			},
			{
				Config: testAccCheckSakuraCloudPrivateHostConfig_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPrivateHostExists("sakuracloud_private_host.foobar", &private_host),
					resource.TestCheckResourceAttr(
						"sakuracloud_private_host.foobar", "name", "after"),
					resource.TestCheckResourceAttr(
						"sakuracloud_private_host.foobar", "description", "after"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_private_host.foobar", "tags"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudPrivateHostExists(n string, private_host *sacloud.PrivateHost) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No PrivateHost ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)
		originalZone := client.Zone
		client.Zone = "tk1a"
		defer func() { client.Zone = originalZone }()

		foundPrivateHost, err := client.PrivateHost.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundPrivateHost.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("PrivateHost not found")
		}

		*private_host = *foundPrivateHost

		return nil
	}
}

func testAccCheckSakuraCloudPrivateHostDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "tk1a"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_private_host" {
			continue
		}

		_, err := client.PrivateHost.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("PrivateHost still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudPrivateHostConfig_Basic = `
resource "sakuracloud_server" "foobar" {
    name = "myswitch"
    private_host_id = "${sakuracloud_private_host.foobar.id}"
    zone = "tk1a"
}
resource "sakuracloud_private_host" "foobar" {
    name = "before"
    description = "before"
    tags = ["tag1", "tag2"]
    zone = "tk1a"

}`

var testAccCheckSakuraCloudPrivateHostConfig_Update = `
resource "sakuracloud_private_host" "foobar" {
    name = "after"
    description = "after"
    zone = "tk1a"
}
`
