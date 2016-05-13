package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccSakuraCloudDisk_Basic(t *testing.T) {
	var disk sacloud.Disk
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDiskConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskExists("sakuracloud_disk.foobar", &disk),
					testAccCheckSakuraCloudDiskAttributes(&disk),
					resource.TestCheckResourceAttr(
						"sakuracloud_disk.foobar", "name", "mydisk"),
					resource.TestCheckResourceAttr(
						"sakuracloud_disk.foobar", "disable_pw_auth", "false"),
				),
			},
		},
	})
}

func TestAccSakuraCloudDisk_Update(t *testing.T) {
	var disk sacloud.Disk
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDiskConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskExists("sakuracloud_disk.foobar", &disk),
					testAccCheckSakuraCloudDiskAttributes(&disk),
					resource.TestCheckResourceAttr(
						"sakuracloud_disk.foobar", "name", "mydisk"),
					resource.TestCheckResourceAttr(
						"sakuracloud_disk.foobar", "disable_pw_auth", "false"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudDiskConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDiskExists("sakuracloud_disk.foobar", &disk),
					testAccCheckSakuraCloudDiskAttributes(&disk),
					resource.TestCheckResourceAttr(
						"sakuracloud_disk.foobar", "name", "mydisk"),
					resource.TestCheckResourceAttr(
						"sakuracloud_disk.foobar", "disable_pw_auth", "true"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudDiskExists(n string, disk *sacloud.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Disk ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)
		originalZone := client.Zone
		client.Zone = "tk1v"
		defer func() { client.Zone = originalZone }()

		foundDisk, err := client.Disk.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundDisk.ID != rs.Primary.ID {
			return fmt.Errorf("Disk not found")
		}

		*disk = *foundDisk

		return nil
	}
}

func testAccCheckSakuraCloudDiskAttributes(disk *sacloud.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if disk.Connection != sacloud.DiskConnectionVirtio {
			return fmt.Errorf("Bad disk connection: %v", disk.Connection)
		}

		if disk.Plan.ID.String() != sacloud.DiskPlanSSD.ID.String() {
			return fmt.Errorf("Bad disk plan: %v", disk.Plan)
		}
		return nil
	}
}

func testAccCheckSakuraCloudDiskDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "tk1v"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_disk" {
			continue
		}

		_, err := client.Disk.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Disk still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDiskConfig_basic = `
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_name = "Ubuntu Server 14"
    description = "Disk from TerraForm for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
    zone = "tk1v"
}`

var testAccCheckSakuraCloudDiskConfig_update = `
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_name = "Ubuntu Server 14"
    description = "Disk from TerraForm for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
    disable_pw_auth = true
    zone = "tk1v"
}`
