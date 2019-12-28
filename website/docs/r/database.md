---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_database"
subcategory: "Appliance"
description: |-
  Manages a SakuraCloud Database.
---

# sakuracloud_database

Manages a SakuraCloud Database.

## Argument Reference

* `backup_time` - (Optional) .
* `backup_weekdays` - (Optional) .
* `database_type` - (Optional) . Changing this forces a new resource to be created. Defaults to `postgres`.
* `description` - (Optional) .
* `gateway` - (Required) . Changing this forces a new resource to be created.
* `icon_id` - (Optional) .
* `ip_addresses` - (Required) . Changing this forces a new resource to be created.
* `name` - (Required) .
* `netmask` - (Required) . Changing this forces a new resource to be created.
* `password` - (Required) .
* `plan` - (Optional) . Changing this forces a new resource to be created. Defaults to `10g`.
* `port` - (Optional) . Defaults to `5432`.
* `replica_password` - (Optional) .
* `replica_user` - (Optional) . Defaults to `replica`.
* `source_ranges` - (Optional) .
* `switch_id` - (Required) . Changing this forces a new resource to be created.
* `tags` - (Optional) .
* `username` - (Required) . Changing this forces a new resource to be created.
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Database

* `read` -   (Defaults to 5 minutes) Used when reading the Database

* `update` - (Defaults to 60 minutes) Used when updating the Database

* `delete` - (Defaults to 20 minutes) Used when deregistering Database



## Attribute Reference

* `id` - The id of the Database.




