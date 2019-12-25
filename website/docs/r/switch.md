---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_switch"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Switch.
---

# sakuracloud_switch

Manages a SakuraCloud Switch.

## Argument Reference

* `bridge_id` - (Optional) .
* `description` - (Optional) .
* `icon_id` - (Optional) .
* `name` - (Required) .
* `tags` - (Optional) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Switch

* `read` -   (Defaults to 5 minutes) Used when reading the Switch

* `update` - (Defaults to 5 minutes) Used when updating the Switch

* `delete` - (Defaults to 20 minutes) Used when deregistering Switch



## Attribute Reference

* `id` - The ID of the Switch.
* `server_ids` - .




