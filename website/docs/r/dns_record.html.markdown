---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_dns_record"
sidebar_current: "docs-sakuracloud-resource-global-dns-record"
description: |-
  Provides a SakuraCloud DNS Record resource. This can be used to create and delete DNS Records.
---

# sakuracloud\_dns\_record

Provides a SakuraCloud DNS Record resource. This can be used to create and delete DNS Records.

## Example Usage

```hcl
# Create a new DNS(zones)
resource "sakuracloud_dns" "foobar" {
  zone = "example.com"
}

# Create a new DNS Record(Type: A)
resource "sakuracloud_dns_record" "foobar" {
  dns_id = sakuracloud_dns.foobar.id
  name   = "test1"
  type   = "A"
  value  = "192.168.2.1"
}

```

## Argument Reference

The following arguments are supported:

* `dns_id` - (Required) The ID of DNS zones to which the Record belongs.
* `name` - (Required) The hostname of target Record. If "@" is specified, it indicates own zone.
* `type` - (Required) The Record type.  
Valid value is one of the following: [ "A" / "AAAA" / "CNAME" / "NS" / "MX" / "TXT" / "SRV" ]
* `value` - (Required) The value of the Record. 
* `ttl` - (Optional) The ttl value of the Record (unit:`second`). 
* `priority` - (Optional) The priority used when `type` is `MX` or `SRV`.
* `weight` - (Optional) The weight used when `type` is `SRV`.
* `port` - (Optional) The port number used when `type` is `SRV`. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The hostname of target Record. 
* `type` - The Record type.
* `value` - The value of the Record. 
* `ttl` - The ttl value of the Record (unit:`second`). 
* `priority` - The priority used when `type` is `MX` or `SRV`.
* `weight` - The weight used when `type` is `SRV`.
* `port` - The port number used when `type` is `SRV`. 

## Import (not supported)

Import of DNS Record is not supported.
