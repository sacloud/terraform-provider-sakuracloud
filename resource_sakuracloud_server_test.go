package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccSakuraCloudServer_Basic(t *testing.T) {
	var server sacloud.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "core", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "memory", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "disks.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "shared_interface", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "switched_interfaces.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "mac_addresses.#", "1"),
				),
			},
		},
	})
}

func TestAccSakuraCloudServer_Update(t *testing.T) {
	var server sacloud.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "core", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "memory", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "disks.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "shared_interface", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "switched_interfaces.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "mac_addresses.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "core", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "memory", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "disks.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "shared_interface", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "switched_interfaces.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "mac_addresses.#", "1"),
				),
			},
		},
	})
}

func TestAccSakuraCloudServer_EditConnections(t *testing.T) {
	var server sacloud.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "shared_interface", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "switched_interfaces.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "mac_addresses.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_swiched_NIC_added,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "shared_interface", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "switched_interfaces.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "mac_addresses.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_swiched_NIC_updated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "shared_interface", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "switched_interfaces.#", "3"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "mac_addresses.#", "4"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_nw_nothing,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributesWithoutSharedInterface(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "shared_interface", "false"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "switched_interfaces.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "mac_addresses.#", "1"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudServerExists(n string, server *sacloud.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Server ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)
		originalZone := client.Zone
		client.Zone = "tk1v"
		defer func() { client.Zone = originalZone }()

		foundServer, err := client.Server.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundServer.ID != rs.Primary.ID {
			return fmt.Errorf("Server not found")
		}

		*server = *foundServer

		return nil
	}
}

func testAccCheckSakuraCloudServerAttributes(server *sacloud.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if !server.Instance.IsUp() {
			return fmt.Errorf("Bad server status. Server must be running.: %v", server.Instance.Status)
		}

		if len(server.Interfaces) == 0 ||
			server.Interfaces[0].Switch == nil ||
			server.Interfaces[0].Switch.Scope != sacloud.ESCopeShared {
			return fmt.Errorf("Bad server NIC status. Server must have are connected to the shared segment.: %v", server)
		}

		return nil
	}
}

func testAccCheckSakuraCloudServerAttributesWithoutSharedInterface(server *sacloud.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if !server.Instance.IsUp() {
			return fmt.Errorf("Bad server status. Server must be running.: %v", server.Instance.Status)
		}

		if len(server.Interfaces) == 0 || server.Interfaces[0].Switch != nil {
			return fmt.Errorf("Bad server NIC status. Server must have NIC which are disconnected the shared segment.: %v", server)
		}

		return nil
	}
}

func testAccCheckSakuraCloudServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "tk1v"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_server" {
			continue
		}

		_, err := client.Server.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Server still exists")
		}
	}

	return nil
}

const testAccCheckSakuraCloudServerConfig_basic = `
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_name = "Ubuntu Server 14"
    zone = "tk1v"
}

resource "sakuracloud_server" "foobar" {
    name = "myserver"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
    zone = "tk1v"
}
`

const testAccCheckSakuraCloudServerConfig_update = `
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_name = "Ubuntu Server 14"
    zone = "tk1v"
}

resource "sakuracloud_server" "foobar" {
    name = "myserver"
    disks = ["${sakuracloud_disk.foobar.id}"]
    core = 2
    memory = 2
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
    zone = "tk1v"
}
`

const testAccCheckSakuraCloudServerConfig_swiched_NIC_added = `
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_name = "Ubuntu Server 14"
    zone = "tk1v"
}

resource "sakuracloud_server" "foobar" {
    name = "myserver"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "Server from TerraForm for SAKURA CLOUD"
    switched_interfaces = [""]
    tags = ["@virtio-net-pci"]
    zone = "tk1v"
}
`
const testAccCheckSakuraCloudServerConfig_swiched_NIC_updated = `
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_name = "Ubuntu Server 14"
    zone = "tk1v"
}

resource "sakuracloud_server" "foobar" {
    name = "myserver"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "Server from TerraForm for SAKURA CLOUD"
    switched_interfaces = ["","",""]
    tags = ["@virtio-net-pci"]
    zone = "tk1v"
}
`

const testAccCheckSakuraCloudServerConfig_nw_nothing = `
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_name = "Ubuntu Server 14"
    zone = "tk1v"
}

resource "sakuracloud_server" "foobar" {
    name = "myserver"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "Server from TerraForm for SAKURA CLOUD"
    shared_interface = false
    tags = ["@virtio-net-pci"]
    zone = "tk1v"
}
`
