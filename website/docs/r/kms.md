---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_kms"
subcategory: "Global"
description: |-
  Manages a SakuraCloud KMS key.
---

# sakuracloud_kms

Manages a SakuraCloud KMS key.

## Example Usage

```hcl
resource "sakuracloud_kms" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]
}
```

## Argument Reference

* `name` - (Required) The name of the KMS key.
* `description` - (Optional) The description of the KMS key.
* `tags` - (Optional) The tags attached to the KMS key.
* `key_origin` - (Optional) Key origin of the KMS key. 'generated' or 'imported'. Default is 'generated'
* `plain_key` - (Optional) Plain key for imported KMS key. Required when 'key_origin' is 'imported'.

## Attribute Reference

* `id` - The id of the KMS key.
