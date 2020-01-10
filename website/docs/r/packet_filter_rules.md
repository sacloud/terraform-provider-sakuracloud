---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_packet_filter_rules"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Packet Filter Rules.
---

# sakuracloud_packet_filter_rules

Manages a SakuraCloud Packet Filter Rules.

## Argument Reference

* `expression` - (Optional) One or more `expression` blocks as defined below. Changing this forces a new resource to be created.
* `packet_filter_id` - (Required) The id of the packet filter that set expressions to. Changing this forces a new resource to be created.
* `zone` - (Optional) The name of zone that the PacketFilter Rule will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.


---

A `expression` block supports the following:

* `allow` - (Optional) The flag to allow the packet through the filter.
* `description` - (Optional) The description of the expression.
* `destination_port` - (Optional) A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`).
* `protocol` - (Required) The protocol used for filtering. This must be one of [`http`/`https`/`tcp`/`udp`/`icmp`/`fragment`/`ip`].
* `source_network` - (Optional) A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`).
* `source_port` - (Optional) A source port number or port range used for filtering (e.g. `1024`, `1024-2048`).


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Packet Filter Rules



* `delete` - (Defaults to 5 minutes) Used when deregistering Packet Filter Rules



## Attribute Reference

* `id` - The id of the Packet Filter Rules.




