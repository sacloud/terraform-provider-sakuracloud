---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_ssh_key"
subcategory: "Misc"
description: |-
  Manages a SakuraCloud SSHKey.
---

# sakuracloud_ssh_key

Manages a SakuraCloud SSHKey.

## Argument Reference

* `description` - (Optional) .
* `name` - (Required) .
* `public_key` - (Required) . Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the SSHKey

* `read` -   (Defaults to 5 minutes) Used when reading the SSHKey

* `update` - (Defaults to 5 minutes) Used when updating the SSHKey

* `delete` - (Defaults to 5 minutes) Used when deregistering SSHKey



## Attribute Reference

* `id` - The ID of the SSHKey.
* `fingerprint` - .




