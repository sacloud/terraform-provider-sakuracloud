---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_ssh_key_gen"
sidebar_current: "docs-sakuracloud-resource-misc-skeygen"
description: |-
  Provides a SakuraCloud SSH Key Gen resource. This can be used to create and delete SSH Keys.
---

# sakuracloud\_ssh_key

Provides a SakuraCloud SSH Key resource. This can be used to create and delete SSH Keys.
The private and public keys is generated on the Sakura Cloud platform.

## Example Usage

```hcl
# Create a new SSH Key
resource "sakuracloud_ssh_key_gen" "foobar" {
  name        = "foobar"
  pass_phrase = "<your-pass-phrase>"
  description = "description"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `pass_phrase` - (Optional) The path phrase of keys. 
* `description` - (Optional) The description of the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `description` - The description of the resource.
* `private_key` - The body of the generated private key. 
* `public_key` - The body of the generated public key. 
* `finger_print` - The fingerprint of the generated public key.

## Import (not supported)

Import of SSH Key Gen is not supported.
