---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_subnet"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Subnet.
---

# sakuracloud_subnet

Manages a SakuraCloud Subnet.

## Example Usage

```hcl
resource sakuracloud_internet "foobar" {
  name = "foobar"
}

resource "sakuracloud_subnet" "foobar" {
  internet_id = sakuracloud_internet.foobar.id
  next_hop    = sakuracloud_internet.foobar.min_ip_address
}
```
## Argument Reference

* `internet_id` - (Required) The id of the switch+router resource that the subnet belongs. Changing this forces a new resource to be created.
* `netmask` - (Optional) The bit length of the subnet to assign to the Subnet. This must be in the range [`26`-`28`]. Changing this forces a new resource to be created. Default:`28`.
* `next_hop` - (Required) The ip address of the next-hop at the subnet.
* `zone` - (Optional) The name of zone that the Subnet will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Subnet


* `update` - (Defaults to 60 minutes) Used when updating the Subnet

* `delete` - (Defaults to 5 minutes) Used when deregistering Subnet



## Attribute Reference

* `id` - The id of the Subnet.
* `ip_addresses` - A list of assigned global address to the subnet.
* `max_ip_address` - Maximum IP address in assigned global addresses to the subnet.
* `min_ip_address` - Minimum IP address in assigned global addresses to the subnet.
* `network_address` - The IPv4 network address assigned to the Subnet.
* `switch_id` - The id of the switch connected from the Subnet.



