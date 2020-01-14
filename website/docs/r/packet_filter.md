---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_packet_filter"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Packet Filter.
---

# sakuracloud_packet_filter

Manages a SakuraCloud Packet Filter.

## Example Usage

```hcl
resource "sakuracloud_packet_filter" "foobar" {
  name        = "foobar"
  description = "description"

  expression {
    protocol  = "tcp"
    dest_port = "22"
  }

  expression {
    protocol  = "tcp"
    dest_port = "80"
  }

  expression {
    protocol  = "tcp"
    dest_port = "443"
  }

  expression {
    protocol = "icmp"
  }

  expression {
    protocol = "fragment"
  }

  expression {
    protocol    = "udp"
    source_port = "123"
  }

  expression {
    protocol  = "tcp"
    dest_port = "32768-61000"
  }

  expression {
    protocol  = "udp"
    dest_port = "32768-61000"
  }

  expression {
    protocol    = "ip"
    allow       = false
    description = "Deny ALL"
  }
}
```
## Argument Reference

* `description` - (Optional) The description of the packetFilter. The length of this value must be in the range [`1`-`512`].
* `expression` - (Optional) One or more `expression` blocks as defined below.
* `name` - (Required) The name of the packetFilter. The length of this value must be in the range [`1`-`64`].
* `zone` - (Optional) The name of zone that the packetFilter will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.


---

A `expression` block supports the following:

* `allow` - (Optional) The flag to allow the packet through the filter.
* `description` - (Optional) The description of the expression.
* `destination_port` - (Optional) A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`).
* `protocol` - (Required) The protocol used for filtering. This must be one of [`http`/`https`/`tcp`/`udp`/`icmp`/`fragment`/`ip`].
* `source_network` - (Optional) A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`).
* `source_port` - (Optional) A source port number or port range used for filtering (e.g. `1024`, `1024-2048`).


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Packet Filter


* `update` - (Defaults to 5 minutes) Used when updating the Packet Filter

* `delete` - (Defaults to 20 minutes) Used when deregistering Packet Filter



## Attribute Reference

* `id` - The id of the Packet Filter.




