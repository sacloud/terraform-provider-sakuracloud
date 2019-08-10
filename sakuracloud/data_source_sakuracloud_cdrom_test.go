package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceCDROM_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudCDROMDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceCDROMConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudCDROMDataSourceID("data.sakuracloud_cdrom.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "name", "Ubuntu Server 18.04.2 LTS 64bit"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "size", "5"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.#", "5"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.0", "arch-64bit"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.1", "current-stable"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.2", "distro-ubuntu"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.3", "distro-ver-18.04.2"),
					resource.TestCheckResourceAttr("data.sakuracloud_cdrom.foobar", "tags.4", "os-linux"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceCDROM_NameSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudCDROMDataSourceID("data.sakuracloud_cdrom.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceCDROM_TagSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudCDROMDataSourceID("data.sakuracloud_cdrom.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceCDROMConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudCDROMDataSourceNotExists("data.sakuracloud_cdrom.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceCDROM_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudCDROMDataSourceNotExists("data.sakuracloud_cdrom.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceCDROM_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudCDROMDataSourceNotExists("data.sakuracloud_cdrom.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudCDROMDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find CDROM data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("CDROM data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudCDROMDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found CDROM data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudCDROMDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_cdrom" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.CDROM.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("CDROM still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceCDROMConfig = `
data "sakuracloud_cdrom" "foobar" {
  filters {
    conditions {
	  name = "Name"
	  values = ["Ubuntu server 18.04.2 LTS 64bit"]
    }
  }
}`

var testAccCheckSakuraCloudDataSourceCDROMConfig_NotExists = `
data "sakuracloud_cdrom" "foobar" {
  filters {
    conditions {
	  name = "Name"
	  values = ["xxxxxxxxxxxxxxxxxx"]
    }
  }
}`

var testAccCheckSakuraCloudDataSourceCDROM_NameSelector_Exists = `
data "sakuracloud_cdrom" "foobar" {
  filters {
    names = ["Ubuntu","server","18"]
  }
}
`
var testAccCheckSakuraCloudDataSourceCDROM_NameSelector_NotExists = `
data "sakuracloud_cdrom" "foobar" {
  filters {
    names = ["xxxxxxxxxx"]
  }
}
`

var testAccCheckSakuraCloudDataSourceCDROM_TagSelector_Exists = `
data "sakuracloud_cdrom" "foobar" {
  filters {
	tags = ["distro-ubuntu","os-unix"]
  }
}`

var testAccCheckSakuraCloudDataSourceCDROM_TagSelector_NotExists = `
data "sakuracloud_cdrom" "foobar" {
  filters {
	tags = ["xxxxxxxxxx"]
  }
}`
