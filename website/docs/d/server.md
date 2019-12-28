---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_server"
subcategory: "Compute"
description: |-
  Get information about an existing Server.
---

# Data Source: sakuracloud_server

Get information about an existing Server.

## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.


---

A `filter` block supports the following:

* `condition` - (Optional) One or more name/values pairs used for filtering. There are several valid keys, for a full reference, check out finding section in the [SakuraCloud API reference](https://developer.sakura.ad.jp/cloud/api/1.1/).
* `id` - (Optional) The resource id on SakuraCloud used for filtering.
* `names` - (Optional) The resource names on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.
* `tags` - (Optional) The resource tags on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.

---

A `condition` block supports the following:

* `name` - (Required) The name of the target field. This value is case-sensitive.
* `values` - (Required) The values of the condition. If multiple values ​​are specified, they combined as AND condition.


## Attribute Reference

* `id` - The id of the Server.
* `cdrom_id` - The id of the CD-ROM attached to the server.
* `commitment` - The policy of how to allocate virtual CPUs to the server. This will be one of [`standard`/`dedicatedcpu`].
* `core` - The number of virtual CPUs.
* `description` - The description of the server.
* `disks` - A list of disk id connected to the server.
* `dns_servers` - A list of IP address of DNS server in the zone.
* `gateway` - The IP address of the gateway used by server.
* `icon_id` - The icon id attached to the server.
* `interface_driver` - The driver name of network interface. This will be one of [`virtio`/`e1000`].
* `ip_address` - The IP address assigned to the server.
* `memory` - The size of memory in GiB.
* `name` - The name of the server.
* `netmask` - The bit length of the subnet assigned to the server.
* `network_address` - The network address which the `ip_address` belongs.
* `network_interface` - A list of `network_interface` blocks as defined below.
* `private_host_id` - The id of the private host which the server is assigned.
* `private_host_name` - The name of the private host which the server is assigned.
* `tags` - Any tags assigned to the server.
* `zone` - The name of zone that the server is in (e.g. `is1a`,`tk1a`).


---

A `network_interface` block exports the following:

* `mac_address` - The MAC address.
* `packet_filter_id` - The id of the packet filter attached to the network interface.
* `upstream` - The upstream type or upstream switch id. This will be one of [`shared`/`disconnect`/`<switch id>`].



