package sakuracloud

import (
	"errors"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"fmt"
	"github.com/sacloud/libsacloud/sacloud"
	"testing"
)

func TestAccResourceSakuraCloudGSLBServer(t *testing.T) {
	var gslb sacloud.GSLB
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudGSLBServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudGSLBServerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBExists("sakuracloud_gslb.foobar", &gslb),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar0", "ipaddress", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar1", "ipaddress", "8.8.4.4"),
				),
			},
			{
				Config: testAccCheckSakuraCloudGSLBServerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBExists("sakuracloud_gslb.foobar", &gslb),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.0", "ipaddress", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.1", "ipaddress", "8.8.4.4"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.2", "ipaddress", "208.67.222.123"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.3", "ipaddress", "208.67.220.123"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudGSLBServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_gslb" {
			continue
		}

		_, err := client.GSLB.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("GSLB still exists")
		}
	}

	return nil
}

func TestAccImportSakuraCloudGSLBServer(t *testing.T) {
	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %#v", s)
		}
		expects := map[string]string{
			"ipaddress": "8.8.8.8",
			"enabled":   "true",
			"weight":    "1",
		}

		if err := compareStateMulti(s[0], expects); err != nil {
			return err
		}
		return stateNotEmptyMulti(s[0], "gslb_id")
	}

	resourceName := "sakuracloud_gslb_server.foobar0"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudGSLBServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudGSLBServerConfig_basic,
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

var testAccCheckSakuraCloudGSLBServerConfig_basic = `
variable "gslb_ip_list" {
    default = ["8.8.8.8","8.8.4.4"]
}
resource "sakuracloud_gslb" "foobar" {
    name = "terraform.io"
    health_check = {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    description = "GSLB from TerraForm for SAKURA CLOUD"
    tags = ["hoge1", "hoge2"]
}
resource "sakuracloud_gslb_server" "foobar0" {
    gslb_id = "${sakuracloud_gslb.foobar.id}"
    ipaddress = "${var.gslb_ip_list[0]}"
}
resource "sakuracloud_gslb_server" "foobar1" {
    gslb_id = "${sakuracloud_gslb.foobar.id}"
    ipaddress = "${var.gslb_ip_list[1]}"
}`

var testAccCheckSakuraCloudGSLBServerConfig_update = `
variable "gslb_ip_list" {
    default = ["8.8.8.8","8.8.4.4","208.67.222.123","208.67.220.123"]
}
resource "sakuracloud_gslb" "foobar" {
    name = "terraform.io"
    health_check = {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    description = "GSLB from TerraForm for SAKURA CLOUD"
    tags = ["hoge1", "hoge2"]
}
resource "sakuracloud_gslb_server" "foobar" {
    count = 4
    gslb_id = "${sakuracloud_gslb.foobar.id}"
    ipaddress = "${var.gslb_ip_list[count.index]}"

}`
