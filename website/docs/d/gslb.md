---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_gslb"
subcategory: "Global"
description: |-
  Get information about an existing GSLB.
---

# Data Source: sakuracloud_gslb

Get information about an existing GSLB.

## Argument Reference

* `filter` - (Optional) A `filter` block as defined below.


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

* `id` - The ID of the GSLB.
* `description` - .
* `fqdn` - .
* `health_check` - A list of `health_check` blocks as defined below.
* `icon_id` - .
* `name` - .
* `server` - A list of `server` blocks as defined below.
* `sorry_server` - .
* `tags` - .
* `weighted` - .


---

A `health_check` block exports the following:

* `delay_loop` - .
* `host_header` - .
* `path` - .
* `port` - .
* `protocol` - .
* `status` - .

---

A `server` block exports the following:

* `enabled` - .
* `ip_address` - .
* `weight` - .



