---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_bridge"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Bridge.
---

# sakuracloud_bridge

Manages a SakuraCloud Bridge.

## Argument Reference

* `description` - (Optional) .
* `name` - (Required) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 minutes) Used when creating the Bridge

* `read` -   (Defaults to 5 minutes) Used when reading the Bridge

* `update` - (Defaults to 20 minutes) Used when updating the Bridge

* `delete` - (Defaults to 20 minutes) Used when deregistering Bridge



## Attribute Reference

* `id` - The ID of the Bridge.




