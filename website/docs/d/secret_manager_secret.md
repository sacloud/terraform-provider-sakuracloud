---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_secret_manager"
subcategory: "Global"
description: |-
  Get information about an existing SecretManager secret.
---

# Data Source: sakuracloud_secret_manager

Get information about an existing SecretManager secret.

## Example Usage

```hcl
data "sakuracloud_secret_manager_secret" "foobar" {
  name     = "foobar"
  vault_id = "secret_manager-resource-id" # e.g. sakuracloud_secret_manager.foobar.id
}
```

## Argument Reference

* `name` - (Required) The name of the SecretManager secret.
* `vault_id` - (Required) The resource id of the SecretManager vault.
* `version` - (Optional) Target version to unveil stored secret. Without this parameter, latest version is used. 

## Attribute Reference

* `id` - The id of the SecretManager secret. This is same as `name`.
* `value` - Unveiled result of stored secret.
