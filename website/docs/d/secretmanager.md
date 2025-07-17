---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_secretmanager"
subcategory: "Global"
description: |-
  Get information about an existing SecretManager vault.
---

# Data Source: sakuracloud_secretmanager

Get information about an existing SecretManager vault.

## Example Usage

```hcl
data "sakuracloud_secretmanager" "foobar" {
  name = "foobar"
}
```

## Argument Reference

One of the following arguments must be specified:

* `name` - (Optional) The name of the SecretManager vault.
* `resource_id` - (Optional) The resource id of the SecretManager vault.

## Attribute Reference

* `id` - The id of the SecretManager vault.
* `name` - The name of the SecretManager vault.
* `kms_key_id` - KMS key id for SecretManager vault.
* `description` - The description of the SecretManager vault.
* `tags` - The tags attached to the SecretManager vault.
