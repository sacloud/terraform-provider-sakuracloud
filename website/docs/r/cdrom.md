---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_cdrom"
subcategory: "Storage"
description: |-
  Manages a SakuraCloud CD-ROM.
---

# sakuracloud_cdrom

Manages a SakuraCloud CD-ROM.

## Example Usage

```hcl
resource "sakuracloud_cdrom" "foobar" {
  name           = "foobar"
  size           = 5
  iso_image_file = "example.iso"
  description    = "description"
  tags           = ["tag1", "tag2"]
}
```

## Argument Reference

* `name` - (Required) The name of the CD-ROM. The length of this value must be in the range [`1`-`64`].
* `content` - (Optional) The content to upload to as the CD-ROM. This conflicts with [`iso_image_file`].
* `content_file_name` - (Optional) The name of content file to upload to as the CD-ROM. This is only used when `content` is specified. This conflicts with [`iso_image_file`]. Default:`config`.
* `iso_image_file` - (Optional) The file path to upload to as the CD-ROM. This conflicts with [`content`].
* `hash` - (Optional) The md5 checksum calculated from the base64 encoded file body.
* `size` - (Optional) The size of CD-ROM in GiB. This must be one of [`5`/`10`/`20`]. Changing this forces a new resource to be created. Default:`5`.

#### Common Arguments

* `description` - (Optional) The description of the CD-ROM. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the CD-ROM.
* `tags` - (Optional) Any tags to assign to the CD-ROM.
* `zone` - (Optional) The name of zone that the CD-ROM will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the CD-ROM
* `update` - (Defaults to 24 hours) Used when updating the CD-ROM
* `delete` - (Defaults to 20 minutes) Used when deleting CD-ROM

## Attribute Reference

* `id` - The id of the CD-ROM.

