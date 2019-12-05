package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceDNS_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	zone := fmt.Sprintf("%s.%s.com", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudDNSDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceDNSBase(zone),
				Check:  testAccCheckSakuraCloudDNSDataSourceID("sakuracloud_dns.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDNSConfig(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSDataSourceID("data.sakuracloud_dns.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "zone", zone),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_dns.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDNSConfig_With_Tag(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSDataSourceID("data.sakuracloud_dns.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDNSConfig_NotExists(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSDataSourceNotExists("data.sakuracloud_dns.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDNSConfig_With_NotExists_Tag(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDNSDataSourceNotExists("data.sakuracloud_dns.foobar"),
				),
				Destroy: true,
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
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found DNS data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudDataSourceDNSBase(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}`, zone)
}

func testAccCheckSakuraCloudDataSourceDNSConfig(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_dns" "foobar" {
  filters {
	names = ["%s"]
  }
}`, zone, zone)
}

func testAccCheckSakuraCloudDataSourceDNSConfig_With_Tag(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_dns" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
}`, zone)
}

func testAccCheckSakuraCloudDataSourceDNSConfig_With_NotExists_Tag(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_dns" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
}`, zone)
}

func testAccCheckSakuraCloudDataSourceDNSConfig_NotExists(zone string) string {
	return fmt.Sprintf(`
resource "sakuracloud_dns" "foobar" {
  zone = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_dns" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, zone)
}
