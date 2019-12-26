---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_nfs"
subcategory: "Appliance"
description: |-
  Get information about an existing NFS.
---

# Data Source: sakuracloud_nfs

Get information about an existing NFS.

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

* `id` - The ID of the NFS.
* `description` - .
* `gateway` - .
* `icon_id` - .
* `ip_address` - .
* `name` - .
* `netmask` - .
* `plan` - .
* `size` - .
* `switch_id` - .
* `tags` - .




