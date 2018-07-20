---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_icon"
sidebar_current: "docs-sakuracloud-datasource-icon"
description: |-
  Get information on a SakuraCloud icon.
---

# sakuracloud\_icon

Use this data source to retrieve information about a SakuraCloud icon.

## Example Usage

```hcl
data sakuracloud_icon "foobar" {
  name_selectors = ["foobar"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The filter value list of name.
 * `tag_selectors` - (Optional) The filter value list of tags.
 * `filter` - (Optional) The map of filter key and value.

## Attributes Reference

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `body` - Base64 encoded icon data(size:`small`).
* `url` - URL to access this resource.
* `tags` - The tag list of the resource.
