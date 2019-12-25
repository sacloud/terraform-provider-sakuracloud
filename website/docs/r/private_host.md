---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_private_host"
subcategory: "Compute"
description: |-
  Manages a SakuraCloud PrivateHost.
---

# sakuracloud_private_host

Manages a SakuraCloud PrivateHost.

## Argument Reference

* `class` - (Optional) . Defaults to `dynamic`.
* `description` - (Optional) .
* `icon_id` - (Optional) .
* `name` - (Required) .
* `tags` - (Optional) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the PrivateHost

* `read` -   (Defaults to 5 minutes) Used when reading the PrivateHost

* `update` - (Defaults to 5 minutes) Used when updating the PrivateHost

* `delete` - (Defaults to 20 minutes) Used when deregistering PrivateHost



## Attribute Reference

* `id` - The ID of the PrivateHost.
* `assigned_core` - .
* `assigned_memory` - .
* `hostname` - .




