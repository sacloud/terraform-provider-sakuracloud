---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_ssh_key"
subcategory: "Misc"
description: |-
  Get information about an existing SSHKey.
---

# Data Source: sakuracloud_ssh_key

Get information about an existing SSHKey.

## Argument Reference

* `filter` - (Optional) A `filter` block as defined below.


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

* `id` - The ID of the SSHKey.
* `description` - .
* `fingerprint` - .
* `name` - .
* `public_key` - .




