---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_switch"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Switch.
---

# sakuracloud_switch

Manages a SakuraCloud Switch.

## Example Usage

```hcl
resource "sakuracloud_switch" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]
}
```
## Argument Reference

* `bridge_id` - (Optional) The bridge id attached to the Switch.
* `description` - (Optional) The description of the Switch. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the Switch.
* `name` - (Required) The name of the Switch. The length of this value must be in the range [`1`-`64`].
* `tags` - (Optional) Any tags to assign to the Switch.
* `zone` - (Optional) The name of zone that the Switch will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Switch


* `update` - (Defaults to 5 minutes) Used when updating the Switch

* `delete` - (Defaults to 20 minutes) Used when deregistering Switch



## Attribute Reference

* `id` - The id of the Switch.
* `server_ids` - A list of server id connected to the switch.



