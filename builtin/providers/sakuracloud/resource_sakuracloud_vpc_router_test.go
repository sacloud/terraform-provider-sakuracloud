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

func TestAccResourceSakuraCloudVPCRouter(t *testing.T) {
	var vpcRouter sacloud.VPCRouter
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudVPCRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudVPCRouterConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "name", "name_before"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "description", "description_before"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "tags.0", "hoge1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "tags.1", "hoge2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "plan", "standard"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "switch_id"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "vip"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "ipaddress1"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "ipaddress2"),
				),
			},
			{
				Config: testAccCheckSakuraCloudVPCRouterConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudVPCRouterExists("sakuracloud_vpc_router.foobar", &vpcRouter),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "name", "name_after"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "description", "description_after"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "tags.0", "hoge1_after"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "tags.1", "hoge2_after"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "plan", "standard"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "switch_id"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "vip"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "ipaddress1"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_vpc_router.foobar", "ipaddress2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_vpc_router.foobar", "syslog_host", "192.168.0.2"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudVPCRouterExists(n string, vpcRouter *sacloud.VPCRouter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No VPCRouter ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)

		foundVPCRouter, err := client.VPCRouter.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundVPCRouter.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("VPCRouter not found")
		}

		*vpcRouter = *foundVPCRouter

		return nil
	}
}

func testAccCheckSakuraCloudVPCRouterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_vpc_router" {
			continue
		}

		_, err := client.VPCRouter.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("VPCRouter still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudVPCRouterConfig_basic = `
resource "sakuracloud_vpc_router" "foobar" {
    name = "name_before"
    description = "description_before"
    tags = ["hoge1" , "hoge2"]
}`

var testAccCheckSakuraCloudVPCRouterConfig_update = `
resource "sakuracloud_vpc_router" "foobar" {
    name = "name_after"
    description = "description_after"
    tags = ["hoge1_after" , "hoge2_after"]
    syslog_host = "192.168.0.2"
}`
