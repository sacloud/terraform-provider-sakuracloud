---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_server"
subcategory: "Compute"
description: |-
  Get information about an existing Server.
---

# Data Source: sakuracloud_server

Get information about an existing Server.

## Example Usage

```hcl
data "sakuracloud_server" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```
## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.
* `zone` - (Optional) The name of zone that the Server is in (e.g. `is1a`, `tk1a`).

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
* `operator` - (Optional) The filtering operator. This must be one of following: `partial_match_and`/`exact_match_or`. Default: `partial_match_and`


## Attribute Reference

* `id` - The id of the Server.
* `cdrom_id` - The id of the CD-ROM attached to the server.
* `commitment` - The policy of how to allocate virtual CPUs to the server. This will be one of [`standard`/`dedicatedcpu`].
* `core` - The number of virtual CPUs.
* `cpu_model` - The model of CPU.
* `description` - The description of the Server.
* `disks` - A list of disk id connected to the server.
* `dns_servers` - A list of IP address of DNS server in the zone.
* `gateway` - The IP address of the gateway used by Server.
* `gpu` - The number of GPUs.
* `gpu_model` - The model of GPU.
* `hostname` - The hostname of the Server.
* `icon_id` - The icon id attached to the Server.
* `interface_driver` - The driver name of network interface. This will be one of [`virtio`/`e1000`].
* `ip_address` - The IP address assigned to the Server.
* `memory` - The size of memory in GiB.
* `name` - The name of the Server.
* `netmask` - The bit length of the subnet assigned to the Server.
* `network_address` - The network address which the `ip_address` belongs.
* `network_interface` - A list of `network_interface` blocks as defined below.
* `private_host_id` - The id of the private host which the server is assigned.
* `private_host_name` - The name of the private host which the server is assigned.
* `tags` - Any tags assigned to the Server.

---

A `network_interface` block exports the following:

* `mac_address` - The MAC address.
* `packet_filter_id` - The id of the packet filter attached to the network interface.
* `upstream` - The upstream type or upstream switch id. This will be one of [`shared`/`disconnect`/`<switch id>`].
* `user_ip_address` - The IP address for only display. This value doesn't affect actual NIC settings.



