---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_kms"
subcategory: "Global"
description: |-
  Get information about an existing KMS key.
---

# Data Source: sakuracloud_kms

Get information about an existing KMS key.

## Example Usage

```hcl
data "sakuracloud_kms" "foobar" {
  name = "foobar"
}
```

## Argument Reference

One of the following arguments must be specified:

* `name` - (Optional) The name of the KMS key.
* `resource_id` - (Optional) The resource id of the KMS key.

## Attribute Reference

* `id` - The id of the KMS key.
* `name` - The name of the KMS key.
* `description` - The description of the KMS key.
* `tags` - The tags attached to the KMS key.
* `key_origin` - The key origin of the KMS key. "generated" or "imported".
