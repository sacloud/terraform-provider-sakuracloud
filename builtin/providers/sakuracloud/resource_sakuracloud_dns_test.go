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

func TestAccResourceSakuraCloudDNS(t *testing.T) {
	var dns sacloud.DNS
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudDNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDNSConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSExists("sakuracloud_dns.foobar", &dns),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns.foobar", "zone", "terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns.foobar", "description", "DNS from TerraForm for SAKURA CLOUD"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDNSConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSExists("sakuracloud_dns.foobar", &dns),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns.foobar", "zone", "terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns.foobar", "description", "DNS from TerraForm for SAKURA CLOUD_upd"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudDNSExists(n string, dns *sacloud.DNS) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No DNS ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)

		foundDNS, err := client.DNS.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundDNS.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("Record not found")
		}

		*dns = *foundDNS

		return nil
	}
}

func testAccCheckSakuraCloudDNSDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_dns" {
			continue
		}

		_, err := client.DNS.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("DNS still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDNSConfig_basic = `
resource "sakuracloud_dns" "foobar" {
    zone = "terraform.io"
    description = "DNS from TerraForm for SAKURA CLOUD"
    tags = ["hoge1"]
}`

var testAccCheckSakuraCloudDNSConfig_update = `
resource "sakuracloud_dns" "foobar" {
    zone = "terraform.io"
    description = "DNS from TerraForm for SAKURA CLOUD_upd"
    tags = ["hoge1" , "hoge2"]
}`
