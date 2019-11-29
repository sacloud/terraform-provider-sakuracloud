package sakuracloud

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func TestAccResourceSakuraCloudServer(t *testing.T) {
	var server sacloud.Server
	resource.ParallelTest(t, resource.TestCase{
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
						"sakuracloud_server.foobar", "hostname", "myserver"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "password", "p@ssw0rd"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "ssh_key_ids.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "ssh_key_ids.0", "100000000000"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "disable_pw_auth", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "note_ids.#", "2"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "note_ids.0", "100000000000"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "nic", "shared"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "additional_nics.#", "0"),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "macaddresses.#", "1"),
					resource.TestMatchResourceAttr("sakuracloud_server.foobar",
						"ipaddress",
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
						"nw_address",
						regexp.MustCompile(".+")), // should be not empty
					resource.TestCheckResourceAttr("sakuracloud_server.foobar", "icon_id", ""),
				),
			},
		},
	})
}

func TestAccSakuraCloudServer_EditConnections(t *testing.T) {
	var server sacloud.Server
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudServerConfig_swiched_NIC_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "nic", "shared"),
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
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "nic", "shared"),
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
					resource.TestCheckResourceAttr(
						"sakuracloud_server.foobar", "nic", "shared"),
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
						"sakuracloud_server.foobar", "nic", "disconnect"),
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
	resource.ParallelTest(t, resource.TestCase{
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

func TestAccSakuraCloudServer_With_PrivateHost(t *testing.T) {

	privateHostID, ok := os.LookupEnv("SAKURACLOUD_PRIVATE_HOST_ID")
	if !ok {
		t.Log("Private host ID($SAKURACLOUD_PRIVATE_HOST_ID) is empty. Skip this test.")
		return
	}

	var server sacloud.Server
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckSakuraCloudServerConfig_with_private_host_template, privateHostID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					resource.TestCheckResourceAttr("sakuracloud_server.foobar", "private_host_id", privateHostID),
					resource.TestMatchResourceAttr("sakuracloud_server.foobar",
						"private_host_name",
						regexp.MustCompile(".+")), // should not empty
				),
			},
		},
	})
}

func TestAccSakuraCloudServer_With_BlankDisk(t *testing.T) {
	var server sacloud.Server
	resource.ParallelTest(t, resource.TestCase{
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
	resource.ParallelTest(t, resource.TestCase{
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
						if server.Interfaces[1].SwitchID.IsEmpty() {
							return errors.New("Server.Interfaces[1].Switch is nil")
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
						if server.Interfaces[0].SwitchID.IsEmpty() {
							return errors.New("Server.Interfaces[0].Switch is nil")
						}
						if server.Interfaces[0].SwitchScope == types.Scopes.Shared {
							return errors.New("Server.Interfaces[0].Switch is connecting to shared segment")
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccSakuraCloudServer_NIC_CustomDiff(t *testing.T) {
	var server sacloud.Server
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudServerConfig_nic_custom_diff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestMatchResourceAttr("sakuracloud_server.foobar",
						"ipaddress",
						regexp.MustCompile(".+")), // should be not empty
				),
			},
		},
	})
}

func TestAccSakuraCloudServer_NIC_CustomDiffReference(t *testing.T) {
	var server sacloud.Server
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudServerConfig_nic_custom_diff_reference,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestMatchResourceAttr("sakuracloud_server.foobar",
						"ipaddress",
						regexp.MustCompile(".+")), // should be not empty
				),
			},
			{
				Config: testAccCheckSakuraCloudServerConfig_nic_custom_diff_reference_upd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					testAccCheckSakuraCloudServerAttributes(&server),
					resource.TestMatchResourceAttr("sakuracloud_server.foobar",
						"ipaddress",
						regexp.MustCompile(".+")), // should be not empty
				),
			},
		},
	})
}

