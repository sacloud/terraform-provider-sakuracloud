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

* `base64content` - (Optional) The base64 encoded content to upload to as the Icon. This conflicts with [`source`]. Changing this forces a new resource to be created.
* `name` - (Required) The name of the Icon. The length of this value must be in the range [`1`-`64`].
* `source` - (Optional) The file path to upload to as the Icon. This conflicts with [`base64content`]. Changing this forces a new resource to be created.
* `tags` - (Optional) Any tags to assign to the Icon.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Icon


* `update` - (Defaults to 5 minutes) Used when updating the Icon

* `delete` - (Defaults to 5 minutes) Used when deregistering Icon



## Attribute Reference

* `id` - The id of the Icon.
* `url` - The URL for getting the icon's raw data.




