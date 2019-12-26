---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_note"
subcategory: "Misc"
description: |-
  Manages a SakuraCloud Note.
---

# sakuracloud_note

Manages a SakuraCloud Note.

## Argument Reference

* `class` - (Optional) . Defaults to `shell`.
* `content` - (Required) .
* `icon_id` - (Optional) .
* `name` - (Required) .
* `tags` - (Optional) .



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Note

* `read` -   (Defaults to 5 minutes) Used when reading the Note

* `update` - (Defaults to 5 minutes) Used when updating the Note

* `delete` - (Defaults to 5 minutes) Used when deregistering Note



## Attribute Reference

* `id` - The ID of the Note.
* `description` - .




