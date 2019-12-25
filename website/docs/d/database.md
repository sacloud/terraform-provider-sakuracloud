---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_database"
subcategory: "Appliance"
description: |-
  Get information about an existing Database.
---

# Data Source: sakuracloud_database

Get information about an existing Database.

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

* `id` - The ID of the Database.
* `backup_time` - .
* `backup_weekdays` - .
* `database_type` - .
* `description` - .
* `gateway` - .
* `icon_id` - .
* `ip_addresses` - .
* `name` - .
* `netmask` - .
* `password` - .
* `plan` - .
* `port` - .
* `replica_password` - .
* `replica_user` - .
* `source_ranges` - .
* `switch_id` - .
* `tags` - .
* `username` - .




