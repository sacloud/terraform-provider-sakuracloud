---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_auto_backup"
subcategory: "Appliance"
description: |-
  Manages a SakuraCloud Auto Backup.
---

# sakuracloud_auto_backup

Manages a SakuraCloud Auto Backup.

## Argument Reference

* `description` - (Optional) .
* `disk_id` - (Required) . Changing this forces a new resource to be created.
* `icon_id` - (Optional) .
* `max_backup_num` - (Optional) . Defaults to `1`.
* `name` - (Required) . Changing this forces a new resource to be created.
* `tags` - (Optional) .
* `weekdays` - (Required) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Auto Backup

* `read` -   (Defaults to 5 minutes) Used when reading the Auto Backup

* `update` - (Defaults to 5 minutes) Used when updating the Auto Backup

* `delete` - (Defaults to 5 minutes) Used when deregistering Auto Backup



## Attribute Reference

* `id` - The id of the Auto Backup.




