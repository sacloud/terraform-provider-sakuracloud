---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_packet_filter"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Packet Filter.
---

# sakuracloud_packet_filter

Manages a SakuraCloud Packet Filter.

## Argument Reference

* `description` - (Optional) .
* `expression` - (Optional) One or more `expression` blocks as defined below.
* `name` - (Required) .
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.


---

A `expression` block supports the following:

* `allow` - (Optional) .
* `description` - (Optional) .
* `destination_port` - (Optional) .
* `protocol` - (Required) .
* `source_network` - (Optional) .
* `source_port` - (Optional) .


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Packet Filter

* `read` -   (Defaults to 5 minutes) Used when reading the Packet Filter

* `update` - (Defaults to 5 minutes) Used when updating the Packet Filter

* `delete` - (Defaults to 20 minutes) Used when deregistering Packet Filter



## Attribute Reference

* `id` - The id of the Packet Filter.




