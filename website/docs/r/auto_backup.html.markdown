---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_auto_backup"
sidebar_current: "docs-sakuracloud-resource-appliance-auto-backup"
description: |-
  Provides a SakuraCloud Auto Backup resource. This can be used to create, update, and delete Auto Backups.
---

# sakuracloud\_auto_backup

Provides a SakuraCloud Auto Backup resource. This can be used to create, update, and delete Auto Backups.

## Example Usage

```hcl
# Create a new Auto Backup
resource "sakuracloud_auto_backup" "foobar" {
  name           = "foobar"
  disk_id        = sakuracloud_disk.disk.id
  weekdays       = ["fri", "sun"]
  max_backup_num = 2
  description    = "description"
  tags           = ["foo", "bar"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `disk_id` - (Optional) The ID of the target disk. 
* `weekdays` - (Optional) Day of the week to get backup.  
Valid values are the following: ["mon", "tue", "wed", "thu", "fri", "sat", "sun"]  
* `max_backup_num` - (Optional) Max number of backups to keep.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `zone` - (Optional) The ID of the zone to which the resource belongs.  
Valid value is one of the following: ["is1b" / "tk1a" / "is1a"]

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `disk_id` - The ID of the target disk. 
* `weekdays` - Day of the week to get backup.  
* `max_backup_num` - Max number of backups to keep.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Auto Backups can be imported using the Auto Backup ID.

```
$ terraform import sakuracloud_auto_backup.foobar <auto_backup_id>
```
