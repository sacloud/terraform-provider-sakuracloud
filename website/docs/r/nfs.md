---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_nfs"
subcategory: "Appliance"
description: |-
  Manages a SakuraCloud NFS.
---

# sakuracloud_nfs

Manages a SakuraCloud NFS.

## Example Usage

```hcl
resource "sakuracloud_nfs" "foobar" {
  name        = "foobar"
  switch_id   = sakuracloud_switch.foobar.id
  plan        = "ssd"
  size        = "500"
  ip_address  = "192.168.11.101"
  netmask     = 24
  gateway     = "192.168.11.1"
  description = "description"
  tags        = ["tag1", "tag2"]
}

resource "sakuracloud_switch" "foobar" {
  name = "foobar"
}
```

## Argument Reference

* `name` - (Required) The name of the NFS. The length of this value must be in the range [`1`-`64`].
* `plan` - (Optional) The plan name of the NFS. This must be one of [`hdd`/`ssd`]. Changing this forces a new resource to be created. Default:`hdd`.
* `size` - (Optional) The size of NFS in GiB. Changing this forces a new resource to be created. Default:`100`.

#### Network

* `switch_id` - (Required) The id of the switch to which the NFS connects. Changing this forces a new resource to be created.
* `ip_address` - (Required) The IP address to assign to the NFS. Changing this forces a new resource to be created.
* `netmask` - (Required) The bit length of the subnet to assign to the NFS. This must be in the range [`8`-`29`]. Changing this forces a new resource to be created.
* `gateway` - (Optional) The IP address of the gateway used by NFS. Changing this forces a new resource to be created.

#### Common Arguments

* `description` - (Optional) The description of the NFS. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the NFS.
* `tags` - (Optional) Any tags to assign to the NFS.
* `zone` - (Optional) The name of zone that the NFS will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the NFS
* `update` - (Defaults to 24 hours) Used when updating the NFS
* `delete` - (Defaults to 20 minutes) Used when deleting NFS

## Attribute Reference

* `id` - The id of the NFS.

