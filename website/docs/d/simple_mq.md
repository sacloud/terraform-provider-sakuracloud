---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_simple_mq"
subcategory: "Global"
description: |-
  Get information about an existing SimpleMQ.
---

# Data Source: sakuracloud_simple_mq

Get information about an existing SimpleMQ.

## Example Usage

```hcl
data "sakuracloud_simple_mq" "foobar" {
  name = "foobar"
}
```

## Argument Reference

One of the following arguments must be specified:

* `name` - (Optional) The name of the queue.
* `tags` - (Optional) The resource tags on SakuraCloud used for filtering. If multiple values are specified, they combined as AND condition.

## Attribute Reference

* `id` - The id of the queue.
* `name` - The name of the queue.
* `visibility_timeout_seconds` - The visibility timeout seconds of the queue.
* `expire_seconds` - The expire seconds of the queue.
* `description` - The description of the queue.
* `tags` - The tags attached to the queue.
* `icon_id` - The icon id attached to the queue.

