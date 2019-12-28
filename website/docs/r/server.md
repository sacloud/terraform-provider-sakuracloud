---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_server"
subcategory: "Compute"
description: |-
  Manages a SakuraCloud Server.
---

# sakuracloud_server

Manages a SakuraCloud Server.

## Argument Reference

* `cdrom_id` - (Optional) .
* `commitment` - (Optional) . Defaults to `standard`.
* `core` - (Optional) . Defaults to `1`.
* `description` - (Optional) .
* `disk_edit_parameter` - (Optional) A `disk_edit_parameter` block as defined below.
* `disks` - (Optional) .
* `force_shutdown` - (Optional) .
* `icon_id` - (Optional) .
* `interface_driver` - (Optional) . Defaults to `virtio`.
* `memory` - (Optional) . Defaults to `1`.
* `name` - (Required) .
* `network_interface` - (Optional) One or more `network_interface` blocks as defined below.
* `private_host_id` - (Optional) .
* `tags` - (Optional) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.


---

A `disk_edit_parameter` block supports the following:

* `change_partition_uuid` - (Optional) .
* `disable_pw_auth` - (Optional) .
* `enable_dhcp` - (Optional) .
* `gateway` - (Optional) .
* `hostname` - (Optional) .
* `ip_address` - (Optional) .
* `netmask` - (Optional) .
* `note_ids` - (Optional) .
* `password` - (Optional) .
* `ssh_key_ids` - (Optional) .

---

A `network_interface` block supports the following:

* `packet_filter_id` - (Optional) .
* `upstream` - (Required) Upstream Network Type: valid value is one of [shared/disconnect/<switch id>].


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Server

* `read` -   (Defaults to 5 minutes) Used when reading the Server

* `update` - (Defaults to 5 minutes) Used when updating the Server

* `delete` - (Defaults to 20 minutes) Used when deregistering Server



## Attribute Reference

* `id` - The id of the Server.
* `dns_servers` - .
* `gateway` - .
* `ip_address` - .
* `netmask` - .
* `network_address` - .
* `private_host_name` - .


---

A `network_interface` block exports the following:

* `mac_address` - .



