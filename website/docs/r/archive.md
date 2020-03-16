---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_archive"
subcategory: "Storage"
description: |-
  Manages a SakuraCloud Archive.
---

# sakuracloud_archive

Manages a SakuraCloud Archive.

## Example Usage

```hcl
# from archive/disk
resource "sakuracloud_archive" "from-archive-or-disk" {
  name         = "foobar"
  description  = "description"
  tags         = ["tag1", "tag2"]

  source_archive_id   = 123456789012
  source_archive_zone = "tk1a"
  # source_disk_id    = 123456789012
}

# from shared archive
resource "sakuracloud_archive" "from-shared-archive" {
  name         = "foobar"
  description  = "description"
  tags         = ["tag1", "tag2"]

  source_shared_key = "is1a:123456789012:xxx"
}


# from local file
resource "sakuracloud_archive" "foobar" {
  name         = "foobar"
  description  = "description"
  tags         = ["tag1", "tag2"]
  size         = 20
  archive_file = "test/dummy.raw"
}
```

## Argument Reference

* `name` - (Required) The name of the archive. The length of this value must be in the range [`1`-`64`].
* `archive_file` - (Optional) The file path to upload to the SakuraCloud.
* `description` - (Optional) The description of the archive. The length of this value must be in the range [`1`-`512`].
* `hash` - (Optional) The md5 checksum calculated from the base64 encoded file body. Changing this forces a new resource to be created.
* `size` - (Optional) The size of archive in GiB. This must be one of [`20`/`40`/`60`/`80`/`100`/`250`/`500`/`750`/`1024`]. Changing this forces a new resource to be created. Default:`20`.
* `source_archive_id` - (Optional) The id of the source archive. This conflicts with [`source_disk_id`]. Changing this forces a new resource to be created.
* `source_archive_zone` - (Optional) The share key of source shared archive. Changing this forces a new resource to be created.
* `source_disk_id` - (Optional) The id of the source disk. This conflicts with [`source_archive_id`]. Changing this forces a new resource to be created.
* `source_shared_key` - (Optional) The share key of source shared archive. Changing this forces a new resource to be created.

#### Common Arguments

* `icon_id` - (Optional) The icon id to attach to the archive.
* `tags` - (Optional) Any tags to assign to the archive.
* `zone` - (Optional) The name of zone that the archive will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the Archive
* `update` - (Defaults to 24 hours) Used when updating the Archive
* `delete` - (Defaults to 5 minutes) Used when deleting Archive

## Attribute Reference

* `id` - The id of the Archive.



