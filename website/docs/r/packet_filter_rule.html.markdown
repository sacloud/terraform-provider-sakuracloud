---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_packet_filter_rule"
sidebar_current: "docs-sakuracloud-resource-networking-pfrule"
description: |-
  Provides a SakuraCloud Packet Filter Rule resource. This can be used to create, update, and delete Packet Filter Rules.
---

# sakuracloud\_packet\_filter\_rule

Provides a SakuraCloud Packet Filter Rule resource. This can be used to create, update, and delete Packet Filter Rules.

## Example Usage

```hcl
# Create a new Packet Filter
resource "sakuracloud_packet_filter" "foobar" {
  name = "foobar"
}

# Create a new packet filter rule
resource "sakuracloud_packet_filter_rule" "rule0" {
  packet_filter_id = sakuracloud_packet_filter.foobar.id

  protocol    = "tcp"
  source_nw   = "0.0.0.0/0"
  source_port = "0-65535"
  dest_port   = "80"
  order       = 0
}

resource "sakuracloud_packet_filter_rule" "rule1" {
  packet_filter_id = sakuracloud_packet_filter.foobar.id

  protocol    = "ip"
  source_nw   = "0.0.0.0/0"
  allow       = false
  description = "deny all"
}

```

## Argument Reference

The following arguments are supported:

* `packet_filter_id` - (Required) The ID of the Packet Filter to which the resource belongs.
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
* `order` - (Optional) The order of the expression (default:`1000`).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `protocol` - Protocol used in health check.  
* `source_nw` - Target source network IP address or CIDR or range.  
* `source_port` - Target source port.
* `dest_port` - Target destination port.
* `allow` - The flag for allow/deny packets (default:`true`).
* `description` - The description of the expression.
* `order` - The order of the expression (default:`1000`).


## Import (not supported)

Import of Packet Filter Rule is not supported.
