---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_packet_filter"
sidebar_current: "docs-sakuracloud-datasource-packet-filter"
description: |-
  Get information on a SakuraCloud Packet Filter.
---

# sakuracloud\_packet_filter

Use this data source to retrieve information about a SakuraCloud Packet Filter.

## Example Usage

```hcl
data "sakuracloud_packet_filter" "foobar" {
  name_selectors = ["foobar"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The list of names to filtering.
 * `filter` - (Optional) The map of filter key and value.
 * `zone` - (Optional) The ID of the zone.

## Attributes Reference

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `expressions` - List of filter-expression. It contains some attributes to [Filter Expressions](#filter-expressions).
* `description` - The description of the resource.
* `zone` - The ID of the zone to which the resource belongs.

### Filter Expressions

Attributes for Filter Expressions:

* `protocol` - The target protocol.
* `source_nw` - The source network address (range).
* `source_port` - The source port (range).
* `dest_port` - The destination port (range).
* `allow` - The flag to allow packets. Default value is `true`. 
* `description` - The description of the expression.
