---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_disk"
subcategory: "Storage"
description: |-
  Manages a SakuraCloud Disk.
---

# sakuracloud_disk

Manages a SakuraCloud Disk.

## Example Usage

```hcl
data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu2004"
}

resource "sakuracloud_disk" "foobar" {
  name              = "foobar"
  plan              = "ssd"
  connector         = "virtio"
  size              = 20
  source_archive_id = data.sakuracloud_archive.ubuntu.id
  #distant_from      = ["111111111111"]

  description = "description"
  tags        = ["tag1", "tag2"]
}
```

## Argument Reference

* `name` - (Required) The name of the disk. The length of this value must be in the range [`1`-`64`].

#### Disk Spec

* `connector` - (Optional) The name of the disk connector. This must be one of [`virtio`/`ide`]. Changing this forces a new resource to be created. Default:`virtio`.
* `plan` - (Optional) The plan name of the disk. This must be one of [`ssd`/`hdd`]. Changing this forces a new resource to be created. Default:`ssd`.
* `size` - (Optional) The size of disk in GiB. Changing this forces a new resource to be created. Default:`20`.
* `distant_from` - (Optional) A list of disk id. The disk will be located to different storage from these disks. Changing this forces a new resource to be created.

The values that can be specified for `size` can be found with the following command.

```bash
# for SSD Plan
usacloud disk-plan read 4 --zone is1a -o json | jq '.[].Size|.[]|select(.Availability == "available")|.SizeMB / 1024'

# for HDD Plan
usacloud disk-plan read 2 --zone is1a -o json | jq '.[].Size|.[]|select(.Availability == "available")|.SizeMB / 1024'
```

#### Disk Source

* `source_archive_id` - (Optional) The id of the source archive. This conflicts with [`source_disk_id`]. Changing this forces a new resource to be created.
* `source_disk_id` - (Optional) The id of the source disk. This conflicts with [`source_archive_id`]. Changing this forces a new resource to be created.

#### Common Arguments

* `description` - (Optional) The description of the disk. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the disk.
* `tags` - (Optional) Any tags to assign to the disk.
* `zone` - (Optional) The name of zone that the disk will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the Disk
* `update` - (Defaults to 24 hours) Used when updating the Disk
* `delete` - (Defaults to 20 minutes) Used when deleting Disk

## Attribute Reference

* `id` - The id of the Disk.
* `server_id` - The id of the Server connected to the disk.

