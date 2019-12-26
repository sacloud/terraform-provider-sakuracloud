---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_dns"
subcategory: "Global"
description: |-
  Get information about an existing DNS.
---

# Data Source: sakuracloud_dns

Get information about an existing DNS.

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

* `id` - The ID of the DNS.
* `description` - .
* `dns_servers` - .
* `icon_id` - .
* `record` - A list of `record` blocks as defined below.
* `tags` - .
* `zone` - .


---

A `record` block exports the following:

* `name` - .
* `port` - .
* `priority` - .
* `ttl` - .
* `type` - .
* `value` - .
* `weight` - .



