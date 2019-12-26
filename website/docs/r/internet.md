---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_internet"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Internet.
---

# sakuracloud_internet

Manages a SakuraCloud Internet.

## Argument Reference

* `band_width` - (Optional) . Defaults to `100`.
* `description` - (Optional) .
* `enable_ipv6` - (Optional) .
* `icon_id` - (Optional) .
* `name` - (Required) .
* `netmask` - (Optional) . Changing this forces a new resource to be created. Defaults to `28`.
* `tags` - (Optional) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Internet

* `read` -   (Defaults to 5 minutes) Used when reading the Internet

* `update` - (Defaults to 60 minutes) Used when updating the Internet

* `delete` - (Defaults to 20 minutes) Used when deregistering Internet



## Attribute Reference

* `id` - The ID of the Internet.
* `gateway` - .
* `ip_addresses` - .
* `ipv6_network_address` - .
* `ipv6_prefix` - .
* `ipv6_prefix_len` - .
* `max_ip_address` - .
* `min_ip_address` - .
* `network_address` - .
* `server_ids` - .
* `switch_id` - .




