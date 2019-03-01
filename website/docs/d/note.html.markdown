---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_note"
sidebar_current: "docs-sakuracloud-datasource-note"
description: |-
  Get information on a SakuraCloud Note.
---

# sakuracloud\_note

Use this data source to retrieve information about a SakuraCloud Note.

## Example Usage

```hcl
data "sakuracloud_note" "foobar" {
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
* `class` - The name of the note class.
* `content` - The body of the note. 
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