func TestAccSakuraCloudServer_Switched_eth0(t *testing.T) {
	var server sacloud.Server
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudServerConfig_switched_eth0,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudServerExists("sakuracloud_server.foobar", &server),
					resource.TestCheckResourceAttr("sakuracloud_server.foobar", "ipaddress", "192.168.0.2"),
					resource.TestCheckResourceAttr("sakuracloud_server.foobar", "nw_mask_len", "24"),
					resource.TestCheckResourceAttr("sakuracloud_server.foobar", "gateway", "192.168.0.1"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudServerExists(n string, server *sacloud.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no Server ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		serverOp := sacloud.NewServerOp(client)
		zone := rs.Primary.Attributes["zone"]

		foundServer, err := serverOp.Read(context.Background(), zone, types.StringID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if foundServer.ID.String() != rs.Primary.ID {
			return fmt.Errorf("not found Server: %s", rs.Primary.ID)
		}

		*server = *foundServer
		return nil
	}
}

func testAccCheckSakuraCloudServerAttributes(server *sacloud.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if !server.InstanceStatus.IsUp() {
			return fmt.Errorf("unexpecterd server status: status=%v", server.InstanceStatus)
		}

		if len(server.Interfaces) == 0 {
			return errors.New("unexpecterd server NIC status: interfaces is nil")
		}

		if server.Interfaces[0].SwitchID.IsEmpty() || server.Interfaces[0].SwitchScope != types.Scopes.Shared {
			return fmt.Errorf("unexpected server NIC status: %#v", server.Interfaces[0])
		}

		return nil
	}
}

func testAccCheckSakuraCloudServerAttributesWithoutSharedInterface(server *sacloud.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if !server.InstanceStatus.IsUp() {
			return fmt.Errorf("unexpected server status: status=%v", server.InstanceStatus)
		}

		if len(server.Interfaces) == 0 || !server.Interfaces[0].SwitchID.IsEmpty() {
			return fmt.Errorf("unexpected server NIC status. %#v", server.Interfaces)
		}

		return nil
	}
}

func testAccCheckSakuraCloudServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	serverOp := sacloud.NewServerOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_server" {
			continue
		}

		zone := rs.Primary.Attributes["zone"]
		_, err := serverOp.Read(context.Background(), zone, types.StringID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("still exists Server:%s", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCheckSakuraCloudServerConfig_basic = `
data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu"
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

    hostname = "myserver"
    password = "p@ssw0rd"
    ssh_key_ids = ["100000000000", "200000000000"]
    disable_pw_auth = true
    note_ids = ["100000000000", "200000000000"]

    graceful_shutdown_timeout = 10
}

resource "sakuracloud_icon" "foobar" {
  name = "myicon"
  base64content = "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"
}
`

const testAccCheckSakuraCloudServerConfig_update = `
data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu"
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
    graceful_shutdown_timeout = 10
}
`

const testAccCheckSakuraCloudServerConfig_swiched_NIC_basic = `
resource "sakuracloud_server" "foobar" {
    name = "myserver"
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["tag1"]

    hostname = "myserver"
    password = "p@ssw0rd"
    ssh_key_ids = ["100000000000", "200000000000"]
    disable_pw_auth = true
    note_ids = ["100000000000", "200000000000"]

    graceful_shutdown_timeout = 2
}
`

const testAccCheckSakuraCloudServerConfig_swiched_NIC_added = `
resource "sakuracloud_server" "foobar" {
    name = "myserver"
    description = "Server from TerraForm for SAKURA CLOUD"
    additional_nics = [sakuracloud_switch.sw1.id]
    graceful_shutdown_timeout = 2
}
resource "sakuracloud_switch" "sw1" {
  name = "terraform-test-switch1"
}
`
const testAccCheckSakuraCloudServerConfig_swiched_NIC_updated = `
resource "sakuracloud_server" "foobar" {
    name = "myserver"
    description = "Server from TerraForm for SAKURA CLOUD"
    additional_nics = [sakuracloud_switch.sw1.id, sakuracloud_switch.sw2.id, sakuracloud_switch.sw3.id]
    graceful_shutdown_timeout = 2
}
resource "sakuracloud_switch" "sw1" {
  name = "terraform-test-switch1"
}
resource "sakuracloud_switch" "sw2" {
  name = "terraform-test-switch2"
}
resource "sakuracloud_switch" "sw3" {
  name = "terraform-test-switch3"
}
`

const testAccCheckSakuraCloudServerConfig_nw_nothing = `
resource "sakuracloud_server" "foobar" {
    name = "myserver"
    description = "Server from TerraForm for SAKURA CLOUD"
    nic = "disconnect"
    additional_nics = [sakuracloud_switch.sw1.id]
    graceful_shutdown_timeout = 2
}
resource "sakuracloud_switch" "sw1" {
  name = "terraform-test-switch1"
}
`

const testAccCheckSakuraCloudServerConfig_with_packet_filter = `
resource "sakuracloud_packet_filter" "foobar" {
  name = "terraform-test-packetfilter"
  expressions {
  	protocol = "tcp"
  	source_network = "0.0.0.0"
  	source_port = "0-65535"
  	destination_port = "80"
  	allow = true
  }
}
resource "sakuracloud_server" "foobar" {
  name = "terraform-test-server"
  nic = "shared"
  additional_nics = [sakurackoud_switch.foobar.id]
  packet_filter_ids = ["" , "${sakuracloud_packet_filter.foobar2.id}"]
  graceful_shutdown_timeout = 10
}

resource "sakuracloud_switch" "foobar" {
  name = "terraform-test-switch"
}
`

const testAccCheckSakuraCloudServerConfig_with_packet_filter_add = `
resource "sakuracloud_packet_filter" "foobar1" {
    name = "mypacket_filter1"
    description = "PacketFilter from TerraForm for SAKURA CLOUD"
    expressions {
    	protocol = "tcp"
    	source_network = "0.0.0.0"
    	source_port = "0-65535"
    	destination_port = "80"
    	allow = true
    }
}
resource "sakuracloud_packet_filter" "foobar2" {
    name = "mypacket_filter2"
    description = "PacketFilter from TerraForm for SAKURA CLOUD"
    expressions {
    	protocol = "tcp"
    	source_network = "0.0.0.0"
    	source_port = "0-65535"
    	destination_port = "80"
    	allow = true
    }
}
resource "sakuracloud_server" "foobar" {
    name = "myserver_upd"
    nic = "shared"
    additional_nics = [""]
    packet_filter_ids = ["${sakuracloud_packet_filter.foobar1.id}" , "${sakuracloud_packet_filter.foobar2.id}"]
    graceful_shutdown_timeout = 10
}

`

const testAccCheckSakuraCloudServerConfig_with_packet_filter_upd = `
resource "sakuracloud_packet_filter" "foobar1" {
    name = "mypacket_filter1"
    description = "PacketFilter from TerraForm for SAKURA CLOUD"
    expressions {
    	protocol = "udp"
    	source_network = "0.0.0.0"
    	source_port = "0-65535"
    	destination_port = "80"
    	allow = true
    }
}
resource "sakuracloud_packet_filter" "foobar2" {
    name = "mypacket_filter2"
    description = "PacketFilter from TerraForm for SAKURA CLOUD"
    expressions {
    	protocol = "udp"
    	source_network = "0.0.0.0"
    	source_port = "0-65535"
    	destination_port = "80"
    	allow = true
    }
}
resource "sakuracloud_server" "foobar" {
    name = "myserver_upd"
    nic = "shared"
    additional_nics = [""]
    packet_filter_ids = ["${sakuracloud_packet_filter.foobar1.id}"]
    graceful_shutdown_timeout = 10
}

`

const testAccCheckSakuraCloudServerConfig_with_packet_filter_del = `
resource "sakuracloud_server" "foobar" {
    name = "myserver_upd"
    nic = "shared"
    additional_nics = [""]
    graceful_shutdown_timeout = 10
}`

const testAccCheckSakuraCloudServerConfig_with_blank_disk = `
resource "sakuracloud_server" "foobar" {
    name = "myserver_with_blank"
    nic = "shared"
    disks = ["${sakuracloud_disk.foobar.id}"]
    graceful_shutdown_timeout = 10
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
    graceful_shutdown_timeout = 10
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
    graceful_shutdown_timeout = 10
}
`

const testAccCheckSakuraCloudServerConfig_with_private_host_template = `
resource "sakuracloud_server" "foobar" {
    name            = "myserver_with_private_host"
    private_host_id = "%s"
    graceful_shutdown_timeout = 10
}
`

const testAccCheckSakuraCloudServerConfig_nic_custom_diff = `
resource "sakuracloud_server" "foobar" {
    name      = "foobar"
    nic       = "shared"
    ipaddress = ""
    graceful_shutdown_timeout = 10
}
`

const testAccCheckSakuraCloudServerConfig_nic_custom_diff_reference = `
resource "sakuracloud_server" "foobar" {
    name      = "foobar"
    graceful_shutdown_timeout = 10
}
resource sakuracloud_simple_monitor "foobar" {
  target = "${sakuracloud_server.foobar.ipaddress}"

  health_check {
    protocol   = "ping"
  }

  notify_email_enabled = true
  enabled              = true
}
`

const testAccCheckSakuraCloudServerConfig_nic_custom_diff_reference_upd = `
resource "sakuracloud_server" "foobar" {
    name        = "foobar"
    nic         = "shared"
    ipaddress   = ""
    gateway     = ""
    nw_mask_len = ""
    graceful_shutdown_timeout = 10
}
resource sakuracloud_simple_monitor "foobar" {
  target = "${sakuracloud_server.foobar.ipaddress}"

  health_check {
    protocol   = "ping"
  }

  notify_email_enabled = true
  enabled              = true
}
`

const testAccCheckSakuraCloudServerConfig_switched_eth0 = `
data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu"
}
resource "sakuracloud_disk" "foobar" {
    name = "mydisk"
    source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}
resource "sakuracloud_switch" "foobar" {
    name = "foobar"
}
resource "sakuracloud_server" "foobar" {
    name        = "foobar"
    disks       = ["${sakuracloud_disk.foobar.id}"]
    nic         = "${sakuracloud_switch.foobar.id}"
    ipaddress   = "192.168.0.2"
    nw_mask_len = 24
    gateway     = "192.168.0.1"
    graceful_shutdown_timeout = 10
}
`

const testAccCheckSakuraCloudServerConfig_display_ipaddress = `
resource sakuracloud_switch "sw" {
  count = 3
  name = "sakuracloud_test_switch"
}

resource sakuracloud_server "switched" {
  name              = "sakuracloud_test_connect_to_switched"
  nic               = "${sakuracloud_switch.sw.0.id}"
  display_ipaddress = "192.2.0.1"

  additional_nics                = ["${sakuracloud_switch.sw.1.id}", "${sakuracloud_switch.sw.2.id}"]
  additional_display_ipaddresses = ["192.2.1.1", "192.2.2.1"]
}
`
const testAccCheckSakuraCloudServerConfig_display_ipaddress_upd = `
resource sakuracloud_switch "sw" {
  count = 3
  name = "sakuracloud_test_switch"
}

resource sakuracloud_server "switched" {
  name              = "sakuracloud_test_connect_to_switched"
  nic               = "${sakuracloud_switch.sw.0.id}"
  display_ipaddress = "192.2.0.2"
  additional_nics   = ["${sakuracloud_switch.sw.1.id}", "${sakuracloud_switch.sw.2.id}"]
  additional_display_ipaddresses = ["192.2.1.2", "192.2.2.2"]
}
`
