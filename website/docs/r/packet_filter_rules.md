---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_packet_filter_rules"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Packet Filter Rules.
---

# sakuracloud_packet_filter_rules

Manages a SakuraCloud Packet Filter Rules.

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

* `create` - (Defaults to 5 minutes) Used when creating the Packet Filter Rules

* `read` -   (Defaults to 5 minutes) Used when reading the Packet Filter Rules


* `delete` - (Defaults to 5 minutes) Used when deregistering Packet Filter Rules



## Attribute Reference

* `id` - The id of the Packet Filter Rules.




