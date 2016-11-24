package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"testing"
)

func TestAccResourceSakuraCloudServer(t *testing.T) {
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
						"sakuracloud_server.foobar", "base_interface", "shared"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_interfaces.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "1"),
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
						"sakuracloud_server.foobar", "base_interface", "shared"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_interfaces.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "1"),
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
						"sakuracloud_server.foobar", "base_interface", "shared"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_interfaces.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_swiched_NIC_added,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "base_interface", "shared"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_interfaces.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_swiched_NIC_updated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "base_interface", "shared"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_interfaces.#", "3"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "4"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_nw_nothing,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributesWithoutSharedInterface(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "base_interface", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_interfaces.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "1"),
				),
			},
		},
	})
}

func TestAccSakuraCloudServer_ConnectPacketFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_with_packet_filter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "packet_filter_ids.0", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "packet_filter_ids.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_with_packet_filter_add,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "name", "myserver_upd"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "packet_filter_ids.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_with_packet_filter_upd,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "packet_filter_ids.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_with_packet_filter_del,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "packet_filter_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccSakuraCloudServer_With_BlankDisk(t *testing.T) {
	var server sacloud.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudServerConfig_with_blank_disk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
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

		foundServer, err := client.Server.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundServer.ID != toSakuraCloudID(rs.Primary.ID) {
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

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_server" {
			continue
		}

		_, err := client.Server.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("Server still exists")
		}
	}

	return nil
}

const testAccCheckSakuraCloudServerConfig_basic = `
data "sakuracloud_archive" "ubuntu" {
    filter = {
	name = "Name"
	values = ["Ubuntu Server 16"]
    }
}
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}

resource "sakuracloud_server" "foobar" {
    name = "myserver"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
}
`

const testAccCheckSakuraCloudServerConfig_update = `
data "sakuracloud_archive" "ubuntu" {
    filter = {
	name = "Name"
	values = ["Ubuntu Server 16"]
    }
}

resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}

resource "sakuracloud_server" "foobar" {
    name = "myserver_after"
    disks = ["${sakuracloud_disk.foobar.id}"]
    core = 2
    memory = 2
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
}
`

const testAccCheckSakuraCloudServerConfig_swiched_NIC_added = `
data "sakuracloud_archive" "ubuntu" {
    filter = {
	name = "Name"
	values = ["Ubuntu Server 16"]
    }
}

resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}

resource "sakuracloud_server" "foobar" {
    name = "myserver"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "Server from TerraForm for SAKURA CLOUD"
    additional_interfaces = [""]
    tags = ["@virtio-net-pci"]
}
`
const testAccCheckSakuraCloudServerConfig_swiched_NIC_updated = `
data "sakuracloud_archive" "ubuntu" {
    filter = {
	name = "Name"
	values = ["Ubuntu Server 16"]
    }
}

resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}

resource "sakuracloud_server" "foobar" {
    name = "myserver"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "Server from TerraForm for SAKURA CLOUD"
    additional_interfaces = ["","",""]
    tags = ["@virtio-net-pci"]
}
`

const testAccCheckSakuraCloudServerConfig_nw_nothing = `
data "sakuracloud_archive" "ubuntu" {
    filter = {
	name = "Name"
	values = ["Ubuntu Server 16"]
    }
}

resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}

resource "sakuracloud_server" "foobar" {
    name = "myserver"
    disks = ["${sakuracloud_disk.foobar.id}"]
    description = "Server from TerraForm for SAKURA CLOUD"
    base_interface = ""
    tags = ["@virtio-net-pci"]
}
`

const testAccCheckSakuraCloudServerConfig_with_packet_filter = `
resource "sakuracloud_packet_filter" "foobar2" {
    name = "mypacket_filter2"
    expressions = {
    	protocol = "tcp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
}
resource "sakuracloud_server" "foobar" {
    name = "myserver"
    base_interface = "shared"
    additional_interfaces = [""]
    packet_filter_ids = ["" , "${sakuracloud_packet_filter.foobar2.id}"]
}
`

const testAccCheckSakuraCloudServerConfig_with_packet_filter_add = `
resource "sakuracloud_packet_filter" "foobar1" {
    name = "mypacket_filter1"
    description = "PacketFilter from TerraForm for SAKURA CLOUD"
    expressions = {
    	protocol = "tcp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
}
resource "sakuracloud_packet_filter" "foobar2" {
    name = "mypacket_filter2"
    description = "PacketFilter from TerraForm for SAKURA CLOUD"
    expressions = {
    	protocol = "tcp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
}
resource "sakuracloud_server" "foobar" {
    name = "myserver_upd"
    base_interface = "shared"
    additional_interfaces = [""]
    packet_filter_ids = ["${sakuracloud_packet_filter.foobar1.id}" , "${sakuracloud_packet_filter.foobar2.id}"]
}

`

const testAccCheckSakuraCloudServerConfig_with_packet_filter_upd = `
resource "sakuracloud_packet_filter" "foobar1" {
    name = "mypacket_filter1"
    description = "PacketFilter from TerraForm for SAKURA CLOUD"
    expressions = {
    	protocol = "tcp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
}
resource "sakuracloud_packet_filter" "foobar2" {
    name = "mypacket_filter2"
    description = "PacketFilter from TerraForm for SAKURA CLOUD"
    expressions = {
    	protocol = "tcp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
}
resource "sakuracloud_server" "foobar" {
    name = "myserver_upd"
    base_interface = "shared"
    additional_interfaces = [""]
    packet_filter_ids = ["${sakuracloud_packet_filter.foobar1.id}"]
}

`

const testAccCheckSakuraCloudServerConfig_with_packet_filter_del = `
resource "sakuracloud_server" "foobar" {
    name = "myserver_upd"
    base_interface = "shared"
    additional_interfaces = [""]
}`

const testAccCheckSakuraCloudServerConfig_with_blank_disk = `
resource "sakuracloud_server" "foobar" {
    name = "myserver_with_blank"
    base_interface = "shared"
    disks = ["${sakuracloud_disk.foobar.id}"]
}
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
}
`
