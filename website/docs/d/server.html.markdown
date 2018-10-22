---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_server"
sidebar_current: "docs-sakuracloud-datasource-server"
description: |-
  Get information on a SakuraCloud Server.
---

# sakuracloud\_server

Use this data source to retrieve information about a SakuraCloud Server.

## Example Usage

```hcl
data "sakuracloud_server" "foobar" {
  name_selectors = ["foobar"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The list of names to filtering.
 * `tag_selectors` - (Optional) The list of tags to filtering.
 * `filter` - (Optional) The map of filter key and value.
 * `zone` - (Optional) The ID of the zone.

## Attributes Reference

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
* `packet_filter_ids` - The ID list of the Packet Filter connected to Server.
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

