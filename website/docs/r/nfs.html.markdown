---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_nfs"
sidebar_current: "docs-sakuracloud-resource-appliance-nfs"
description: |-
  Provides a SakuraCloud NFS Appliance resource. This can be used to create, update, and delete NFS Appliances.
---

# sakuracloud\_nfs

Provides a SakuraCloud NFS Appliance resource. This can be used to create, update, and delete NFS Appliances.

## Example Usage

```hcl
# Create a new NFS Appliance
resource "sakuracloud_nfs" "foobar" {
  name = "foobar"
  plan = 100

  switch_id     = sakuracloud_switch.foobar.id
  ipaddress     = "192.168.11.101"
  nw_mask_len   = 24
  default_route = "192.168.11.1"

  description = "description"
  tags        = ["foo", "bar"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `plan` - (Optional) The plan (size) of the NFS Appliance (unit:`GB`).  
Valid value is one of the following: [ 100 (default) / 500 / 1024 / 2048 / 4096 ]
* `switch_id` - (Required) The ID of the switch connected to the NFS Appliance.
* `ipaddress` - (Required) The IP address of the NFS Appliance.
* `nw_mask_len` - (Required) The network mask length of the NFS Appliance.
* `default_route` - (Required) The default route IP address of the NFS Appliance.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the NFS Appliance.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `plan` - The plan (size) of the NFS Appliance (unit:`GB`).
* `switch_id` - The ID of the switch connected to the NFS Appliance.
* `ipaddress` - The IP address of the NFS Appliance.
* `nw_mask_len` - The network mask length of the NFS Appliance.
* `default_route` - The default route IP address of the NFS Appliance.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import

NFS Appliances can be imported using the NFS Appliance ID.

```
$ terraform import sakuracloud_nfs.foobar <nfs_id>
```
