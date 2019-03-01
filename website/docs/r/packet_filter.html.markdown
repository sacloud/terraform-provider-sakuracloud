---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_packet_filter"
sidebar_current: "docs-sakuracloud-resource-networking-packet-filter"
description: |-
  Provides a SakuraCloud Packet Filter resource. This can be used to create, update, and delete Packet Filters.
---

# sakuracloud\_packet\_filter

Provides a SakuraCloud Packet Filter resource. This can be used to create, update, and delete Packet Filters.

## Example Usage

```hcl
# Create a new Packet Filter
resource "sakuracloud_packet_filter" "foobar" {
  name        = "foobar"
  description = "description"

  expressions {
    protocol    = "tcp"
    source_nw   = "0.0.0.0/0"
    source_port = "0-65535"
    dest_port   = "80"
  }

  expressions {
    protocol    = "ip"
    source_nw   = "0.0.0.0/0"
    allow       = false
    description = "deny all"
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `description` - (Optional) The description of the resource.
* `expressions` - (Required) Health check rules. It contains some attributes to [Expressions](#expressions).

### Expressions

Attributes for Expressions:

* `protocol` - (Required) Protocol used in health check.  
Valid value is one of the following: [ "tcp" / "udp" / "icmp" / "fragment" / "ip" ]
* `source_nw` - (Optional) Target source network IP address or CIDR or range.  
Valid format is one of the following:   
  * IP address: `"xxx.xxx.xxx.xxx"`
  * CIDR: `"xxx.xxx.xxx.xxx/nn"`
  * Range: `"xxx.xxx.xxx.xxx/yyy.yyy.yyy.yyy"`
* `source_port` - (Optional) Target source port.
Valid format is one of the following:
  * Number: `"0"` - `"65535"`
  * Range: `"xx-yy"`
  * Range (hex): `"0xPPPP/0xMMMM"`
* `dest_port` - (Optional) Target destination port.
Valid format is one of the following:
  * Number: `"0"` - `"65535"`
  * Range: `"xx-yy"`
  * Range (hex): `"0xPPPP/0xMMMM"`
* `allow` - (Optional) The flag for allow/deny packets (default:`true`).
* `description` - (Optional) The description of the expression.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `description` - The description of the resource.
* `expressions` - Health check rules. It contains some attributes to [Expressions](#expressions).

## Import

Packet Filters can be imported using the Packet Filter ID.

```
$ terraform import sakuracloud_packet_filter.foobar <packet_filter_id>
```
