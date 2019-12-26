---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_ipv4_ptr"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud IPv4 PTR.
---

# sakuracloud_ipv4_ptr

Manages a SakuraCloud IPv4 PTR.

## Argument Reference

* `hostname` - (Required) .
* `ip_address` - (Required) .
* `retry_interval` - (Optional) . Defaults to `10`.
* `retry_max` - (Optional) . Defaults to `30`.
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the IPv4 PTR

* `read` -   (Defaults to 5 minutes) Used when reading the IPv4 PTR

* `update` - (Defaults to 60 minutes) Used when updating the IPv4 PTR

* `delete` - (Defaults to 5 minutes) Used when deregistering IPv4 PTR



## Attribute Reference

* `id` - The ID of the IPv4 PTR.




