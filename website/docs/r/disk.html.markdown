---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_disk"
sidebar_current: "docs-sakuracloud-resource-storage-disk"
description: |-
  Provides a SakuraCloud Disk resource. This can be used to create, update, and delete Disks.
---

# sakuracloud\_disk

Provides a SakuraCloud Disk resource. This can be used to create, update, and delete Disks.

## Example Usage

```hcl
# Create a new Disk
resource "sakuracloud_disk" "foobar" {
  name              = "foobar"
  plan              = "ssd"
  connector         = "virtio"
  size              = 20
  
  source_archive_id = data.sakuracloud_archive.ubuntu.id
  # or
  #source_disk_id = "<your-disk-id>"

  # for storage isolation
  #distant_from = ["<your-disk-id>"]

  description = "description"
  tags        = ["foo", "bar"]
}

# Source archive
data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `plan` - The plan of the resource.  
Valid value is one of the following: [ "ssd" (default) / "hdd"]
* `conector` - (Optional) The disk connector of the resource.  
Valid value is one of the following: [ "virtio" (default) / "ide"]
* `size` - (Optional) Size of the resource (unit:`GB`).
* `distant_from` - (Optional) The ID list of the Disks isolated from Disk.
* `source_archive_id` - (Optional) The ID of source Archive.
* `source_disk_id` - (Optional) The ID of source Disk.
* `hostname` - (**Deprecated**) The hostname to set with `"Modify Disk"` API.
* `password` - (**Deprecated**) The password of OS's administrator to set with `"Modify Disk"` API.
* `ssh_key_ids` - (**Deprecated**) The ID list of SSH Keys to set with `"Modify Disk"` API.
* `note_ids` - (**Deprecated**) The ID list of Notes (Startup-Scripts) to set with `"Modify Disk"` API.
* `disable_pw_auth` - (**Deprecated**) The flag of disable password auth via SSH.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the server connected to the resource.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `plan` - The plan of the resource (`ssd`/`hdd`).
* `conector` - The disk connector of the resource (`virtio`/`ide`).
* `size` - Size of the resource (unit:`GB`).
* `server_id` - The ID of the server connected to the disk.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Disks can be imported using the Disk ID.

```
$ terraform import sakuracloud_disk.foobar <disk_id>
```
