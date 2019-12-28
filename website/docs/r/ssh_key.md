---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_ssh_key"
subcategory: "Misc"
description: |-
  Manages a SakuraCloud SSH Key.
---

# sakuracloud_ssh_key

Manages a SakuraCloud SSH Key.

## Argument Reference

* `description` - (Optional) .
* `name` - (Required) .
* `public_key` - (Required) . Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the SSH Key

* `read` -   (Defaults to 5 minutes) Used when reading the SSH Key

* `update` - (Defaults to 5 minutes) Used when updating the SSH Key

* `delete` - (Defaults to 5 minutes) Used when deregistering SSH Key



## Attribute Reference

* `id` - The id of the SSH Key.
* `fingerprint` - .




