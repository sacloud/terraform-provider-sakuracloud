---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_server"
subcategory: "Compute"
description: |-
  Get information about an existing Server.
---

# Data Source: sakuracloud_server

Get information about an existing Server.

## Argument Reference

* `filter` - (Optional) A `filter` block as defined below.
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.


---

A `filter` block supports the following:

* `condition` - (Optional) One or more `condition` blocks as defined below.
* `id` - (Optional) .
* `names` - (Optional) .
* `tags` - (Optional) .

---

A `condition` block supports the following:

* `name` - (Required) .
* `values` - (Required) .


## Attribute Reference

* `id` - The ID of the Server.
* `cdrom_id` - .
* `commitment` - .
* `core` - .
* `description` - .
* `disks` - .
* `dns_servers` - .
* `gateway` - .
* `icon_id` - .
* `interface_driver` - .
* `ip_address` - .
* `memory` - .
* `name` - .
* `netmask` - .
* `network_address` - .
* `network_interface` - A list of `network_interface` blocks as defined below.
* `private_host_id` - .
* `private_host_name` - .
* `tags` - .


---

A `network_interface` block exports the following:

* `mac_address` - .
* `packet_filter_id` - .
* `upstream` - .



