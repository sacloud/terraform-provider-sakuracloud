// Copyright 2016-2020 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/sacloud"
)

const (
	envTestDomain = "SAKURACLOUD_TEST_DOMAIN"
)

var (
	testDomain string
)

func TestAccResourceSakuraCloudIPv4Ptr_basic(t *testing.T) {
	var ip sacloud.IPAddress
	if domain, ok := os.LookupEnv(envTestDomain); ok {
		testDomain = domain
	} else {
		t.Skipf("ENV %q is requilred. skip", envTestDomain)
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudIPv4PtrDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckSakuraCloudIPv4PtrConfig_basic, testDomain, testDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudIPv4PtrExists("sakuracloud_ipv4_ptr.foobar", &ip),
					resource.TestCheckResourceAttr(
						"sakuracloud_ipv4_ptr.foobar", "hostname", fmt.Sprintf("terraform-test-domain01.%s", testDomain)),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckSakuraCloudIPv4PtrConfig_update, testDomain, testDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudIPv4PtrExists("sakuracloud_ipv4_ptr.foobar", &ip),
					resource.TestCheckResourceAttr(
						"sakuracloud_ipv4_ptr.foobar", "hostname", fmt.Sprintf("terraform-test-domain02.%s", testDomain)),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudIPv4PtrExists(n string, ip *sacloud.IPAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No IPv4Ptr ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)

		foundIPv4Ptr, err := client.IPAddress.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundIPv4Ptr.IPAddress != rs.Primary.ID {
			return errors.New("IPv4Ptr not found")
		}
		if foundIPv4Ptr.HostName == "" {
			return errors.New("IPv4Ptr hostname is empty")
		}

		*ip = *foundIPv4Ptr

		return nil
	}
}

func testAccCheckSakuraCloudIPv4PtrDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_ipv4_ptr" {
			continue
		}

		ip, err := client.IPAddress.Read(rs.Primary.ID)

		if err == nil && ip.HostName != "" {
			return errors.New("IPv4Ptr still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudIPv4PtrConfig_basic = `
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

resource "sakuracloud_ipv4_ptr" "foobar" {
  ipaddress = "${sakuracloud_server.server.ipaddress}"
  hostname  = "terraform-test-domain01.%s"
}
`

var testAccCheckSakuraCloudIPv4PtrConfig_update = `
data sakuracloud_dns "dns" {
  name_selectors = ["%s"]
}

resource sakuracloud_dns_record "record01" {
  dns_id = "${data.sakuracloud_dns.dns.id}"
  name   = "terraform-test-domain02"
  type   = "A"
  value  = sakuracloud_server.server.ipaddress
}

resource sakuracloud_server "server" {
  name = "server"
  graceful_shutdown_timeout = 5
}

resource "sakuracloud_ipv4_ptr" "foobar" {
  ipaddress = "${sakuracloud_server.server.ipaddress}"
  hostname  = "terraform-test-domain02.%s"
}
`
