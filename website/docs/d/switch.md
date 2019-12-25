---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_switch"
subcategory: "Networking"
description: |-
  Get information about an existing Switch.
---

# Data Source: sakuracloud_switch

Get information about an existing Switch.

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

* `id` - The ID of the Switch.
* `bridge_id` - .
* `description` - .
* `icon_id` - .
* `name` - .
* `server_ids` - .
* `tags` - .




