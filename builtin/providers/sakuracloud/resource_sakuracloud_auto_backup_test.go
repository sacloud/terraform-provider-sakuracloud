package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccResourceSakuraCloudAutoBackup(t *testing.T) {
	var autoBackup sacloud.AutoBackup
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudAutoBackupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudAutoBackupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudAutoBackupExists("sakuracloud_auto_backup.foobar", &autoBackup),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "name", "name_before"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "backup_hour", "12"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.#", "3"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.0", "mon"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.1", "tue"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.2", "wed"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "max_backup_num", "1"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "description", "description_before"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "tags.0", "hoge1"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "tags.1", "hoge2"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "zone", "tk1a"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudAutoBackupConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudAutoBackupExists("sakuracloud_auto_backup.foobar", &autoBackup),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "name", "name_after"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "backup_hour", "18"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.0", "sat"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.1", "sun"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "max_backup_num", "2"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "description", "description_after"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "tags.0", "hoge1_after"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "tags.1", "hoge2_after"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "zone", "tk1a"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudAutoBackupExists(n string, auto_backup *sacloud.AutoBackup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AutoBackup ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)
		originalZone := client.Zone
		client.Zone = "tk1a"
		defer func() { client.Zone = originalZone }()

		foundAutoBackup, err := client.AutoBackup.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundAutoBackup.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found")
		}

		*auto_backup = *foundAutoBackup

		return nil
	}
}

func testAccCheckSakuraCloudAutoBackupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "tk1a"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_auto_backup" {
			continue
		}

		_, err := client.AutoBackup.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("AutoBackup still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudAutoBackupConfig_basic = `
resource "sakuracloud_disk" "disk" {
    name = "disk01"
    zone = "tk1a"
}
resource "sakuracloud_auto_backup" "foobar" {
    name = "name_before"
    disk_id = "${sakuracloud_disk.disk.id}"
    backup_hour = 12
    weekdays = ["mon","tue","wed"]
    max_backup_num = 1
    description = "description_before"
    tags = ["hoge1", "hoge2"]
    zone = "tk1a"
}`

var testAccCheckSakuraCloudAutoBackupConfig_update = `
resource "sakuracloud_disk" "disk" {
    name = "disk01"
    zone = "tk1a"
}
resource "sakuracloud_auto_backup" "foobar" {
    name = "name_after"
    disk_id = "${sakuracloud_disk.disk.id}"
    backup_hour = 18
    weekdays = ["sat","sun"]
    max_backup_num = 2
    description = "description_after"
    tags = ["hoge1_after", "hoge2_after"]
    zone = "tk1a"
}`
