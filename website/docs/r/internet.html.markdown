---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_internet"
sidebar_current: "docs-sakuracloud-resource-networking-internet"
description: |-
  Provides a SakuraCloud Internet resource. This can be used to create, update, and delete Internet.
---

# sakuracloud\_internet

Provides a SakuraCloud Internet resource. This can be used to create, update, and delete Internet.

## Example Usage

```hcl
# Create a new Internet
resource "sakuracloud_internet" "foobar" {
  name         = "foobar"
  nw_mask_ken  = 28
  band_width   = 100 
  enable_ipv6  = true
  description  = "description"
  tags         = ["foo", "bar"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `nw_mask_len` - (Optional) Network mask length.  
Valid value is one of the following: [ 28 (default) / 27 / 26 ]
* `band_width` - (Optional) Bandwidth of outbound traffic.(unit:`Mbps`)  
Valid value is one of the following: [ 100 (default) / 250 / 500 / 1000 / 1500 / 2000 / 2500 / 3000 ]
* `enable_ipv6` - (Optional) The ipv6 enabled flag.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the server connected to the resource.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `nw_mask_len` - Network mask length.
* `band_width` - Bandwidth of outbound traffic.
* `switch_id` - The ID of the switch.
* `server_ids` - The IDs of the server connected to the switch.
* `nw_address` - The network address.
* `gateway` - The network gateway address of the switch.
* `min_ipaddress` - Min global IP address.
* `max_ipaddress` - Max global IP address.
* `ipaddresses` - Global IP address list.
* `enable_ipv6` - The ipv6 enabled flag.
* `ipv6_prefix` - Address prefix of ipv6 network.
* `ipv6_prefix_len` - Address prefix length of ipv6 network.
* `ipv6_nw_address` - The ipv6 network address.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Internet can be imported using Internet ID.

```
$ terraform import sakuracloud_internet.foobar <internet_id>
```
