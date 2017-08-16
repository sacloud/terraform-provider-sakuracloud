package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"regexp"
	"testing"
)

func TestAccResourceSakuraCloudServer(t *testing.T) {
	var server sacloud.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			{
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
						"sakuracloud_server.foobar", "interface_driver", "virtio"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "nic", "shared"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_nics.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "1"),
					resource.TestMatchResourceAttr("sakuracloud_server.foobar",
						"base_nw_ipaddress",
						regexp.MustCompile(".+")), // should be not empty
					resource.TestCheckResourceAttrPair(
						"sakuracloud_server.foobar", "icon_id",
						"sakuracloud_icon.foobar", "id",
					),
				),
			},
			{
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
						"sakuracloud_server.foobar", "interface_driver", "e1000"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "tags.0", "tag2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "nic", "shared"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_nics.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "1"),
					resource.TestMatchResourceAttr("sakuracloud_server.foobar",
						"base_nw_ipaddress",
						regexp.MustCompile(".+")), // should be not empty
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "icon_id", ""),
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
			{
				Config: testAccCheckSakuraCloudServerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_server.foobar", "base_interface"),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_server.foobar", "additional_nics"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "1"),
				),
			},
			{
				Config: testAccCheckSakuraCloudServerConfig_swiched_NIC_added,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_server.foobar", "base_interface"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_nics.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "2"),
				),
			},
			{
				Config: testAccCheckSakuraCloudServerConfig_swiched_NIC_updated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckNoResourceAttr(
						"sakuracloud_server.foobar", "base_interface"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_nics.#", "3"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "4"),
				),
			},
			{
				Config: testAccCheckSakuraCloudServerConfig_nw_nothing,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributesWithoutSharedInterface(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "nic", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_nics.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "2"),
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
			{
				Config: testAccCheckSakuraCloudServerConfig_with_packet_filter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "packet_filter_ids.0", ""),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "packet_filter_ids.#", "2"),
				),
			},
			{
				Config: testAccCheckSakuraCloudServerConfig_with_packet_filter_add,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "name", "myserver_upd"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "packet_filter_ids.#", "2"),
				),
			},
			{
				Config: testAccCheckSakuraCloudServerConfig_with_packet_filter_upd,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "packet_filter_ids.#", "1"),
				),
			},
			{
				Config: testAccCheckSakuraCloudServerConfig_with_packet_filter_del,
			},
			{
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
			{
				Config: testAccCheckSakuraCloudServerConfig_with_blank_disk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
				),
			},
		},
	})
}

func TestAccSakuraCloudServer_EditConnect_With_Same_Switch(t *testing.T) {
	var server sacloud.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudServerConfig_connect_same_sw_before,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "nic", "shared"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_nics.#", "1"),
					func(s *terraform.State) error {
						if server.Interfaces[1].GetSwitch() == nil {
							return errors.New("Server.Interfaces[1].Switch is nil")
						}
						if server.Interfaces[1].GetSwitch().GetID() == 0 {
							return errors.New("Server.Interfaces[1].Switch has invalid ID")
						}
						return nil
					},
				),
			},
			{
				Config: testAccCheckSakuraCloudServerConfig_connect_same_sw_after,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_nics.#", "1"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_nics.0", ""),
					func(s *terraform.State) error {
						if server.Interfaces[0].GetSwitch() == nil {
							return errors.New("Server.Interfaces[0].Switch is nil")
						}
						if server.Interfaces[0].GetSwitch().GetID() == 0 {
							return errors.New("Server.Interfaces[0].Switch has invalid ID")
						}
						if server.Interfaces[0].GetSwitch().Scope == sacloud.ESCopeShared {
							return errors.New("Server.Interfaces[0].Switch is connecting to shared segment")
						}
						return nil
					},
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
			return errors.New("No Server ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)

		foundServer, err := client.Server.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundServer.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("Server not found")
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
			return fmt.Errorf("Bad server NIC status. Server must have NIC which are disconnected the shared segment.: %#v", server.Interfaces)
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
			return errors.New("Server still exists")
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
    tags = ["tag1"]
    icon_id = "${sakuracloud_icon.foobar.id}"
}

resource "sakuracloud_icon" "foobar" {
  name = "myicon"
  base64content = "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"
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
    tags = ["tag2"]
    interface_driver = "e1000"
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
    additional_nics = [""]
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
    additional_nics = ["","",""]
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
    nic = ""
    additional_nics = [""]
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
    nic = "shared"
    additional_nics = [""]
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
    nic = "shared"
    additional_nics = [""]
    packet_filter_ids = ["${sakuracloud_packet_filter.foobar1.id}" , "${sakuracloud_packet_filter.foobar2.id}"]
}

`

const testAccCheckSakuraCloudServerConfig_with_packet_filter_upd = `
resource "sakuracloud_packet_filter" "foobar1" {
    name = "mypacket_filter1"
    description = "PacketFilter from TerraForm for SAKURA CLOUD"
    expressions = {
    	protocol = "udp"
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
    	protocol = "udp"
    	source_nw = "0.0.0.0"
    	source_port = "0-65535"
    	dest_port = "80"
    	allow = true
    }
}
resource "sakuracloud_server" "foobar" {
    name = "myserver_upd"
    nic = "shared"
    additional_nics = [""]
    packet_filter_ids = ["${sakuracloud_packet_filter.foobar1.id}"]
}

`

const testAccCheckSakuraCloudServerConfig_with_packet_filter_del = `
resource "sakuracloud_server" "foobar" {
    name = "myserver_upd"
    nic = "shared"
    additional_nics = [""]
}`

const testAccCheckSakuraCloudServerConfig_with_blank_disk = `
resource "sakuracloud_server" "foobar" {
    name = "myserver_with_blank"
    nic = "shared"
    disks = ["${sakuracloud_disk.foobar.id}"]
}
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
}
`

const testAccCheckSakuraCloudServerConfig_connect_same_sw_before = `
resource "sakuracloud_switch" "foobar" {
    name = "foobar"
}
resource "sakuracloud_server" "foobar" {
    name = "foobar"
    nic = "shared"
    additional_nics = ["${sakuracloud_switch.foobar.id}"]
}
`

const testAccCheckSakuraCloudServerConfig_connect_same_sw_after = `
resource "sakuracloud_switch" "foobar" {
    name = "foobar"
}
resource "sakuracloud_server" "foobar" {
    name = "foobar"
    nic = "${sakuracloud_switch.foobar.id}"
    additional_nics = [""]
}
`
