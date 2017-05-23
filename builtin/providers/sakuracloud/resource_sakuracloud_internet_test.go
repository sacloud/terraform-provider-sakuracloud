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

func TestAccResourceSakuraCloudInternet(t *testing.T) {
	var internet sacloud.Internet
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudInternetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudInternetConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetExists("sakuracloud_internet.foobar", &internet),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "name", "myinternet"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "nw_mask_len", "28"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "band_width", "100"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "server_ids.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "nw_ipaddresses.#", "11"),
				),
			},
			{
				Config: testAccCheckSakuraCloudInternetConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetExists("sakuracloud_internet.foobar", &internet),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "name", "myinternet_upd"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "nw_mask_len", "28"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "band_width", "500"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "server_ids.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "nw_ipaddresses.#", "11"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "enable_ipv6", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "ipv6_prefix_len", "64"),
				),
			},
			{
				Config: testAccCheckSakuraCloudInternetConfig_with_server,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetExists("sakuracloud_internet.foobar", &internet),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "name", "myinternet_upd"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "nw_mask_len", "28"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "band_width", "500"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "server_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "nw_ipaddresses.#", "11"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "enable_ipv6", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "ipv6_prefix_len", "64"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudInternetExists(n string, internet *sacloud.Internet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Internet ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)

		foundInternet, err := client.Internet.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundInternet.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("Internet not found")
		}

		*internet = *foundInternet

		return nil
	}
}

func testAccCheckSakuraCloudInternetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_internet" {
			continue
		}

		_, err := client.Internet.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Internet still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudInternetConfig_basic = `
resource "sakuracloud_internet" "foobar" {
    name = "myinternet"
}`

var testAccCheckSakuraCloudInternetConfig_update = `
resource "sakuracloud_server" "foobar" {
    name = "myserver"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
    nic = "${sakuracloud_internet.foobar.switch_id}"
    base_nw_ipaddress = "${sakuracloud_internet.foobar.nw_ipaddresses.0}"
    base_nw_gateway = "${sakuracloud_internet.foobar.nw_gateway}"
    base_nw_mask_len = "${sakuracloud_internet.foobar.nw_mask_len}"
}
data "sakuracloud_archive" "ubuntu" {
    filter = {
	name = "Name"
	values = ["Ubuntu Server 16"]
    }
}
resource "sakuracloud_disk" "foobar"{
    name = "mydisk"
    source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}

resource "sakuracloud_internet" "foobar" {
    name = "myinternet_upd"
    band_width = 500
    enable_ipv6 = true
}`

var testAccCheckSakuraCloudInternetConfig_with_server = `
resource "sakuracloud_server" "foobar" {
    name = "myserver"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
    nic = "${sakuracloud_internet.foobar.switch_id}"
    base_nw_ipaddress = "${sakuracloud_internet.foobar.nw_ipaddresses.0}"
    base_nw_gateway = "${sakuracloud_internet.foobar.nw_gateway}"
    base_nw_mask_len = "${sakuracloud_internet.foobar.nw_mask_len}"
}
data "sakuracloud_archive" "ubuntu" {
    filter = {
	name = "Name"
	values = ["Ubuntu Server 16"]
    }
}
resource "sakuracloud_disk" "foobar"{
    name = "mydisk"
    source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}
resource "sakuracloud_internet" "foobar" {
    name = "myinternet_upd"
    band_width = 500
    enable_ipv6 = true
}
`
