---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_simple_mq"
subcategory: "Global"
description: |-
  Manages a SakuraCloud SimpleMQ.
---

# sakuracloud_simple_mq

Manages a SakuraCloud SimpleMQ.

!> **Warning:** Queue's API key cannot be created with Terraform. To get one, head over to [SakuraCloud Console](https://secure.sakura.ad.jp/cloud/) or use [API](https://manual.sakura.ad.jp/api/cloud/simplemq/sacloud/#operation/rotateAPIKey).

## Example Usage

```hcl
resource "sakuracloud_simple_mq" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]

  visibility_timeout_seconds = 30
  expire_seconds             = 345600
}
```

## Argument Reference

* `name` - (Required) The name of the queue.
* `visibility_timeout_seconds` - (Optional) The visibility timeout seconds of the queue.
* `expire_seconds` - (Optional) The expire seconds of the queue.
* `description` - (Optional) The description of the queue.
* `tags` - (Optional) The tags attached to the queue.
* `icon_id` - (Optional) The icon id attached to the queue.

## Attribute Reference

* `id` - The id of the queue.

