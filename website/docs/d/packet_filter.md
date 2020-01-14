---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_packet_filter"
subcategory: "Networking"
description: |-
  Get information about an existing Packet Filter.
---

# Data Source: sakuracloud_packet_filter

Get information about an existing Packet Filter.

## Example Usage

```hcl
data "sakuracloud_packet_filter" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```
## Argument Reference

* `expression` - (Optional) One or more `expression` blocks as defined below.
* `filter` - (Optional) One or more values used for filtering, as defined below.


---

A `filter` block supports the following:

* `condition` - (Optional) One or more name/values pairs used for filtering. There are several valid keys, for a full reference, check out finding section in the [SakuraCloud API reference](https://developer.sakura.ad.jp/cloud/api/1.1/).
* `id` - (Optional) The resource id on SakuraCloud used for filtering.
* `names` - (Optional) The resource names on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.

---

A `condition` block supports the following:

* `name` - (Required) The name of the target field. This value is case-sensitive.
* `values` - (Required) The values of the condition. If multiple values ​​are specified, they combined as AND condition.


## Attribute Reference

* `id` - The id of the Packet Filter.
* `description` - The description of the PacketFilter.
* `name` - The name of the PacketFilter.
* `zone` - The name of zone that the PacketFilter is in (e.g. `is1a`, `tk1a`).


---

A `expression` block exports the following:

* `allow` - The flag to allow the packet through the filter.
* `description` - The description of the expression.
* `destination_port` - A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`).
* `protocol` - The protocol used for filtering. This will be one of [`http`/`https`/`tcp`/`udp`/`icmp`/`fragment`/`ip`].
* `source_network` - A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`).
* `source_port` - A source port number or port range used for filtering (e.g. `1024`, `1024-2048`).



