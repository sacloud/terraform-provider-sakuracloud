---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_icon"
sidebar_current: "docs-sakuracloud-datasource-icon"
description: |-
  Get information on a SakuraCloud Icon.
---

# sakuracloud\_icon

Use this data source to retrieve information about a SakuraCloud Icon.

## Example Usage

```hcl
data "sakuracloud_icon" "foobar" {
  name_selectors = ["foobar"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The list of names to filtering.
 * `tag_selectors` - (Optional) The list of tags to filtering.
 * `filter` - (Optional) The map of filter key and value.

## Attributes Reference

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `body` - Base64 encoded icon data (size:`small`).
* `url` - URL to access this resource.
* `tags` - The tag list of the resources.
