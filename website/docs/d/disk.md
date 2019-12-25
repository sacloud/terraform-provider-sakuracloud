---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_disk"
subcategory: "Storage"
description: |-
  Get information about an existing Disk.
---

# Data Source: sakuracloud_disk

Get information about an existing Disk.

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

* `id` - The ID of the Disk.
* `connector` - .
* `description` - .
* `icon_id` - .
* `name` - .
* `plan` - .
* `server_id` - .
* `size` - .
* `source_archive_id` - .
* `source_disk_id` - .
* `tags` - .




