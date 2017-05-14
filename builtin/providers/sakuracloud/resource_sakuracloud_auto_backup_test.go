package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"testing"
)

func TestAccResourceSakuraCloudAutoBackup(t *testing.T) {
	var autoBackup sacloud.AutoBackup
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudAutoBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudAutoBackupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudAutoBackupExists("sakuracloud_auto_backup.foobar", &autoBackup),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "name", "name_before"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.0", "wed"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.1", "thu"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "max_backup_num", "1"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "description", "description_before"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "tags.0", "hoge1"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "tags.1", "hoge2"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "zone", "is1b"),
				),
			},
			{
				Config: testAccCheckSakuraCloudAutoBackupConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudAutoBackupExists("sakuracloud_auto_backup.foobar", &autoBackup),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "name", "name_after"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.0", "thu"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "weekdays.1", "fri"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "max_backup_num", "2"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "description", "description_after"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "tags.0", "hoge1_after"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "tags.1", "hoge2_after"),
					resource.TestCheckResourceAttr("sakuracloud_auto_backup.foobar", "zone", "is1b"),
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
			return errors.New("No AutoBackup ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)
		originalZone := client.Zone
		client.Zone = "is1b"
		defer func() { client.Zone = originalZone }()

		foundAutoBackup, err := client.AutoBackup.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundAutoBackup.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("Resource not found")
		}

		*auto_backup = *foundAutoBackup

		return nil
	}
}

func testAccCheckSakuraCloudAutoBackupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "is1b"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_auto_backup" {
			continue
		}

		_, err := client.AutoBackup.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("AutoBackup still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudAutoBackupConfig_basic = `
resource "sakuracloud_disk" "disk" {
    name = "disk01"
    zone = "is1b"
}
resource "sakuracloud_auto_backup" "foobar" {
    name = "name_before"
    disk_id = "${sakuracloud_disk.disk.id}"
    weekdays = ["wed","thu"]
    max_backup_num = 1
    description = "description_before"
    tags = ["hoge1", "hoge2"]
    zone = "is1b"
}`

var testAccCheckSakuraCloudAutoBackupConfig_update = `
resource "sakuracloud_disk" "disk" {
    name = "disk01"
    zone = "is1b"
}
resource "sakuracloud_auto_backup" "foobar" {
    name = "name_after"
    disk_id = "${sakuracloud_disk.disk.id}"
    weekdays = ["thu","fri"]
    max_backup_num = 2
    description = "description_after"
    tags = ["hoge1_after", "hoge2_after"]
    zone = "is1b"
}`
