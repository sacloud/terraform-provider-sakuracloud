---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_load_balancer"
subcategory: "Appliance"
description: |-
  Get information about an existing LoadBalancer.
---

# Data Source: sakuracloud_load_balancer

Get information about an existing LoadBalancer.

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

* `id` - The ID of the LoadBalancer.
* `description` - .
* `gateway` - .
* `high_availability` - .
* `icon_id` - .
* `ip_addresses` - .
* `is_double` - .
* `name` - .
* `netmask` - .
* `plan` - .
* `switch_id` - .
* `tags` - .
* `vip` - A list of `vip` blocks as defined below.
* `vrid` - .


---

A `vip` block exports the following:

* `delay_loop` - .
* `description` - .
* `port` - .
* `server` - A list of `server` blocks as defined below.
* `sorry_server` - .
* `vip` - .

---

A `server` block exports the following:

* `check_path` - .
* `check_protocol` - .
* `check_status` - .
* `enabled` - .
* `ip_address` - .



