// Copyright 2016-2019 terraform-provider-sakuracloud authors
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
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

const (
	envTestDomain = "SAKURACLOUD_TEST_DOMAIN"
)

var (
	testDomain string
)

func TestAccResourceSakuraCloudIPv4Ptr(t *testing.T) {
	skipIfFakeModeEnabled(t)

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
		CheckDestroy: testCheckSakuraCloudIPv4PtrDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckSakuraCloudIPv4PtrConfig_basic, testDomain, testDomain),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudIPv4PtrExists("sakuracloud_ipv4_ptr.foobar", &ip),
					resource.TestCheckResourceAttr(
						"sakuracloud_ipv4_ptr.foobar", "hostname", fmt.Sprintf("terraform-test-domain01.%s", testDomain)),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckSakuraCloudIPv4PtrConfig_update, testDomain, testDomain),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudIPv4PtrExists("sakuracloud_ipv4_ptr.foobar", &ip),
					resource.TestCheckResourceAttr(
						"sakuracloud_ipv4_ptr.foobar", "hostname", fmt.Sprintf("terraform-test-domain02.%s", testDomain)),
				),
			},
		},
	})
}

func testCheckSakuraCloudIPv4PtrExists(n string, ip *sacloud.IPAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no IPv4Ptr ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		zone := rs.Primary.Attributes["zone"]
		ipAddrOp := sacloud.NewIPAddressOp(client)

		foundIPv4Ptr, err := ipAddrOp.Read(context.Background(), zone, rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundIPv4Ptr.IPAddress != rs.Primary.ID {
			return fmt.Errorf("not found IPv4Ptr: %s", rs.Primary.ID)
		}
		if foundIPv4Ptr.HostName == "" {
			return fmt.Errorf("hostname is empty IPv4Ptr: %s", foundIPv4Ptr.IPAddress)
		}

		*ip = *foundIPv4Ptr
		return nil
	}
}

func testCheckSakuraCloudIPv4PtrDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	ipAddrOp := sacloud.NewIPAddressOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_ipv4_ptr" {
			continue
		}

		zone := rs.Primary.Attributes["zone"]
		ip, err := ipAddrOp.Read(context.Background(), zone, rs.Primary.ID)

		if err == nil && ip.HostName != "" {
			return fmt.Errorf("still exists IPv4Ptr: %s", ip.IPAddress)
		}
	}

	return nil
}

var testAccCheckSakuraCloudIPv4PtrConfig_basic = `
data sakuracloud_dns "dns" {
  filters {
    names = ["%s"]
  }
}

resource sakuracloud_dns_record "record01" {
  dns_id = data.sakuracloud_dns.dns.id
  name   = "terraform-test-domain01"
  type   = "A"
  value  = sakuracloud_server.server.ip_address
}

resource sakuracloud_server "server" {
  name = "server"
  
  force_shutdown = true
}

resource "sakuracloud_ipv4_ptr" "foobar" {
  ip_address = sakuracloud_server.server.ip_address
  hostname  = "terraform-test-domain01.%s"
}
`

var testAccCheckSakuraCloudIPv4PtrConfig_update = `
data sakuracloud_dns "dns" {
  filters {
    names = ["%s"]
  }
}

resource sakuracloud_dns_record "record01" {
  dns_id = data.sakuracloud_dns.dns.id
  name   = "terraform-test-domain02"
  type   = "A"
  value  = sakuracloud_server.server.ip_address
}

resource sakuracloud_server "server" {
  name = "server"

  force_shutdown = true
}

resource "sakuracloud_ipv4_ptr" "foobar" {
  ip_address = sakuracloud_server.server.ip_address
  hostname  = "terraform-test-domain02.%s"
}
`
