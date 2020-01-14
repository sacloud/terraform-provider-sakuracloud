---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_ssh_key_gen"
subcategory: "Misc"
description: |-
  Manages a SakuraCloud SSH Key Gen.
---

# sakuracloud_ssh_key_gen

Manages a SakuraCloud SSH Key Gen.

## Example Usage

```hcl
resource "sakuracloud_ssh_key_gen" "foobar" {
  name = "foobar"
  #pass_phrase = "your-pass-phrase"
  description = "description"
}
```
## Argument Reference

* `description` - (Optional) The description of the SSHKey. The length of this value must be in the range [`1`-`512`]. Changing this forces a new resource to be created.
* `name` - (Required) The name of the SSHKey. The length of this value must be in the range [`1`-`64`]. Changing this forces a new resource to be created.
* `pass_phrase` - (Optional) The pass phrase of the private key. The length of this value must be in the range [`8`-`64`]. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the SSH Key Gen


* `update` - (Defaults to 5 minutes) Used when updating the SSH Key Gen

* `delete` - (Defaults to 5 minutes) Used when deregistering SSH Key Gen



## Attribute Reference

* `id` - The id of the SSH Key Gen.
* `fingerprint` - The fingerprint of the public key.
* `private_key` - The body of the private key.
* `public_key` - The body of the public key.




