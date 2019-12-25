---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_internet"
subcategory: "Networking"
description: |-
  Get information about an existing Internet.
---

# Data Source: sakuracloud_internet

Get information about an existing Internet.

## Argument Reference

* `filter` - (Optional) A `filter` block as defined below.
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.


---

A `filter` block supports the following:

* `condition` - (Optional) One or more `condition` blocks as defined below.
* `id` - (Optional) .
* `names` - (Optional) .
* `tags` - (Optional) .

---

A `condition` block supports the following:

* `name` - (Required) .
* `values` - (Required) .


## Attribute Reference

* `id` - The ID of the Internet.
* `band_width` - .
* `description` - .
* `enable_ipv6` - .
* `gateway` - .
* `icon_id` - .
* `ip_addresses` - .
* `ipv6_network_address` - .
* `ipv6_prefix` - .
* `ipv6_prefix_len` - .
* `max_ip_address` - .
* `min_ip_address` - .
* `name` - .
* `netmask` - .
* `network_address` - .
* `server_ids` - .
* `switch_id` - .
* `tags` - .




