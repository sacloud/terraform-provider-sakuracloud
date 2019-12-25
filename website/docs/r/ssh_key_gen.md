---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_ssh_key_gen"
subcategory: "Misc"
description: |-
  Manages a SakuraCloud SSHKey Gen.
---

# sakuracloud_ssh_key_gen

Manages a SakuraCloud SSHKey Gen.

## Argument Reference

* `description` - (Optional) . Changing this forces a new resource to be created.
* `name` - (Required) . Changing this forces a new resource to be created.
* `pass_phrase` - (Optional) . Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the SSHKey Gen

* `read` -   (Defaults to 5 minutes) Used when reading the SSHKey Gen

* `update` - (Defaults to 5 minutes) Used when updating the SSHKey Gen

* `delete` - (Defaults to 5 minutes) Used when deregistering SSHKey Gen



## Attribute Reference

* `id` - The ID of the SSHKey Gen.
* `fingerprint` - .
* `private_key` - .
* `public_key` - .




