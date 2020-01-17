---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_disk"
subcategory: "Storage"
description: |-
  Get information about an existing Disk.
---

# Data Source: sakuracloud_disk

Get information about an existing Disk.

## Example Usage

```hcl
data "sakuracloud_disk" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```
## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.
* `zone` - (Optional) The name of zone that the Disk is in (e.g. `is1a`, `tk1a`).

---

A `filter` block supports the following:

* `condition` - (Optional) One or more name/values pairs used for filtering. There are several valid keys, for a full reference, check out finding section in the [SakuraCloud API reference](https://developer.sakura.ad.jp/cloud/api/1.1/).
* `id` - (Optional) The resource id on SakuraCloud used for filtering.
* `names` - (Optional) The resource names on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.
* `tags` - (Optional) The resource tags on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.

---

A `condition` block supports the following:

* `name` - (Required) The name of the target field. This value is case-sensitive.
* `values` - (Required) The values of the condition. If multiple values ​​are specified, they combined as AND condition.


## Attribute Reference

* `id` - The id of the Disk.
* `connector` - The name of the disk connector. This will be one of [`virtio`/`ide`].
* `description` - The description of the Disk.
* `icon_id` - The icon id attached to the Disk.
* `name` - The name of the Disk.
* `plan` - The plan name of the Disk. This will be one of [`ssd`/`hdd`].
* `server_id` - The id of the Server connected to the Disk.
* `size` - The size of Disk in GiB.
* `source_archive_id` - The id of the source archive.
* `source_disk_id` - The id of the source disk.
* `tags` - Any tags assigned to the Disk.
* `zone` - The name of zone that the Disk is in (e.g. `is1a`, `tk1a`).



