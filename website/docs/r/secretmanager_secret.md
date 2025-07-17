---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_kms"
subcategory: "Global"
description: |-
  Manages a SakuraCloud SecretManager secret.
---

# sakuracloud_kms

Manages a SakuraCloud SecretManager secret.

## Example Usage

```hcl
resource "sakuracloud_secretmanager_secret" "foobar" {
  name     = "foobar"
  vault_id = "secretmanager-resource-id" # e.g. sakuracloud_secretmanager.foobar.id
  value    = "Secret value!"
}
```

## Argument Reference

* `name` - (Required) The name of the SecretManager secret.
* `vault_id` - (Required) The resource id of the SecretManager vault.
* `value` - (Required) Secret value

## Attribute Reference

* `id` - The id of the SecretManager secret. This is same as `name`.
* `version` - The version of stored secret.
