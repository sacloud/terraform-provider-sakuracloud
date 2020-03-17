---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_archive_share"
subcategory: "Storage"
description: |-
  Manages a SakuraCloud Archive Sharing.
---

# sakuracloud_archive_share

Manages a SakuraCloud Archive Sharing.

## Example Usage

```hcl
resource "sakuracloud_archive" "source" {
  name         = "foobar"
  size         = 20
  archive_file = "test/dummy.raw"
}

resource "sakuracloud_archive_share" "share_info" {
  archive_id = sakuracloud_archive.source.id
}
```
## Argument Reference

* `archive_id` - (Required) The id of the archive. Changing this forces a new resource to be created.
* `zone` - (Optional) The name of zone that the ArchiveShare will be created (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Archive
* `update` - (Defaults to 5 minutes) Used when updating the Archive
* `delete` - (Defaults to 5 minutes) Used when deleting Archive

## Attribute Reference

* `id` - The id of the Archive.
* `share_key` - The key to use sharing the Archive.

