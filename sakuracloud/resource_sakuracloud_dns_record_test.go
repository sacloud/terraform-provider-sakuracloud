package sakuracloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func TestAccResourceSakuraCloudDNSRecord_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	zone := fmt.Sprintf("%s.%s.com", randString1, randString2)

	var dns sacloud.DNS
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDNSRecordConfig_basic(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSExists("sakuracloud_dns.foobar", &dns),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar", "name", "test1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar", "type", "A"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar", "value", "192.168.0.1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar1", "name", "_sip._tls"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar1", "type", "SRV"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar1", "value", "www.sakura.ne.jp."),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar1", "priority", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar1", "weight", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar1", "port", "3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDNSRecordConfig_update(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSExists("sakuracloud_dns.foobar", &dns),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar", "name", "test2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar", "type", "A"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar", "value", "192.168.0.2"),
				),
			},
		},
	})
}

func TestAccResourceSakuraCloudDNSRecord_With_Count(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	zone := fmt.Sprintf("%s.%s.com", randString1, randString2)

	var dns sacloud.DNS
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDNSRecordConfig_with_count(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSExists("sakuracloud_dns.foobar", &dns),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar.0", "name", "test"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar.0", "type", "A"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar.0", "value", "192.168.0.1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar.1", "name", "test"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar.1", "type", "A"),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns_record.foobar.1", "value", "192.168.0.2"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudDNSRecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	dnsOp := sacloud.NewDNSOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_dns_record" {
			continue
		}
		dnsID := rs.Primary.Attributes["dns_id"]
		if dnsID != "" {
			dns, err := dnsOp.Read(context.Background(), types.StringID(dnsID))
			if err == nil {
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

func testAccCheckSakuraCloudDNSRecordConfig_basic(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "DNS from TerraForm for SAKURA CLOUD"
  tags = ["hoge1"]
}

resource "sakuracloud_dns_record" "foobar" {
  dns_id = "${sakuracloud_dns.foobar.id}"
  name = "test1"
  type = "A"
  value = "192.168.0.1"
}

resource "sakuracloud_dns_record" "foobar1" {
  dns_id = "${sakuracloud_dns.foobar.id}"
  name = "_sip._tls"
  type = "SRV"
  value = "www.sakura.ne.jp."
  priority = 1
  weight = 2
  port = 3
}
`, zone)
}

func testAccCheckSakuraCloudDNSRecordConfig_update(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "DNS from TerraForm for SAKURA CLOUD"
  tags = ["hoge1"]
}

resource "sakuracloud_dns_record" "foobar" {
  dns_id = "${sakuracloud_dns.foobar.id}"
  name = "test2"
  type = "A"
  value = "192.168.0.2"
}`, zone)
}

func testAccCheckSakuraCloudDNSRecordConfig_with_count(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "DNS from TerraForm for SAKURA CLOUD"
  tags = ["hoge1"]
}
variable "ip_list" {
  default = ["192.168.0.1","192.168.0.2"]
}
resource "sakuracloud_dns_record" "foobar" {
  count = 2
  dns_id = "${sakuracloud_dns.foobar.id}"
  name = "test"
  type = "A"
  value = "${var.ip_list[count.index]}"
}`, zone)
}
