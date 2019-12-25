---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_packet_filter"
subcategory: "Networking"
description: |-
  Get information about an existing PacketFilter.
---

# Data Source: sakuracloud_packet_filter

Get information about an existing PacketFilter.

## Argument Reference

* `expression` - (Optional) One or more `expression` blocks as defined below.
* `filter` - (Optional) A `filter` block as defined below.
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.


---

A `filter` block supports the following:

* `condition` - (Optional) One or more `condition` blocks as defined below.
* `id` - (Optional) .
* `names` - (Optional) .

---

A `condition` block supports the following:

* `name` - (Required) .
* `values` - (Required) .


## Attribute Reference

* `id` - The ID of the PacketFilter.
* `description` - .
* `name` - .


---

A `expression` block exports the following:

* `allow` - .
* `description` - .
* `destination_port` - .
* `protocol` - .
* `source_network` - .
* `source_port` - .



