---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_kms"
subcategory: "Global"
description: |-
  Manages a SakuraCloud SecretManager vault.
---

# sakuracloud_kms

Manages a SakuraCloud SecretManager vault.

## Example Usage

```hcl
resource "sakuracloud_secretmanager" "foobar" {
  name        = "foobar"
  kms_key_id  = "kms-resource-id" # e.g. sakuracloud_kms.foobar.id
  description = "description"
  tags        = ["tag1", "tag2"]
}
```

## Argument Reference

* `name` - (Required) The name of the SecretManager vault.
* `kms_key_id` - (Required) KMS key id for SecretManager vault.
* `description` - (Optional) The description of the SecretManager vault.
* `tags` - (Optional) The tags attached to the SecretManager vault.

## Attribute Reference

* `id` - The id of the SecretManager vault.
