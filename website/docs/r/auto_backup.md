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

* `description` - (Optional) The description of the AutoBackup. The length of this value must be in the range [`1`-`512`].
* `disk_id` - (Required) The disk id to backed up. Changing this forces a new resource to be created.
* `icon_id` - (Optional) The icon id to attach to the AutoBackup.
* `max_backup_num` - (Optional) The number backup files to keep. This must be in the range [`1`-`10`]. Default:`1`.
* `name` - (Required) The name of the AutoBackup. The length of this value must be in the range [`1`-`64`].
* `tags` - (Optional) Any tags to assign to the AutoBackup.
* `weekdays` - (Required) A list of weekdays to backed up. The values in the list must be in [`sun`/`mon`/`tue`/`wed`/`thu`/`fri`/`sat`].
* `zone` - (Optional) The name of zone that the AutoBackup will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Auto Backup


* `update` - (Defaults to 5 minutes) Used when updating the Auto Backup

* `delete` - (Defaults to 5 minutes) Used when deregistering Auto Backup



## Attribute Reference

* `id` - The id of the Auto Backup.




