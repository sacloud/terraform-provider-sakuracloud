---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_archive"
subcategory: "Storage"
description: |-
  Get information about an existing Archive.
---

# Data Source: sakuracloud_archive

Get information about an existing Archive.

## Argument Reference

* `filter` - (Optional) A `filter` block as defined below.
* `os_type` - (Optional) .
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

* `id` - The ID of the Archive.
* `description` - .
* `icon_id` - .
* `name` - .
* `size` - .
* `tags` - .




