---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_archive"
sidebar_current: "docs-sakuracloud-resource-storage-archive"
description: |-
  Provides a SakuraCloud Archive resource. This can be used to create, update, and delete Archives.
---

# sakuracloud\_archive

Provides a SakuraCloud Archive resource. This can be used to create, update, and delete Archives.

## Example Usage

```hcl
# Create a new Archive
resource "sakuracloud_archive" "foobar" {
  name         = "foobar"
  size         = 20
  archive_file = "your/archive/file.raw"
  hash         = md5(file("your/archive/file.raw"))
  description  = "description"
  tags         = ["foo", "bar"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `size` - (Optional) The size of the resource (unit:`GB`).   
Valid value is one of the following: [ 20 (default) / 40 / 60 / 80 / 100 / 250 / 500 / 750 / 1024 ]
* `archive_file` - (Optional) Archive file to upload (format:`raw`).
* `hash` - (Optional) MD5 hash value of the archive file.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `size` - The size of the resource (unit:`GB`).
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Archives can be imported using the Archive ID.

```
$ terraform import sakuracloud_archive.foobar <archive_id>
```
