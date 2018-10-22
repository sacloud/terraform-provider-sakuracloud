---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_icon"
sidebar_current: "docs-sakuracloud-resource-misc-icon"
description: |-
  Provides a SakuraCloud Icon resource. This can be used to create, update, and delete Icons.
---

# sakuracloud\_icon

Provides a SakuraCloud Icon resource. This can be used to create, update, and delete Icons.

## Example Usage

```hcl
# Create a new Icon
resource "sakuracloud_icon" "foobar" {
  name = "foobar"

  source = "path/to/your/file"
  # or
  #base64content = "<base64-encoded-content-body>"
  
  tags = ["foo", "bar"]
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `source` - (Optional) The path of source content file.
* `base64content` - (Optional) The base64 encoded body of source content.
* `tags` - (Optional) The tag list of the resources.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `body` - Base64 encoded icon data (size:`small`).
* `url` - URL to access this resource.
* `tags` - The tag list of the resources.

## Import

Icons can be imported using the Icon ID.

```
$ terraform import sakuracloud_icon.foobar <icon_id>
```
