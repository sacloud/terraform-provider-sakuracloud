package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccSakuraCloudGSLBServer_Basic(t *testing.T) {
	var gslb sacloud.GSLB
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudGSLBServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudGSLBServerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBExists("sakuracloud_gslb.foobar", &gslb),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.0", "ipaddress", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.1", "ipaddress", "8.8.4.4"),
				),
			},
		},
	})
}

func TestAccSakuraCloudGSLBServer_Update(t *testing.T) {
	var gslb sacloud.GSLB
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudGSLBServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudGSLBServerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBExists("sakuracloud_gslb.foobar", &gslb),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.0", "ipaddress", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb_server.foobar.1", "ipaddress", "8.8.4.4"),
				),
			},
			resource.TestStep{
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
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_gslb" {
			continue
		}

		_, err := client.GSLB.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("GSLB still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudGSLBServerConfig_basic = `
variable "gslb_ip_list" {
    default = "8.8.8.8,8.8.4.4"
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
    count = 2
    gslb_id = "${sakuracloud_gslb.foobar.id}"
    ipaddress = "${element(split("," , var.gslb_ip_list),count.index)}"
}`

var testAccCheckSakuraCloudGSLBServerConfig_update = `
variable "gslb_ip_list" {
    default = "8.8.8.8,8.8.4.4,208.67.222.123,208.67.220.123"
}
resource "sakuracloud_gslb" "foobar" {
    name = "terraform.io"
    health_check = {
        protocol = "https"
        delay_loop = 20
        host_header = "update.terraform.io"
        path = "/"
        status = "200"
    }
    description = "GSLB from TerraForm for SAKURA CLOUD"
    tags = ["hoge1", "hoge2"]
}
resource "sakuracloud_gslb_server" "foobar" {
    count = 4
    gslb_id = "${sakuracloud_gslb.foobar.id}"
    ipaddress = "${element(split("," , var.gslb_ip_list),count.index)}"

}`
