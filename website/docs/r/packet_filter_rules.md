---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_packet_filter_rules"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Packet Filter Rules.
---

# sakuracloud_packet_filter_rules

Manages a SakuraCloud Packet Filter Rules.

## Example Usage

```hcl
resource "sakuracloud_packet_filter" "foobar" {
  name        = "foobar"
  description = "description"
}

resource "sakuracloud_packet_filter_rules" "rules" {
  packet_filter_id = sakuracloud_packet_filter.foobar.id

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

* `packet_filter_id` - (Required) The id of the packet filter that set expressions to. Changing this forces a new resource to be created.
* `expression` - (Optional) One or more `expression` blocks as defined below. Changing this forces a new resource to be created.
* `zone` - (Optional) The name of zone that the PacketFilter Rule will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.

---

A `expression` block supports the following:

* `protocol` - (Required) The protocol used for filtering. This must be one of [`http`/`https`/`tcp`/`udp`/`icmp`/`fragment`/`ip`].
* `allow` - (Optional) The flag to allow the packet through the filter.
* `destination_port` - (Optional) A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`).
* `source_network` - (Optional) A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`).
* `source_port` - (Optional) A source port number or port range used for filtering (e.g. `1024`, `1024-2048`).
* `description` - (Optional) The description of the expression.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Packet Filter Rules
* `delete` - (Defaults to 5 minutes) Used when deleting Packet Filter Rules

## Attribute Reference

* `id` - The id of the Packet Filter Rules.

