package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TTestAccResourceSakuraCloudInternet_Basic(t *testing.T) {
	var internet sacloud.Internet
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudInternetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
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
		},
	})
}

func TTestAccResourceSakuraCloudInternet_Update(t *testing.T) {
	var internet sacloud.Internet
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudInternetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
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
			resource.TestStep{
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
				),
			},
		},
	})
}

func TTestAccResourceSakuraCloudInternet_WithServer(t *testing.T) {
	var internet sacloud.Internet
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudInternetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudInternetConfig_with_server,
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
			resource.TestStep{
				Config: testAccCheckSakuraCloudInternetConfig_with_server,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetExists("sakuracloud_internet.foobar", &internet),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "name", "myinternet"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "nw_mask_len", "28"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "band_width", "100"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "server_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_internet.foobar", "nw_ipaddresses.#", "11"),
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
			return fmt.Errorf("No Internet ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)

		foundInternet, err := client.Internet.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundInternet.ID != rs.Primary.ID {
			return fmt.Errorf("Internet not found")
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

		_, err := client.Internet.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Internet still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudInternetConfig_basic = `
resource "sakuracloud_internet" "foobar" {
    name = "myinternet"
}`

var testAccCheckSakuraCloudInternetConfig_update = `
resource "sakuracloud_internet" "foobar" {
    name = "myinternet_upd"
    band_width = 500
}`

var testAccCheckSakuraCloudInternetConfig_with_server = `
resource "sakuracloud_server" "foobar" {
    name = "myserver"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
    base_interface = "${sakuracloud_internet.foobar.switch_id}"
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
    name = "myinternet"
}
`
