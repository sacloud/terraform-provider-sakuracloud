---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_packet_filter_rules"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud PacketFilter Rules.
---

# sakuracloud_packet_filter_rules

Manages a SakuraCloud PacketFilter Rules.

## Argument Reference

* `expression` - (Optional) One or more `expression` blocks as defined below. Changing this forces a new resource to be created.
* `packet_filter_id` - (Required) . Changing this forces a new resource to be created.
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

* `create` - (Defaults to 5 minutes) Used when creating the PacketFilter Rules

* `read` -   (Defaults to 5 minutes) Used when reading the PacketFilter Rules


* `delete` - (Defaults to 5 minutes) Used when deregistering PacketFilter Rules



## Attribute Reference

* `id` - The ID of the PacketFilter Rules.




