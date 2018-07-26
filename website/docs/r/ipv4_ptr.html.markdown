---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_ipv4_ptr"
sidebar_current: "docs-sakuracloud-resource-ipv4-ptr"
description: |-
  Provides a SakuraCloud IPv4 PTR record resource. This can be used to create, update, and delete IPv4 PTR records.
---

# sakuracloud\_ipv4\_ptr

Provides a SakuraCloud IPv4 PTR Record resource. This can be used to create, update, and delete IPv4 PTR records.

## Example Usage

```hcl
# Create a new IPv4 PTR Record
resource "sakuracloud_ipv4_ptr" "foobar" {
  ipaddress = "192.2.0.1"
  hostname  = "www.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `ipaddress` - (Required) The target IP address.
* `hostname` - (Required) The hostname of target.
* `retry_max` - (Optional) Max count of API call retry.(default:`30`)
* `retry_interval` - (Optional) Interval of API call retry.(unit:`second`, default:`10`)
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `ipaddress` - The target IP address.
* `hostname` - The hostname of target.
* `retry_max` - Max count of API call retry.
* `retry_interval` - Interval of API call retry.(unit:`second`)
* `zone` - The ID of the zone to which the resource belongs.

## Import

IPv4 PTR records can be imported using the IPv4 PTR Record ID.

```
$ terraform import sakuracloud_ipv4_ptr.foobar <ipv4_ptr_id>
```
