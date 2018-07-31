package sakuracloud

import (
	"errors"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"

	"github.com/sacloud/libsacloud/sacloud"
)

func TestAccResourceSakuraCloudDNSRecord_Basic(t *testing.T) {
	var dns sacloud.DNS
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDNSRecordConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSExists("sakuracloud_dns.foobar", &dns),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns.foobar", "zone", "terraform.io"),
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
				Config: testAccCheckSakuraCloudDNSRecordConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSExists("sakuracloud_dns.foobar", &dns),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns.foobar", "zone", "terraform.io"),
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
	var dns sacloud.DNS
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDNSRecordConfig_with_count,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSExists("sakuracloud_dns.foobar", &dns),
					resource.TestCheckResourceAttr(
						"sakuracloud_dns.foobar", "zone", "terraform.io"),
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

var testAccCheckSakuraCloudDNSRecordConfig_basic = `
resource "sakuracloud_dns" "foobar" {
    zone = "terraform.io"
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
`

var testAccCheckSakuraCloudDNSRecordConfig_update = `
resource "sakuracloud_dns" "foobar" {
    zone = "terraform.io"
    description = "DNS from TerraForm for SAKURA CLOUD"
    tags = ["hoge1"]
}

resource "sakuracloud_dns_record" "foobar" {
    dns_id = "${sakuracloud_dns.foobar.id}"
    name = "test2"
    type = "A"
    value = "192.168.0.2"
}`

var testAccCheckSakuraCloudDNSRecordConfig_with_count = `

resource "sakuracloud_dns" "foobar" {
    zone = "terraform.io"
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
}`
