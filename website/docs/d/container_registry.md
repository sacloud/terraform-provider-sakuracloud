---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_container_registry"
subcategory: "Lab"
description: |-
  Get information about an existing Container Registry.
---

# Data Source: sakuracloud_container_registry

Get information about an existing Container Registry.

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

* `id` - The ID of the Container Registry.
* `access_level` - .
* `description` - .
* `fqdn` - .
* `icon_id` - .
* `name` - .
* `subdomain_label` - .
* `tags` - .




