---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_subnet"
subcategory: "Networking"
description: |-
  Get information about an existing Subnet.
---

# Data Source: sakuracloud_subnet

Get information about an existing Subnet.

## Argument Reference

* `index` - (Required) The index of the subnet in assigned to the switch+router. Changing this forces a new resource to be created.
* `internet_id` - (Required) The id of the switch+router resource that the subnet belongs. Changing this forces a new resource to be created.



## Attribute Reference

* `id` - The id of the Subnet.
* `ip_addresses` - A list of assigned global address to the subnet.
* `max_ip_address` - Maximum IP address in assigned global addresses to the subnet.
* `min_ip_address` - Minimum IP address in assigned global addresses to the subnet.
* `netmask` - The bit length of the subnet assigned to the subnet.
* `network_address` - The IPv4 network address assigned to the subnet.
* `next_hop` - The ip address of the next-hop at the subnet.
* `switch_id` - The id of the switch connected from the subnet.
* `zone` - The name of zone that the subnet is in (e.g. `is1a`,`tk1a`).




