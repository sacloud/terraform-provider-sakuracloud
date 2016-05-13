package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccSakuraCloudGSLB_Basic(t *testing.T) {
	var gslb sacloud.GSLB
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudGSLBDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudGSLBConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBExists("sakuracloud_gslb.foobar", &gslb),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "name", "terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "health_check.1802742300.protocol", "http"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "health_check.1802742300.delay_loop", "10"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "health_check.1802742300.host_header", "terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.0.ipaddress", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.1.ipaddress", "8.8.4.4"),
				),
			},
		},
	})
}

func TestAccSakuraCloudGSLB_Update(t *testing.T) {
	var gslb sacloud.GSLB
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudGSLBDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudGSLBConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBExists("sakuracloud_gslb.foobar", &gslb),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "name", "terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "health_check.1802742300.protocol", "http"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "health_check.1802742300.delay_loop", "10"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "health_check.1802742300.host_header", "terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.0.ipaddress", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.1.ipaddress", "8.8.4.4"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudGSLBConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBExists("sakuracloud_gslb.foobar", &gslb),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "name", "terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "health_check.755645870.protocol", "https"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "health_check.755645870.delay_loop", "20"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "health_check.755645870.host_header", "update.terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.0.ipaddress", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.1.ipaddress", "8.8.4.4"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudGSLBConfig_added,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBExists("sakuracloud_gslb.foobar", &gslb),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "name", "terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "health_check.755645870.protocol", "https"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "health_check.755645870.delay_loop", "20"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "health_check.755645870.host_header", "update.terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.#", "4"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.0.ipaddress", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.1.ipaddress", "8.8.4.4"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.2.ipaddress", "208.67.222.123"),
					resource.TestCheckResourceAttr(
						"sakuracloud_gslb.foobar", "servers.3.ipaddress", "208.67.220.123"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudGSLBExists(n string, gslb *sacloud.GSLB) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No GSLB ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)

		foundGSLB, err := client.GSLB.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundGSLB.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found")
		}

		*gslb = *foundGSLB

		return nil
	}
}

func testAccCheckSakuraCloudGSLBDestroy(s *terraform.State) error {
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

var testAccCheckSakuraCloudGSLBConfig_basic = `
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
    servers = {
      ipaddress = "8.8.8.8"
    }
    servers = {
      ipaddress = "8.8.4.4"
    }
}`

var testAccCheckSakuraCloudGSLBConfig_update = `
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
    servers = {
      ipaddress = "8.8.8.8"
    }
    servers = {
      ipaddress = "8.8.4.4"
    }
}`

var testAccCheckSakuraCloudGSLBConfig_added = `
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
    servers = {
      ipaddress = "8.8.8.8"
    }
    servers = {
      ipaddress = "8.8.4.4"
    }
    servers = {
      ipaddress = "208.67.222.123"
    }
    servers = {
      ipaddress = "208.67.220.123"
    }
}`
