---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_cdrom"
subcategory: "Storage"
description: |-
  Manages a SakuraCloud CD-ROM.
---

# sakuracloud_cdrom

Manages a SakuraCloud CD-ROM.

## Argument Reference

* `content` - (Optional) .
* `content_file_name` - (Optional) . Defaults to `config`.
* `description` - (Optional) .
* `hash` - (Optional) .
* `icon_id` - (Optional) .
* `iso_image_file` - (Optional) .
* `name` - (Required) .
* `size` - (Optional) . Changing this forces a new resource to be created. Defaults to `5`.
* `tags` - (Optional) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the CD-ROM

* `read` -   (Defaults to 5 minutes) Used when reading the CD-ROM

* `update` - (Defaults to 24 hours) Used when updating the CD-ROM

* `delete` - (Defaults to 20 minutes) Used when deregistering CD-ROM



## Attribute Reference

* `id` - The ID of the CD-ROM.




