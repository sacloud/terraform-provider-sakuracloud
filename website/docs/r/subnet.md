---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_subnet"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Subnet.
---

# sakuracloud_subnet

Manages a SakuraCloud Subnet.

## Argument Reference

* `internet_id` - (Required) . Changing this forces a new resource to be created.
* `netmask` - (Optional) . Changing this forces a new resource to be created. Defaults to `28`.
* `next_hop` - (Required) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Subnet

* `read` -   (Defaults to 5 minutes) Used when reading the Subnet

* `update` - (Defaults to 60 minutes) Used when updating the Subnet

* `delete` - (Defaults to 5 minutes) Used when deregistering Subnet



## Attribute Reference

* `id` - The ID of the Subnet.
* `ip_addresses` - .
* `max_ip_address` - .
* `min_ip_address` - .
* `network_address` - .
* `switch_id` - .




