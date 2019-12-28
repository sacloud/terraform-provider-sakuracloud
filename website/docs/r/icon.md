---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_icon"
subcategory: "Misc"
description: |-
  Manages a SakuraCloud Icon.
---

# sakuracloud_icon

Manages a SakuraCloud Icon.

## Argument Reference

* `base64content` - (Optional) . Changing this forces a new resource to be created.
* `name` - (Required) .
* `source` - (Optional) . Changing this forces a new resource to be created.
* `tags` - (Optional) .



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Icon

* `read` -   (Defaults to 5 minutes) Used when reading the Icon

* `update` - (Defaults to 5 minutes) Used when updating the Icon

* `delete` - (Defaults to 5 minutes) Used when deregistering Icon



## Attribute Reference

* `id` - The id of the Icon.
* `url` - .




