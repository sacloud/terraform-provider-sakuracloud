---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_private_host"
subcategory: "Compute"
description: |-
  Get information about an existing PrivateHost.
---

# Data Source: sakuracloud_private_host

Get information about an existing PrivateHost.

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

* `id` - The ID of the PrivateHost.
* `assigned_core` - .
* `assigned_memory` - .
* `class` - .
* `description` - .
* `hostname` - .
* `icon_id` - .
* `name` - .
* `tags` - .




