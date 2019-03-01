---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_server_connector"
sidebar_current: "docs-sakuracloud-resource-computing-sconnector"
description: |-
  Provides a SakuraCloud Server Connection resource. This can be used to create, update, and delete Servers.
---

# sakuracloud\_server

Provides a SakuraCloud Server Connection resource. This can be used to create, update, and delete Servers.

## Example Usage

```hcl
# Connect Disk or CDROM or PacketFilter to Server
resource "sakuracloud_server_connector" "foobar" {
  server_id = sakuracloud_server.foobar.id

  disks = [sakuracloud_disk.foobar.id]

  cdrom_id = sakuracloud_cdrom.foobar.id

  packet_filter_ids = [
    sakuracloud_packet_filter.foobar.id, # for primary NIC
    sakuracloud_packet_filter.foobar.id, # for secondly NIC
  ]
}

```

## Argument Reference

The following arguments are supported:

* `server_id` - (Required) The name of the resource.
* `disks` - (Optional) The ID list of the Disks connected to Server.
* `cdrom_id` - (Optional) The ID of the CD-ROM inserted to Server.
* `packet_filter_ids` - (Optional) The ID list of the Packet Filters connected to Server.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the Server.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `disks` - The ID list of the Disks connected to Server.
* `cdrom_id` - The ID of the CD-ROM inserted to Server.
* `packet_filter_ids` - The ID list of the Packet Filters connected to Server.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Server Connections can be imported using the Server ID.

```
$ terraform import sakuracloud_server_connector.foobar <server_id>
```
