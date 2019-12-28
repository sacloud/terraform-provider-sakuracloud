---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_archive"
subcategory: "Storage"
description: |-
  Manages a SakuraCloud Archive.
---

# sakuracloud_archive

Manages a SakuraCloud Archive.

## Argument Reference

* `archive_file` - (Required) .
* `description` - (Optional) .
* `hash` - (Optional) .
* `icon_id` - (Optional) .
* `name` - (Required) .
* `size` - (Optional) . Changing this forces a new resource to be created. Defaults to `20`.
* `tags` - (Optional) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the Archive

* `read` -   (Defaults to 5 minutes) Used when reading the Archive

* `update` - (Defaults to 24 hours) Used when updating the Archive

* `delete` - (Defaults to 5 minutes) Used when deregistering Archive



## Attribute Reference

* `id` - The id of the Archive.




