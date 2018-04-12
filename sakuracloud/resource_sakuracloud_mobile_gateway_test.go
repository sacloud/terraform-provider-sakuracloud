package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/sacloud/libsacloud/sacloud"
	"testing"
)

func TestAccSakuraCloudMobileGateway(t *testing.T) {
	var mgs sacloud.MobileGateway
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudMobileGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudMobileGatewayConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudMobileGatewayExists("sakuracloud_mobile_gateway.foobar", &mgs),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "name", "name_before"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "description", "description_before"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "internet_connection", "true"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "tags.0", "hoge1"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "tags.1", "hoge2"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "private_ipaddress", "192.168.11.101"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "private_nw_mask_len", "24"),
					resource.TestCheckResourceAttrPair(
						"sakuracloud_mobile_gateway.foobar", "icon_id",
						"sakuracloud_icon.foobar", "id",
					),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway_static_route.r1", "prefix", "192.168.10.0/24"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway_static_route.r1", "next_hop", "192.168.11.1"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway_static_route.r2", "prefix", "192.168.10.0/25"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway_static_route.r2", "next_hop", "192.168.11.2"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway_static_route.r3", "prefix", "192.168.10.0/26"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway_static_route.r3", "next_hop", "192.168.11.3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudMobileGatewayConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudMobileGatewayExists("sakuracloud_mobile_gateway.foobar", &mgs),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "name", "name_after"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "description", "description_after"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "internet_connection", "false"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "tags.0", "hoge1_after"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "tags.1", "hoge2_after"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "private_ipaddress", ""),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "private_nw_mask_len", "0"),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "switch_id", ""),
					resource.TestCheckResourceAttr("sakuracloud_mobile_gateway.foobar", "icon_id", ""),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudMobileGatewayExists(n string, mgs *sacloud.MobileGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No MobileGateway ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)

		foundMobileGateway, err := client.MobileGateway.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundMobileGateway.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("MobileGateway not found")
		}

		*mgs = *foundMobileGateway

		return nil
	}
}

func testAccCheckSakuraCloudMobileGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_mobile_gateway" {
			continue
		}

		_, err := client.MobileGateway.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("MobileGateway still exists")
		}
	}

	return nil
}

func TestAccImportSakuraCloudMobileGateway(t *testing.T) {
	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %#v", s)
		}
		expects := map[string]string{
			"name":                      "name_before",
			"private_ipaddress":         "192.168.11.101",
			"private_nw_mask_len":       "24",
			"description":               "description_before",
			"tags.0":                    "hoge1",
			"tags.1":                    "hoge2",
			"graceful_shutdown_timeout": "60",
		}

		if err := compareStateMulti(s[0], expects); err != nil {
			return err
		}
		return stateNotEmptyMulti(s[0], "switch_id", "icon_id")
	}

	resourceName := "sakuracloud_mobile_gateway.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudMobileGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudMobileGatewayConfig_basic,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateCheck:  checkFn,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCheckSakuraCloudMobileGatewayConfig_basic = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
}
resource "sakuracloud_mobile_gateway" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    private_ipaddress = "192.168.11.101"
    private_nw_mask_len = 24
	internet_connection = true
    name = "name_before"
    description = "description_before"
    tags = ["hoge1" , "hoge2"]
    icon_id = "${sakuracloud_icon.foobar.id}"
}

resource sakuracloud_mobile_gateway_static_route "r1" {
	mobile_gateway_id = "${sakuracloud_mobile_gateway.foobar.id}"
    prefix = "192.168.10.0/24"
    next_hop = "192.168.11.1"
}

resource sakuracloud_mobile_gateway_static_route "r2" {
	mobile_gateway_id = "${sakuracloud_mobile_gateway.foobar.id}"
    prefix = "192.168.10.0/25"
    next_hop = "192.168.11.2"
}

resource sakuracloud_mobile_gateway_static_route "r3" {
	mobile_gateway_id = "${sakuracloud_mobile_gateway.foobar.id}"
    prefix = "192.168.10.0/26"
    next_hop = "192.168.11.3"
}
resource "sakuracloud_icon" "foobar" {
  name = "myicon"
  base64content = "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"
}
`

const testAccCheckSakuraCloudMobileGatewayConfig_update = `
resource "sakuracloud_switch" "sw" {
    name = "sw"
}
resource "sakuracloud_mobile_gateway" "foobar" {
    internet_connection = false
    name = "name_after"
    description = "description_after"
    tags = ["hoge1_after" , "hoge2_after"]
}`
