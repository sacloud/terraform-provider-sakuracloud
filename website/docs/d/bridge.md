---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_bridge"
subcategory: "Networking"
description: |-
  Get information about an existing Bridge.
---

# Data Source: sakuracloud_bridge

Get information about an existing Bridge.

## Argument Reference

* `filter` - (Optional) A `filter` block as defined below.
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.


---

A `filter` block supports the following:

* `condition` - (Optional) One or more `condition` blocks as defined below.
* `id` - (Optional) .
* `names` - (Optional) .

---

A `condition` block supports the following:

* `name` - (Required) .
* `values` - (Required) .


## Attribute Reference

* `id` - The ID of the Bridge.
* `description` - .
* `name` - .




