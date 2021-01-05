// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func TestAccSakuraCloudDNSRecord_basic(t *testing.T) {
	resourceName1 := "sakuracloud_dns_record.foobar1"
	resourceName2 := "sakuracloud_dns_record.foobar2"

	zone := fmt.Sprintf("%s.com", randomName())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudDNSDestroy,
			testCheckSakuraCloudDNSRecordDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDNSRecord_basic, zone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "name", "www"),
					resource.TestCheckResourceAttr(resourceName1, "type", "A"),
					resource.TestCheckResourceAttr(resourceName1, "value", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName2, "name", "_sip._tls"),
					resource.TestCheckResourceAttr(resourceName2, "type", "SRV"),
					resource.TestCheckResourceAttr(resourceName2, "value", "www.sakura.ne.jp."),
					resource.TestCheckResourceAttr(resourceName2, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName2, "weight", "2"),
					resource.TestCheckResourceAttr(resourceName2, "port", "3"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDNSRecord_update, zone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "name", "www2"),
					resource.TestCheckResourceAttr(resourceName1, "type", "A"),
					resource.TestCheckResourceAttr(resourceName1, "value", "192.168.0.2"),
				),
			},
		},
	})
}

func TestAccSakuraCloudDNSRecord_withCount(t *testing.T) {
	resourceName := "sakuracloud_dns_record.foobar"
	zone := fmt.Sprintf("%s.com", randomName())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudDNSDestroy,
			testCheckSakuraCloudDNSRecordDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDNSRecord_withCount, zone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName+".0", "name", "www"),
					resource.TestCheckResourceAttr(resourceName+".0", "type", "A"),
					resource.TestCheckResourceAttr(resourceName+".0", "value", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName+".1", "name", "www"),
					resource.TestCheckResourceAttr(resourceName+".1", "type", "A"),
					resource.TestCheckResourceAttr(resourceName+".1", "value", "192.168.0.2"),
				),
			},
		},
	})
}

func testCheckSakuraCloudDNSRecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	dnsOp := sacloud.NewDNSOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_dns_record" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		dnsID := rs.Primary.Attributes["dns_id"]
		if dnsID != "" {
			dns, err := dnsOp.Read(context.Background(), sakuraCloudID(dnsID))
			if err != nil && !sacloud.IsNotFoundError(err) {
				return fmt.Errorf("resource still exists: DNS: %s", rs.Primary.ID)
			}
			if dns != nil {
				record := &sacloud.DNSRecord{
					Name:  rs.Primary.Attributes["name"],
					Type:  types.EDNSRecordType(rs.Primary.Attributes["type"]),
					RData: rs.Primary.Attributes["value"],
					TTL:   forceAtoI(rs.Primary.Attributes["ttl"]),
				}

				for _, r := range dns.Records {
					if isSameDNSRecord(r, record) {
						return fmt.Errorf("resource still exists: DNSRecord: %s", rs.Primary.ID)
					}
				}
			}
		}
	}

	return nil
}

var testAccSakuraCloudDNSRecord_basic = `
resource "sakuracloud_dns" "foobar" {
  zone        = "{{ .arg0 }}"
}

resource "sakuracloud_dns_record" "foobar1" {
  dns_id = sakuracloud_dns.foobar.id
  name   = "www"
  type   = "A"
  value  = "192.168.0.1"
}

resource "sakuracloud_dns_record" "foobar2" {
  dns_id   = sakuracloud_dns.foobar.id
  name     = "_sip._tls"
  type     = "SRV"
  value    = "www.sakura.ne.jp."
  priority = 1
  weight   = 2
  port     = 3
}
`

var testAccSakuraCloudDNSRecord_update = `
resource "sakuracloud_dns" "foobar" {
  zone = "{{ .arg0 }}"
}

resource "sakuracloud_dns_record" "foobar1" {
  dns_id = sakuracloud_dns.foobar.id
  name   = "www2"
  type   = "A"
  value  = "192.168.0.2"
}`

var testAccSakuraCloudDNSRecord_withCount = `
resource "sakuracloud_dns" "foobar" {
  zone = "{{ .arg0 }}"
}

variable "addresses" {
  default = ["192.168.0.1", "192.168.0.2"]
}

resource "sakuracloud_dns_record" "foobar" {
  count  = 2
  dns_id = sakuracloud_dns.foobar.id
  name   = "www"
  type   = "A"
  value  = var.addresses[count.index]
}`
