package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourcePacketFilter_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudPacketFilterDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourcePacketFilterBase(name),
				Check:  testAccCheckSakuraCloudPacketFilterDataSourceID("sakuracloud_packet_filter.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePacketFilterConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPacketFilterDataSourceID("data.sakuracloud_packet_filter.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.#", "2"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.source_nw", "0.0.0.0"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.source_port", "0-65535"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.dest_port", "80"),
					resource.TestCheckResourceAttr("data.sakuracloud_packet_filter.foobar", "expressions.0.allow", "true"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePacketFilter_NameSelector_Exists(name, randString1, randString2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPacketFilterDataSourceID("sakuracloud_packet_filter.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePacketFilterConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPacketFilterDataSourceNotExists("data.sakuracloud_packet_filter.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourcePacketFilter_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudPacketFilterDataSourceNotExists("data.sakuracloud_packet_filter.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudPacketFilterDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find PacketFilter data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("PacketFilter data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudPacketFilterDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found PacketFilter data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudPacketFilterDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_packet_filter" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.PacketFilter.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("PacketFilter still exists")
		}
	}

	return nil
}

func testAccCheckSakuraCloudDataSourcePacketFilterBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_packet_filter" "foobar" {
    name = "%s"
    description = "description_test"
    expressions {
    	protocol = "tcp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
    expressions {
    	protocol = "udp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
}`, name)
}

func testAccCheckSakuraCloudDataSourcePacketFilterConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_packet_filter" "foobar" {
    name = "%s"
    description = "description_test"
    expressions {
    	protocol = "tcp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
    expressions {
    	protocol = "udp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
}
data "sakuracloud_packet_filter" "foobar" {
    filter {
	name = "Name"
	values = ["%s"]
    }
}`, name, name)
}

var testAccCheckSakuraCloudDataSourcePacketFilterConfig_NotExists = `
data "sakuracloud_packet_filter" "foobar" {
    filter {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`

func testAccCheckSakuraCloudDataSourcePacketFilter_NameSelector_Exists(name, p1, p2 string) string {
	return fmt.Sprintf(`
resource "sakuracloud_packet_filter" "foobar" {
    name = "%s"
    description = "description_test"
    expressions {
    	protocol = "tcp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
    expressions {
    	protocol = "udp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
}
data "sakuracloud_packet_filter" "foobar" {
    name_selectors = ["%s", "%s"]
}`, name, p1, p2)
}

var testAccCheckSakuraCloudDataSourcePacketFilter_NameSelector_NotExists = `
data "sakuracloud_packet_filter" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}
`
