---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_server"
sidebar_current: "docs-sakuracloud-resource-computing-server"
description: |-
  Provides a SakuraCloud Server resource. This can be used to create, update, and delete Servers.
---

# sakuracloud\_server

Provides a SakuraCloud Server resource. This can be used to create, update, and delete Servers.

## Example Usage

```hcl
# Create a new Server
resource "sakuracloud_server" "foobar" {
  name                = "foobar"
  # core              = 1
  # memory            = 1
  disks               = [sakuracloud_disk.foobar.id]
  # interface_driver  = "virtio"
  # nic               = "shared"
  
  # cdrom_id          = sakuracloud_cdrom.foobar.id
  # private_host_id   = sakuracloud_private_host.foobar.id
  
  # additional_nics   = [
  #  sakuracloud_switch.foobar.id,  # connect to switch
  #  "", # disconnected
  # ] 
  
  # packet_filter_ids = [
  #  sakuracloud_packet_filter.foobar.id, # for primary NIC
  #  sakuracloud_packet_filter.foobar.id, # for secondly NIC
  #]
 
  # only when nic != shared
  # ipaddress       = "192.2.0.1"
  # nw_mask_len     = 24
  # gateway         = "192.2.0.254" 
  
  description       = "description"
  tags              = ["foo", "bar"]
  
  # For "Modify Disk" API
  hostname          = "your-host-name"
  password          = ""
  ssh_key_ids       = ["<your-ssh-key-id>"]
  note_ids          = ["<your-note-id"]
  disable_pw_auth   = true
}

# Create a new Disk
resource "sakuracloud_disk" "foobar" {
  name              = "foobar"
  plan              = "ssd"
  connector         = "virtio"
  size              = 20
  source_archive_id = data.sakuracloud_archive.ubuntu.id
  # or
  # source_disk_id  = "<your-disk-id>"
  
  description       = "description"
  tags              = ["foo", "bar"]
}

# Source archive
data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `core` - (Optional) The number of cores (default:`1`).
* `memory` - (Optional) The size of memory (unit:`GB`, default:`1`).
* `disks` - (Optional) The ID list of the Disks connected to Server.
* `interface_driver` - (Optional) The name of network interface driver.  
Valid value is one of the following: [ "virtio" (default) / "e1000"]
* `nic` - (Optional) The primary NIC's connection destination.  
Valid value is one of the following: [ "shared" (default) / <Switch ID> / "" (disconnected) ]
* `cdrom_id` - (Optional) The ID of the CD-ROM inserted to Server.
* `private_host_id` - (Optional) The ID of the Private Host to which the Server belongs.
* `additional_nics` - (Optional) The ID list of the Switches connected to NICs (excluding primary NIC) of Server.  
Valid values are one of the following: [ <Switch ID> / "" (disconnected) ]
* `packet_filter_ids` - (Optional) The ID list of the Packet Filters connected to Server.
* `ipaddress` - (Optional) The IP address of primary NIC to set with `"Modify Disk"` API.
* `gateway` - (Optional) Default gateway address of the Server to set with `"Modify Disk"` API.	 
* `nw_mask_len` - (Optional) Network mask length of the Server to set with `"Modify Disk"` API.
* `hostname` - (Optional) The hostname to set with `"Modify Disk"` API.
* `password` - (Optional) The password of OS's administrator to set with `"Modify Disk"` API.
* `ssh_key_ids` - (Optional) The ID list of SSH Keys to set with `"Modify Disk"` API.
* `note_ids` - (Optional) The ID list of Notes (Startup-Scripts) to set with `"Modify Disk"` API.
* `disable_pw_auth` - (Optional) The flag of disable password auth via SSH.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the Server.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `core` - The number of cores.
* `memory` - The size of memory (unit:`GB`).
* `disks` - The ID list of the Disks connected to Server.
* `interface_driver` - The name of network interface driver.
* `nic` - The primary NIC's connection destination.
* `cdrom_id` - The ID of the CD-ROM inserted to Server.
* `private_host_id` - The ID of the Private Host to which the Server belongs.
* `private_host_name` - The name of the Private Host to which the Server belongs.
* `additional_nics` - The ID list of the Switches connected to NICs (excluding primary NIC) of Server.
* `packet_filter_ids` - The ID list of the Packet Filters connected to Server.
* `macaddresses` - The MAC address list of NICs connected to Server.
* `ipaddress` - The IP address of primary NIC.
* `dns_servers` - List of default DNS servers for the zone to which the Server belongs.
* `gateway` - Default gateway address of the Server.	 
* `nw_address` - The network address of the Server.
* `nw_mask_len` - Network mask length of the Server.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Servers can be imported using the Server ID.

```
$ terraform import sakuracloud_server.foobar <server_id>
```
