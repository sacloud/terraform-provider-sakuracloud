---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_disk"
subcategory: "Storage"
description: |-
  Manages a SakuraCloud Disk.
---

# sakuracloud_disk

Manages a SakuraCloud Disk.

## Argument Reference

* `connector` - (Optional) . Changing this forces a new resource to be created. Defaults to `virtio`.
* `description` - (Optional) .
* `distant_from` - (Optional) . Changing this forces a new resource to be created.
* `icon_id` - (Optional) .
* `name` - (Required) .
* `plan` - (Optional) . Changing this forces a new resource to be created. Defaults to `ssd`.
* `size` - (Optional) . Changing this forces a new resource to be created. Defaults to `20`.
* `source_archive_id` - (Optional) . Changing this forces a new resource to be created.
* `source_disk_id` - (Optional) . Changing this forces a new resource to be created.
* `tags` - (Optional) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the Disk

* `read` -   (Defaults to 5 minutes) Used when reading the Disk

* `update` - (Defaults to 24 hours) Used when updating the Disk

* `delete` - (Defaults to 20 minutes) Used when deregistering Disk



## Attribute Reference

* `id` - The id of the Disk.
* `server_id` - .




