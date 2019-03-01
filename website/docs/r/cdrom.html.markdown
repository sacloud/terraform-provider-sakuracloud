---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_cdrom"
sidebar_current: "docs-sakuracloud-resource-storage-cdrom"
description: |-
  Provides a SakuraCloud CDROM (ISO-Image) resource. This can be used to create, update, and delete CDROMs.
---

# sakuracloud\_cdrom

Provides a SakuraCloud CDROM (ISO-Image) resource. This can be used to create, update, and delete CDROMs.

## Example Usage

```hcl
# Create a new CDROM(ISO-Image)
resource "sakuracloud_cdrom" "foobar" {
  name           = "foobar"
  size           = 5
 
  iso_image_file = "your/cdrom/file.iso"
  hash           = md5(file("your/cdrom/file.iso"))
  
  # or
  # content      = "your-content" 
  
  # or
  # content_file_path = file("your/cdrom/content.json")
  
  description = "description"
  tags        = ["foo", "bar"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `size` - (Optional) The size of the resource (unit:`GB`).   
Valid value is one of the following: [ 5 (default) / 10 ]
* `iso_image_file` - (Optional) CDROM file to upload (format:`raw`).
* `hash` - (Optional) MD5 hash value of the CDROM file.
* `content` - (Optional) String of the value of the CDROM. 
* `content_file_path` - (Optional) The source file path of the CDROM. 
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

CDROMs can be imported using the CDROM ID.

```
$ terraform import sakuracloud_cdrom.foobar <cdrom_id>
```
