---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_subnet"
subcategory: "Networking"
description: |-
  Get information about an existing Subnet.
---

# Data Source: sakuracloud_subnet

Get information about an existing Subnet.

## Example Usage

```hcl
variable internet_id {}
data sakuracloud_subnet "foobar" {
  internet_id = var.internet_id
  index       = 1
}
```
## Argument Reference

* `index` - (Required) The index of the subnet in assigned to the Switch+Router. Changing this forces a new resource to be created.
* `internet_id` - (Required) The id of the switch+router resource that the Subnet belongs. Changing this forces a new resource to be created.



## Attribute Reference

* `id` - The id of the Subnet.
* `ip_addresses` - A list of assigned global address to the Subnet.
* `max_ip_address` - Maximum IP address in assigned global addresses to the Subnet.
* `min_ip_address` - Minimum IP address in assigned global addresses to the Subnet.
* `netmask` - The bit length of the subnet assigned to the Subnet.
* `network_address` - The IPv4 network address assigned to the Subnet.
* `next_hop` - The ip address of the next-hop at the Subnet.
* `switch_id` - The id of the switch connected from the Subnet.
* `zone` - The name of zone that the Subnet is in (e.g. `is1a`, `tk1a`).




