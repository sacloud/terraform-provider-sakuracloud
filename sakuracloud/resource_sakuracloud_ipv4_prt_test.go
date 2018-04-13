package sakuracloud

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/sacloud/libsacloud/sacloud"
	"testing"
)

const (
	envTestDomain = "SAKURACLOUD_TEST_DOMAIN"
)

var (
	testDomain string
)

func TestAccResourceSakuraCloudIPv4Prt(t *testing.T) {
	var ip sacloud.IPAddress
	if domain, ok := os.LookupEnv(envTestDomain); ok {
		testDomain = domain
	} else {
		t.Skipf("ENV %q is requilred. skip", envTestDomain)
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudIPv4PrtDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckSakuraCloudIPv4PrtConfig_basic, testDomain, testDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudIPv4PrtExists("sakuracloud_ipv4_prt.foobar", &ip),
					resource.TestCheckResourceAttr(
						"sakuracloud_ipv4_prt.foobar", "hostname", fmt.Sprintf("terraform-test-domain01.%s", testDomain)),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckSakuraCloudIPv4PrtConfig_update, testDomain, testDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudIPv4PrtExists("sakuracloud_ipv4_prt.foobar", &ip),
					resource.TestCheckResourceAttr(
						"sakuracloud_ipv4_prt.foobar", "hostname", fmt.Sprintf("terraform-test-domain02.%s", testDomain)),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudIPv4PrtExists(n string, ip *sacloud.IPAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No IPv4Prt ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)

		foundIPv4Prt, err := client.IPAddress.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundIPv4Prt.IPAddress != rs.Primary.ID {
			return errors.New("IPv4Prt not found")
		}
		if foundIPv4Prt.HostName == "" {
			return errors.New("IPv4Prt hostname is empty")
		}

		*ip = *foundIPv4Prt

		return nil
	}
}

func testAccCheckSakuraCloudIPv4PrtDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_ipv4_prt" {
			continue
		}

		ip, err := client.IPAddress.Read(rs.Primary.ID)

		if err == nil && ip.HostName != "" {
			return errors.New("IPv4Prt still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudIPv4PrtConfig_basic = `
data sakuracloud_dns "dns" {
  name_selectors = ["%s"]
}

resource sakuracloud_dns_record "record01" {
  dns_id = "${data.sakuracloud_dns.dns.id}"
  name   = "terraform-test-domain01"
  type   = "A"
  value  = "${sakuracloud_server.server.ipaddress}"
}

resource sakuracloud_server "server" {
  name = "server"
  graceful_shutdown_timeout = 5
}

resource "sakuracloud_ipv4_prt" "foobar" {
  ipaddress = "${sakuracloud_server.server.ipaddress}"
  hostname  = "terraform-test-domain01.%s"
}
`

var testAccCheckSakuraCloudIPv4PrtConfig_update = `
data sakuracloud_dns "dns" {
  name_selectors = ["%s"]
}

resource sakuracloud_dns_record "record01" {
  dns_id = "${data.sakuracloud_dns.dns.id}"
  name   = "terraform-test-domain02"
  type   = "A"
  value  = "${sakuracloud_server.server.ipaddress}"
}

resource sakuracloud_server "server" {
  name = "server"
  graceful_shutdown_timeout = 5
}

resource "sakuracloud_ipv4_prt" "foobar" {
  ipaddress = "${sakuracloud_server.server.ipaddress}"
  hostname  = "terraform-test-domain02.%s"
}
`
