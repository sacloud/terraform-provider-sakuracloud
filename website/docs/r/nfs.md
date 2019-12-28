---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_nfs"
subcategory: "Appliance"
description: |-
  Manages a SakuraCloud NFS.
---

# sakuracloud_nfs

Manages a SakuraCloud NFS.

## Argument Reference

* `description` - (Optional) .
* `gateway` - (Optional) . Changing this forces a new resource to be created.
* `icon_id` - (Optional) .
* `ip_address` - (Required) . Changing this forces a new resource to be created.
* `name` - (Required) .
* `netmask` - (Required) . Changing this forces a new resource to be created.
* `plan` - (Optional) . Changing this forces a new resource to be created. Defaults to `hdd`.
* `size` - (Optional) . Changing this forces a new resource to be created. Defaults to `100`.
* `switch_id` - (Required) . Changing this forces a new resource to be created.
* `tags` - (Optional) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the NFS

* `read` -   (Defaults to 5 minutes) Used when reading the NFS

* `update` - (Defaults to 24 hours) Used when updating the NFS

* `delete` - (Defaults to 20 minutes) Used when deregistering NFS



## Attribute Reference

* `id` - The id of the NFS.




