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
  source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
  # or
  # source_disk_id  = "<your-disk-id>"
  
  # For "Modify Disk" API
  hostname          = "your-host-name"
  password          = ""
  ssh_key_ids       = ["<your-ssh-key-id>"]
  note_ids          = ["<your-note-id"]
  disable_pw_auth   = true
  
  description       = "description"
  tags              = ["foo", "bar"]
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
Valid value is one of the following: [ "ssd"(default) / "hdd"]
* `conector` - The disk connector of the resource.  
Valid value is one of the following: [ "virtio"(default) / "ide"]
* `size` - Size of the resource(unit:`GB`).
* `source_archive_id` - The ID of source Archive.
* `source_disk_id` - The ID of source Disk.
* `hostname` - The hostname to set with `"Modify Disk"` API.
* `password` - The password of OS's administrator to set with `"Modify Disk"` API.
* `ssh_key_ids` - The ID list of SSH Keys to set with `"Modify Disk"` API.
* `note_ids` - The ID list of Notes(Startup-Scripts) to set with `"Modify Disk"` API.
* `disable_pw_auth` - The flag of disable password auth via SSH.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time(seconds) to do graceful shutdown the server connected to the resource.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `plan` - The plan of the resource(`ssd`/`hdd`).
* `conector` - The disk connector of the resource(`virtio`/`ide`).
* `size` - Size of the resource(unit:`GB`).
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
