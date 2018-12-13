package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceZone_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceZoneBase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudZoneDataSourceID("data.sakuracloud_zone.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "name", "is1a"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "zone_id", "31001"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "description", "石狩第1ゾーン"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "region_id", "310"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "region_name", "石狩"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "dns_servers.0", "133.242.0.3"),
					resource.TestCheckResourceAttr("data.sakuracloud_zone.foobar", "dns_servers.1", "133.242.0.4"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudZoneDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Zone data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Zone data source ID not set")
		}
		return nil
	}
}

var testAccCheckSakuraCloudDataSourceZoneBase = `data "sakuracloud_zone" "foobar" { name = "is1a"}`
