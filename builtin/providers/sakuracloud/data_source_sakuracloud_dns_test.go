package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudDNSDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudDNSDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceDNSBase,
				Check:  testAccCheckSakuraCloudDNSDataSourceID("sakuracloud_dns.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDNSConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSDataSourceID("data.sakuracloud_dns.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "zone", "test-terraform-sakuracloud.com"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "tags.2", "tag3"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceDNSConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSDataSourceID("data.sakuracloud_dns.foobar"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceDNSConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSDataSourceNotExists("data.sakuracloud_dns.foobar"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceDNSConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSDataSourceNotExists("data.sakuracloud_dns.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudDNSDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find DNS data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("DNS data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudDNSDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found DNS data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudDNSDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_dns" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.DNS.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("DNS still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceDNSBase = `
resource "sakuracloud_dns" "foobar" {
    zone = "test-terraform-sakuracloud.com"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}`

var testAccCheckSakuraCloudDataSourceDNSConfig = `
resource "sakuracloud_dns" "foobar" {
    zone = "test-terraform-sakuracloud.com"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_dns" "foobar" {
    filter = {
	name = "Name"
	values = ["test-terraform-sakuracloud.com"]
    }
}`

var testAccCheckSakuraCloudDataSourceDNSConfig_With_Tag = `
resource "sakuracloud_dns" "foobar" {
    zone = "test-terraform-sakuracloud.com"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_dns" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`

var testAccCheckSakuraCloudDataSourceDNSConfig_With_NotExists_Tag = `
resource "sakuracloud_dns" "foobar" {
    zone = "test-terraform-sakuracloud.com"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_dns" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceDNSConfig_NotExists = `
resource "sakuracloud_dns" "foobar" {
    zone = "test-terraform-sakuracloud.com"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_dns" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`
