---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_database_read_replica"
subcategory: "Appliance"
description: |-
  Manages a SakuraCloud Database Read Replica.
---

# sakuracloud_database_read_replica

Manages a SakuraCloud Database Read Replica.

## Argument Reference

* `description` - (Optional) .
* `gateway` - (Optional) . Changing this forces a new resource to be created.
* `icon_id` - (Optional) .
* `ip_addresses` - (Required) . Changing this forces a new resource to be created.
* `master_id` - (Required) . Changing this forces a new resource to be created.
* `name` - (Required) .
* `netmask` - (Optional) . Changing this forces a new resource to be created.
* `source_ranges` - (Optional) .
* `switch_id` - (Optional) . Changing this forces a new resource to be created.
* `tags` - (Optional) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Database Read Replica

* `read` -   (Defaults to 5 minutes) Used when reading the Database Read Replica

* `update` - (Defaults to 60 minutes) Used when updating the Database Read Replica

* `delete` - (Defaults to 20 minutes) Used when deregistering Database Read Replica



## Attribute Reference

* `id` - The id of the Database Read Replica.




